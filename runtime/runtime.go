package runtime

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type runtime struct {
	validator *validator.Validate
}

// validatable marks gqlgen structs that carry validation rules.
type validatable interface {
	IsValidatable()
}

// Middleware validates all resolver arguments that satisfy the validatable interface
// after gqlgen unmarshalling.
func Middleware() func(ctx context.Context, next graphql.Resolver) (any, error) {
	r := newRuntime()

	return func(ctx context.Context, next graphql.Resolver) (any, error) {
		if fc := graphql.GetFieldContext(ctx); fc != nil {
			for _, arg := range fc.Args {
				err := r.validate(ctx, arg)
				if err != nil {
					return nil, err
				}
			}
		}

		return next(ctx)
	}
}

func newRuntime() *runtime {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Use the JSON tag name in error messages instead of the Go struct field name because
	// this is the actual name used in the GraphQL schema.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if jsonTag := fld.Tag.Get("json"); jsonTag != "" {
			name := strings.Split(jsonTag, ",")[0]
			if name != "" && name != "-" {
				return name
			}
		}

		return fld.Name
	})

	return &runtime{
		validator: v,
	}
}

// validate runs go-playground/validator against the supplied value and maps
// the errors into the GraphQL response.
func (d *runtime) validate(ctx context.Context, value any) error {
	if !isValidatable(value) {
		return nil
	}

	err := d.validator.StructCtx(ctx, value)
	if err != nil {
		var ves validator.ValidationErrors

		ok := errors.As(err, &ves)
		if !ok || len(ves) == 0 {
			return graphql.ErrorOnPath(ctx, err)
		}

		for i, ve := range ves {
			// We skip the first segment because it's always the struct name.
			path := strings.Split(ve.Namespace(), ".")[1:]
			pathCtx := extendGraphQLPath(ctx, value, path)
			override := fieldMessage(value, ve)
			msg := formatValidationMessage(ve, override)

			err = &gqlerror.Error{
				Message: msg,
				Path:    graphql.GetPath(pathCtx),
				Extensions: map[string]interface{}{
					"code":  "BAD_USER_INPUT",
					"field": ve.Field(),
					"rule":  ve.Tag(),
				},
			}

			// If it's the last error we need to return it so that the request fails.
			if i == len(ves)-1 {
				return graphql.ErrorOnPath(pathCtx, err)
			}

			graphql.AddError(pathCtx, err)
		}
	}

	return nil
}

func isValidatable(value any) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(validatable); !ok {
		return false
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return false
	}

	return true
}

func formatValidationMessage(fe validator.FieldError, override string) string {
	if override != "" {
		return override
	}

	if fe.Param() != "" {
		return fmt.Sprintf("%s failed on the '%s' rule (param: %s)", fe.Field(), fe.Tag(), fe.Param())
	}

	return fmt.Sprintf("%s failed on the '%s' rule", fe.Field(), fe.Tag())
}

func fieldMessage(root any, fe validator.FieldError) string {
	if root == nil {
		return ""
	}

	rt := reflect.TypeOf(root)
	for rt.Kind() == reflect.Pointer {
		if rt.Elem() == nil {
			return ""
		}
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		return ""
	}

	path := fe.StructNamespace()
	if path == "" {
		return fieldTagMessage(rt, fe.StructField())
	}

	segments := strings.Split(path, ".")
	if len(segments) == 0 {
		return ""
	}

	segments = segments[1:]

	current := rt
	for i, segment := range segments {
		segment = trimCollectionIndex(segment)
		if segment == "" {
			continue
		}

		field, ok := current.FieldByName(segment)
		if !ok {
			return ""
		}

		if i == len(segments)-1 {
			return field.Tag.Get("message")
		}

		current = derefType(field.Type)
		if current.Kind() == reflect.Slice || current.Kind() == reflect.Array {
			current = derefType(current.Elem())
		}
		if current.Kind() == reflect.Map {
			current = derefType(current.Elem())
		}
		if current.Kind() != reflect.Struct {
			return ""
		}
	}

	return ""
}

func fieldTagMessage(rt reflect.Type, fieldName string) string {
	field, ok := rt.FieldByName(fieldName)
	if !ok {
		return ""
	}

	return field.Tag.Get("message")
}

func trimCollectionIndex(segment string) string {
	if segment == "" {
		return ""
	}

	idx := strings.Index(segment, "[")
	if idx == -1 {
		return segment
	}

	return segment[:idx]
}

func extendGraphQLPath(ctx context.Context, root any, target []string) context.Context {
	if len(target) == 0 {
		return ctx
	}

	pathCtx := ctx
	typ := reflect.TypeOf(root)
	val := reflect.ValueOf(root)

	for _, raw := range target {
		name, idx := splitNamespaceSegment(raw)
		fieldName, nextType, nextValue := resolvePathSegment(typ, val, name)

		pathCtx = graphql.WithPathContext(pathCtx, graphql.NewPathWithField(fieldName))

		typ = nextType
		val = nextValue

		if idx != nil {
			typ, val = advanceCollection(typ, val, *idx)
			pathCtx = graphql.WithPathContext(pathCtx, graphql.NewPathWithIndex(*idx))
		}
	}

	return pathCtx
}

func resolvePathSegment(typ reflect.Type, val reflect.Value, name string) (string, reflect.Type, reflect.Value) {
	baseName := name
	if typ == nil {
		return baseName, nil, reflect.Value{}
	}

	rootType := derefType(typ)
	if rootType == nil || rootType.Kind() != reflect.Struct {
		return baseName, nil, reflect.Value{}
	}

	field, ok := rootType.FieldByName(name)
	if !ok {
		return baseName, nil, reflect.Value{}
	}

	jsonName := jsonFieldName(field)
	if jsonName == "" {
		jsonName = field.Name
	}

	var nextVal reflect.Value
	currentVal := derefValue(val)
	if currentVal.IsValid() && currentVal.Kind() == reflect.Struct {
		fv := currentVal.FieldByName(field.Name)
		if fv.IsValid() {
			nextVal = fv
		}
	}

	return jsonName, field.Type, nextVal
}

func advanceCollection(typ reflect.Type, val reflect.Value, idx int) (reflect.Type, reflect.Value) {
	if typ == nil {
		return nil, reflect.Value{}
	}

	collectionType := derefType(typ)
	if collectionType == nil {
		return nil, reflect.Value{}
	}

	currentVal := derefValue(val)

	switch collectionType.Kind() {
	case reflect.Slice, reflect.Array:
		elemType := collectionType.Elem()
		if currentVal.IsValid() && idx >= 0 && idx < currentVal.Len() {
			return elemType, currentVal.Index(idx)
		}
		return elemType, reflect.Value{}
	case reflect.Map:
		elemType := collectionType.Elem()
		return elemType, reflect.Value{}
	default:
		return collectionType, currentVal
	}
}

func splitNamespaceSegment(segment string) (string, *int) {
	if segment == "" {
		return "", nil
	}

	idx := strings.Index(segment, "[")
	if idx == -1 {
		return segment, nil
	}

	name := segment[:idx]
	end := strings.Index(segment[idx:], "]")
	if end == -1 {
		return name, nil
	}

	value := segment[idx+1 : idx+end]
	if value == "" {
		return name, nil
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return name, nil
	}

	return name, &n
}

func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return ""
	}
	name := strings.Split(tag, ",")[0]
	if name == "" || name == "-" {
		return ""
	}
	return name
}

func derefValue(v reflect.Value) reflect.Value {
	for v.IsValid() && v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

func derefType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Pointer {
		if t.Elem() == nil {
			return t
		}
		t = t.Elem()
	}
	return t
}
