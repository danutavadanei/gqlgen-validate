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
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Middleware validates all resolver arguments that satisfy the validatable interface
// after gqlgen unmarshalling.
func Middleware(opts ...Option) func(ctx context.Context, next graphql.Resolver) (any, error) {
	r, err := newRuntime(opts...)
	if err != nil {
		panic(fmt.Sprintf("runtime middleware configuration failed: %v", err))
	}
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
	validator    *validator.Validate
	fieldCache   sync.Map // map[reflect.Type]map[string]*field
	translations translationState
}

type field struct {
	goName   string
	jsonName string
	message  string
	typ      reflect.Type
}

type (
	// Option customises the runtime middleware during construction.
	Option func(*settings) error

	settings struct {
		validator    *validator.Validate
		translations translationState
	}

	translationState struct {
		translators map[string]ut.Translator
		defaultLang string
		pickLang    func(context.Context) string
	}
)

// TranslationRegistration wires a translator for a specific language.
type TranslationRegistration struct {
	Lang       string
	Translator ut.Translator
	Init       func(*validator.Validate, ut.Translator) error
}

// TranslationConfig controls translation behaviour for validation errors.
type TranslationConfig struct {
	Registrations []TranslationRegistration
	DefaultLang   string
	PickLang      func(context.Context) string
}

// WithTranslations registers translation handlers that complement directive messages.
func WithTranslations(cfg TranslationConfig) Option {
	return func(s *settings) error {
		if len(cfg.Registrations) == 0 {
			return errors.New("runtime: WithTranslations requires at least one registration")
		}

		translators := make(map[string]ut.Translator, len(cfg.Registrations))
		for _, reg := range cfg.Registrations {
			if reg.Lang == "" {
				return errors.New("runtime: translator language must not be empty")
			}
			if reg.Translator == nil {
				return fmt.Errorf("runtime: translator missing for %q", reg.Lang)
			}
			if _, exists := translators[reg.Lang]; exists {
				return fmt.Errorf("runtime: duplicate translator for %q", reg.Lang)
			}

			if reg.Init != nil {
				if err := reg.Init(s.validator, reg.Translator); err != nil {
					return fmt.Errorf("runtime: initializing translator %q failed: %w", reg.Lang, err)
				}
			}

			translators[reg.Lang] = reg.Translator
		}

		defaultLang := cfg.DefaultLang
		if defaultLang == "" {
			defaultLang = cfg.Registrations[0].Lang
		}
		if _, ok := translators[defaultLang]; !ok {
			return fmt.Errorf("runtime: default language %q is not registered", defaultLang)
		}

		picker := cfg.PickLang
		if picker == nil {
			picker = func(context.Context) string { return defaultLang }
		}

		s.translations = translationState{
			translators: translators,
			defaultLang: defaultLang,
			pickLang:    picker,
		}
		return nil
	}
}

func newRuntime(opts ...Option) (*runtime, error) {
	s := &settings{
		validator: defaultValidator(),
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return &runtime{
		validator:    s.validator,
		translations: s.translations,
	}, nil
}

func defaultValidator() *validator.Validate {
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

	return v
}

// validate runs go-playground/validator against the supplied value and maps
// the errors into the GraphQL response.
func (r *runtime) validate(ctx context.Context, root any) error {
	if !isValidatable(root) {
		return nil
	}

	translator := r.translatorFor(ctx)

	if err := r.validator.StructCtx(ctx, root); err != nil {
		var ves validator.ValidationErrors
		if !errors.As(err, &ves) || len(ves) == 0 {
			return graphql.ErrorOnPath(ctx, err)
		}

		for i, ve := range ves {
			pctx := r.getPathContext(ctx, root, ve)
			msg := r.messageFor(root, ve, translator)

			err = &gqlerror.Error{
				Message: msg,
				Path:    graphql.GetPath(pctx),
				Extensions: map[string]any{
					"code":  "BAD_USER_INPUT",
					"field": ve.Field(),
					"rule":  ve.Tag(),
					"param": ve.Param(),
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
	if len(segments) <= 1 {
		return ctx
	}
	segments = segments[1:]

	pctx := ctx
	rt := reflect.TypeOf(root)
	rv := reflect.ValueOf(root)

	for _, raw := range segments {
		name, idx := parseSegment(raw)

		jsonName, nextT, nextV := r.resolve(rt, rv, name)
		if jsonName == "" {
			jsonName = name
		}
		pctx = graphql.WithPathContext(pctx, graphql.NewPathWithField(jsonName))

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
	if f == nil {
		return goName, nil, reflect.Value{}
	}

	jsonName := f.jsonName
	nextTyp := f.typ

	value = derefValue(value)
	if value.IsValid() && value.Kind() == reflect.Struct {
		if fv := value.FieldByName(goName); fv.IsValid() {
			return jsonName, nextTyp, fv
		}
	}
	return jsonName, nextTyp, reflect.Value{}
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

		jsonName := f.Tag.Get("json")
		jsonName = strings.Split(jsonName, ",")[0]
		if jsonName == "" || jsonName == "-" {
			jsonName = f.Name
		}

		fld := &field{
			goName:   f.Name,
			jsonName: jsonName,
			message:  f.Tag.Get("message"),
			typ:      f.Type,
		}

		out[f.Name] = fld
		if jsonName != f.Name {
			out[jsonName] = fld
		}
	}

	r.fieldCache.Store(typ, out)
	return out[goName]
}

func (r *runtime) messageFor(root any, fieldError validator.FieldError, translator ut.Translator) string {
	if msg := r.lookupMessage(root, fieldError); msg != "" {
		return msg
	}
	if translator != nil {
		if translated := fieldError.Translate(translator); translated != "" {
			return translated
		}
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

	if rv := reflect.ValueOf(value); rv.Kind() == reflect.Ptr && rv.IsNil() {
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

	name, rest, ok := strings.Cut(segment, "[")
	if !ok {
		return segment, nil
	}

	rest = strings.TrimSuffix(rest, "]")
	if rest == "" {
		return name, nil
	}

	if n, err := strconv.Atoi(rest); err == nil {
		return name, &n
	}

	return name, nil
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
		t = t.Elem()
	}
	return t
}

func (r *runtime) translatorFor(ctx context.Context) ut.Translator {
	ts := r.translations
	if ts.pickLang == nil || len(ts.translators) == 0 {
		return nil
	}

	lang := ts.pickLang(ctx)
	if lang != "" {
		if tr, ok := ts.translators[lang]; ok {
			return tr
		}
	}

	if ts.defaultLang == "" {
		return nil
	}
	return ts.translators[ts.defaultLang]
}
