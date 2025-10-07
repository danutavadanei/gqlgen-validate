package runtime

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Middleware validates all resolver arguments that satisfy the validatable interface
// after gqlgen unmarshalling.
func Middleware() func(ctx context.Context, next graphql.Resolver) (any, error) {
	r := newRuntime()
	return func(ctx context.Context, next graphql.Resolver) (any, error) {
		if fc := graphql.GetFieldContext(ctx); fc != nil {
			for _, arg := range fc.Args {
				if err := r.validate(ctx, arg); err != nil {
					return nil, err
				}
			}
		}
		return next(ctx)
	}
}

// validatable marks gqlgen structs that carry validation rules.
type validatable interface {
	IsValidatable()
}

type runtime struct {
	validator  *validator.Validate
	fieldCache sync.Map // map[reflect.Type]map[string]*field
}

type field struct {
	goName   string
	jsonName string
	message  string
	typ      reflect.Type
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

	return &runtime{validator: v}
}

// validate runs go-playground/validator against the supplied value and maps
// the errors into the GraphQL response.
func (r *runtime) validate(ctx context.Context, root any) error {
	if !isValidatable(root) {
		return nil
	}

	if err := r.validator.StructCtx(ctx, root); err != nil {
		var ves validator.ValidationErrors
		if !errors.As(err, &ves) || len(ves) == 0 {
			return graphql.ErrorOnPath(ctx, err)
		}

		for i, ve := range ves {
			pctx := r.getPathContext(ctx, root, ve)
			msg := r.messageFor(root, ve)

			err = &gqlerror.Error{
				Message: msg,
				Path:    graphql.GetPath(pctx),
				Extensions: map[string]any{
					"code":  "BAD_USER_INPUT",
					"field": ve.Field(),
					"rule":  ve.Tag(),
				},
			}

			// If it's the last error we need to return it so that the request fails.
			if i == len(ves)-1 {
				return graphql.ErrorOnPath(pctx, err)
			}
			graphql.AddError(pctx, err)
		}
	}

	return nil
}

func (r *runtime) getPathContext(ctx context.Context, root any, fieldError validator.FieldError) context.Context {
	segments := strings.Split(fieldError.Namespace(), ".")
	if len(segments) == 0 {
		return ctx
	}
	segments = segments[1:]

	pctx := ctx
	rt := reflect.TypeOf(root)
	rv := reflect.ValueOf(root)

	for _, raw := range segments {
		name, idx := parseSegment(raw)

		json, nextT, nextV := r.resolve(rt, rv, name)
		if json == "" {
			json = name
		}
		pctx = graphql.WithPathContext(pctx, graphql.NewPathWithField(json))

		rt, rv = nextT, nextV
		if idx != nil {
			rt, rv = advanceCollection(rt, rv, *idx)
			pctx = graphql.WithPathContext(pctx, graphql.NewPathWithIndex(*idx))
		}
	}
	return pctx
}

func (r *runtime) resolve(typ reflect.Type, value reflect.Value, goName string) (string, reflect.Type, reflect.Value) {
	typ = derefType(typ)
	if typ == nil || typ.Kind() != reflect.Struct {
		return goName, nil, reflect.Value{}
	}

	f := r.fieldFor(typ, goName)
	json := f.jsonName
	nextTyp := f.typ

	value = derefValue(value)
	if value.IsValid() && value.Kind() == reflect.Struct {
		if fv := value.FieldByName(goName); fv.IsValid() {
			return json, nextTyp, fv
		}
	}

	return json, nextTyp, reflect.Value{}
}

func (r *runtime) fieldFor(typ reflect.Type, goName string) *field {
	typ = derefType(typ)

	if cached, ok := r.fieldCache.Load(typ); ok {
		return cached.(map[string]*field)[goName]
	}

	out := make(map[string]*field)
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue
		}

		json := f.Tag.Get("json")
		json = strings.Split(json, ",")[0]
		if json == "" || json == "-" {
			json = f.Name
		}

		fld := &field{
			goName:   f.Name,
			jsonName: json,
			message:  f.Tag.Get("message"),
			typ:      f.Type,
		}

		out[f.Name] = fld
		if json != f.Name {
			out[json] = fld
		}
	}

	r.fieldCache.Store(typ, out)
	return out[goName]
}

func (r *runtime) messageFor(root any, fieldError validator.FieldError) string {
	if msg := r.lookupMessage(root, fieldError); msg != "" {
		return msg
	}
	if p := fieldError.Param(); p != "" {
		return fmt.Sprintf("%s failed on the '%s' rule (param: %s)", fieldError.Field(), fieldError.Tag(), p)
	}
	return fmt.Sprintf("%s failed on the '%s' rule", fieldError.Field(), fieldError.Tag())
}

func (r *runtime) lookupMessage(root any, fieldError validator.FieldError) string {
	if root == nil {
		return ""
	}

	rt := derefType(reflect.TypeOf(root))
	if rt == nil || rt.Kind() != reflect.Struct {
		return ""
	}

	// Direct field override.
	f := r.fieldFor(rt, fieldError.StructField())
	if f != nil && f.message != "" {
		return f.message
	}

	segments := strings.Split(fieldError.StructNamespace(), ".")
	if len(segments) == 0 {
		return ""
	}
	segments = segments[1:]

	curr := rt
	for i, raw := range segments {
		name, _ := parseSegment(raw)
		if name == "" {
			continue
		}

		f = r.fieldFor(curr, name)
		if f == nil {
			return ""
		}
		if i == len(segments)-1 {
			return f.message
		}

		nextT := derefType(f.typ)
		if nextT.Kind() == reflect.Slice || nextT.Kind() == reflect.Array || nextT.Kind() == reflect.Map {
			nextT = derefType(nextT.Elem())
		}
		if nextT == nil || nextT.Kind() != reflect.Struct {
			return ""
		}
		curr = nextT
	}
	return ""
}

func isValidatable(value any) bool {
	if value == nil {
		return false
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return false
	}

	_, ok := value.(validatable)
	return ok
}

func advanceCollection(typ reflect.Type, value reflect.Value, idx int) (reflect.Type, reflect.Value) {
	typ = derefType(typ)
	value = derefValue(value)

	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		et := typ.Elem()
		if value.IsValid() && idx >= 0 && idx < value.Len() {
			return et, value.Index(idx)
		}
		return et, reflect.Value{}
	case reflect.Map:
		return typ.Elem(), reflect.Value{}
	default:
		return typ, value
	}
}

func parseSegment(segment string) (string, *int) {
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
