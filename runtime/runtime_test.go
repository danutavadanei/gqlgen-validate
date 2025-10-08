package runtime

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	enlocale "github.com/go-playground/locales/en"
	rolocale "github.com/go-playground/locales/ro"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type simpleInput struct {
	Name string `json:"name" validate:"required" message:"name must not be empty"`
	Age  int    `json:"age"`
}

func (simpleInput) IsValidatable() {}

type ownershipInput struct {
	UserOwned bool    `json:"userOwned"`
	LegalName *string `json:"legalName" validate:"required_if=UserOwned false"`
}

func (ownershipInput) IsValidatable() {}

type nestedInner struct {
	Message string `json:"message" validate:"min=2" message:"message too short"`
}

type nestedOuter struct {
	Inner nestedInner `json:"inner"`
}

func (nestedOuter) IsValidatable() {}

type listRoot struct {
	Items []nestedInner `json:"items" validate:"dive"`
}

func (listRoot) IsValidatable() {}

type simplePointer struct {
	Name string `json:"name" validate:"required"`
}

func (simplePointer) IsValidatable() {}

type translationInput struct {
	Name string `json:"name" validate:"required"`
}

func (translationInput) IsValidatable() {}

type langCtxKey struct{}

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		wantErr  string
		wantPath string
	}{
		{
			name:  "simple ok",
			value: &simpleInput{Name: "Alice"},
		},
		{
			name:     "simple missing name",
			value:    &simpleInput{},
			wantErr:  "name must not be empty",
			wantPath: "input.name",
		},
		{
			name:  "ownership ok",
			value: &ownershipInput{UserOwned: true},
		},
		{
			name:     "ownership missing legal name",
			value:    &ownershipInput{UserOwned: false},
			wantErr:  "legalName failed on the 'required_if' rule (param: UserOwned false)",
			wantPath: "input.legalName",
		},
		{
			name:  "nested ok",
			value: &nestedOuter{Inner: nestedInner{Message: "ok"}},
		},
		{
			name:     "nested too short",
			value:    &nestedOuter{Inner: nestedInner{Message: "a"}},
			wantErr:  "message too short",
			wantPath: "input.inner.message",
		},
		{
			name:     "list element",
			value:    &listRoot{Items: []nestedInner{{Message: "a"}}},
			wantErr:  "message too short",
			wantPath: "input.items[0].message",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := newTestRuntime(t)
			ctx := graphql.WithPathContext(context.Background(), graphql.NewPathWithField("input"))

			err := d.validate(ctx, tc.value)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)
			var gqlErr *gqlerror.Error
			require.True(t, errors.As(err, &gqlErr))
			assert.Equal(t, tc.wantErr, gqlErr.Message)
			if tc.wantPath != "" {
				assert.Equal(t, tc.wantPath, gqlErr.Path.String())
			}
		})
	}
}

func TestValidateWithTranslations(t *testing.T) {
	opt := translationOption(t)
	r := newTestRuntime(t, opt)

	baseCtx := graphql.WithPathContext(context.Background(), graphql.NewPathWithField("input"))

	t.Run("uses requested translator", func(t *testing.T) {
		ctx := context.WithValue(baseCtx, langCtxKey{}, "ro")

		err := r.validate(ctx, &translationInput{})
		require.Error(t, err)

		var gqlErr *gqlerror.Error
		require.True(t, errors.As(err, &gqlErr))
		assert.Equal(t, "name trebuie completat", gqlErr.Message)
		assert.Equal(t, "input.name", gqlErr.Path.String())
	})

	t.Run("falls back to default translator", func(t *testing.T) {
		ctx := context.WithValue(baseCtx, langCtxKey{}, "fr")

		err := r.validate(ctx, &translationInput{})
		require.Error(t, err)

		var gqlErr *gqlerror.Error
		require.True(t, errors.As(err, &gqlErr))
		assert.Equal(t, "name is required (en)", gqlErr.Message)
		assert.Equal(t, "input.name", gqlErr.Path.String())
	})

	t.Run("directive message overrides translations", func(t *testing.T) {
		ctx := context.WithValue(baseCtx, langCtxKey{}, "ro")

		err := r.validate(ctx, &simpleInput{})
		require.Error(t, err)

		var gqlErr *gqlerror.Error
		require.True(t, errors.As(err, &gqlErr))
		assert.Equal(t, "name must not be empty", gqlErr.Message)
	})
}

func TestIsValidatable(t *testing.T) {
	cases := []struct {
		name  string
		value any
		want  bool
	}{
		{"pointer", &simplePointer{}, true},
		{"value", simplePointer{}, true},
		{"nil", (*simplePointer)(nil), false},
		{"non-struct", 42, false},
		{"struct no marker", struct{}{}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, isValidatable(tc.value))
		})
	}
}

func TestMiddleware(t *testing.T) {
	mw := Middleware()

	ctx := graphql.WithPathContext(context.Background(), graphql.NewPathWithField("input"))

	t.Run("valid", func(t *testing.T) {
		fc := &graphql.FieldContext{Args: map[string]any{"input": &simpleInput{Name: "Alice"}}}
		localCtx := graphql.WithFieldContext(ctx, fc)
		called := false
		res, err := mw(localCtx, func(ctx context.Context) (any, error) {
			called = true
			return "ok", nil
		})
		require.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, "ok", res)
	})

	t.Run("invalid", func(t *testing.T) {
		fc := &graphql.FieldContext{Args: map[string]any{"input": &simpleInput{}}}
		localCtx := graphql.WithFieldContext(ctx, fc)
		res, err := mw(localCtx, func(ctx context.Context) (any, error) { return "ok", nil })
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestLookupMessage(t *testing.T) {
	type child struct {
		Name string `json:"name" validate:"required" message:"child message"`
	}

	type parent struct {
		Direct    string            `json:"direct" validate:"required" message:"direct message"`
		Nested    child             `json:"nested"`
		NestedPtr *child            `json:"nestedPtr"`
		List      []child           `json:"list" validate:"dive"`
		ListPtr   []*child          `json:"listPtr" validate:"dive"`
		Map       map[string]child  `json:"map" validate:"dive"`
		MapPtr    map[string]*child `json:"mapPtr" validate:"dive"`
		Value     []string          `json:"value" validate:"dive,required"`
	}

	root := &parent{
		Nested:    child{},
		NestedPtr: &child{},
		List:      []child{{}},
		ListPtr:   []*child{&child{}},
		Map:       map[string]child{"key": {}},
		MapPtr:    map[string]*child{"key": &child{}},
		Value:     []string{""},
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if jsonTag := fld.Tag.Get("json"); jsonTag != "" {
			name := strings.Split(jsonTag, ",")[0]
			if name != "" && name != "-" {
				return name
			}
		}

		return fld.Name
	})

	err := v.Struct(root)
	var ves validator.ValidationErrors
	require.Error(t, err)
	require.True(t, errors.As(err, &ves))

	byNamespace := make(map[string]validator.FieldError, len(ves))
	for _, fe := range ves {
		byNamespace[fe.StructNamespace()] = fe
	}

	get := func(suffix string) validator.FieldError {
		t.Helper()
		for ns, fe := range byNamespace {
			if strings.HasSuffix(ns, suffix) {
				return fe
			}
		}

		keys := make([]string, 0, len(byNamespace))
		for ns := range byNamespace {
			keys = append(keys, ns)
		}
		t.Fatalf("no field error with suffix %q (namespaces: %v)", suffix, keys)
		return nil
	}

	varValidator := validator.New()
	varErr := varValidator.Var("", "required")
	var blank validator.ValidationErrors
	require.Error(t, varErr)
	require.True(t, errors.As(varErr, &blank))
	blanks := blank[0]

	direct := get(".Direct")

	tests := []struct {
		name string
		root any
		err  validator.FieldError
		want string
	}{
		{"nil root", nil, direct, ""},
		{"non struct root", 42, direct, ""},
		{"direct field", root, direct, "direct message"},
		{"missing field on root", struct{}{}, direct, ""},
		{"nested struct", root, get(".Nested.Name"), "child message"},
		{"nested pointer", root, get(".NestedPtr.Name"), "child message"},
		{"slice element", root, get(".List[0].Name"), "child message"},
		{"slice pointer element", root, get(".ListPtr[0].Name"), "child message"},
		{"map element", root, get(".Map[key].Name"), "child message"},
		{"map pointer element", root, get(".MapPtr[key].Name"), "child message"},
		{"non struct collection element", root, get(".Value[0]"), ""},
		{"empty namespace", root, blanks, ""},
		{"empty segment", root, overrideNamespace(direct, "parent..Direct"), "direct message"},
		{"no segments after root", root, overrideNamespace(direct, "parent"), "direct message"},
	}

	r := newTestRuntime(t)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := r.lookupMessage(tc.root, tc.err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func overrideNamespace(fe validator.FieldError, ns string) validator.FieldError {
	return namespaceOverride{FieldError: fe, structNamespace: ns}
}

type namespaceOverride struct {
	validator.FieldError
	structNamespace string
}

func (n namespaceOverride) StructNamespace() string {
	if n.structNamespace != "" {
		return n.structNamespace
	}

	return n.FieldError.StructNamespace()
}

func TestFieldMessageTag(t *testing.T) {
	typ := reflect.TypeOf(struct {
		First string `message:"first"`
		Last  string
	}{})

	r := newTestRuntime(t)
	assert.Equal(t, "first", r.fieldFor(typ, "First").message)
	assert.Nil(t, r.fieldFor(typ, "Missing"))
}

func TestWithTranslationsConfigValidation(t *testing.T) {
	_, err := newRuntime(WithTranslations(TranslationConfig{}))
	assert.Error(t, err)

	en := enlocale.New()
	uni := ut.New(en, en)
	enTrans, found := uni.GetTranslator(en.Locale())
	require.True(t, found)

	cfg := TranslationConfig{
		Registrations: []TranslationRegistration{
			{Lang: "en", Translator: enTrans},
		},
		DefaultLang: "unknown",
	}

	_, err = newRuntime(WithTranslations(cfg))
	assert.Error(t, err)
}

func newTestRuntime(t *testing.T, opts ...Option) *runtime {
	t.Helper()

	r, err := newRuntime(opts...)
	require.NoError(t, err)
	return r
}

func translationOption(t *testing.T) Option {
	t.Helper()

	en := enlocale.New()
	enUni := ut.New(en, en)
	enTrans, found := enUni.GetTranslator(en.Locale())
	require.True(t, found)

	ro := rolocale.New()
	roUni := ut.New(ro, ro)
	roTrans, found := roUni.GetTranslator(ro.Locale())
	require.True(t, found)

	return WithTranslations(TranslationConfig{
		Registrations: []TranslationRegistration{
			{
				Lang:       "en",
				Translator: enTrans,
				Init:       requiredInit("{0} is required (en)"),
			},
			{
				Lang:       "ro",
				Translator: roTrans,
				Init:       requiredInit("{0} trebuie completat"),
			},
		},
		DefaultLang: "en",
		PickLang: func(ctx context.Context) string {
			if lang, ok := ctx.Value(langCtxKey{}).(string); ok {
				return lang
			}
			return ""
		},
	})
}

func requiredInit(message string) func(*validator.Validate, ut.Translator) error {
	return func(v *validator.Validate, trans ut.Translator) error {
		return v.RegisterTranslation("required", trans,
			func(ut ut.Translator) error {
				return ut.Add("required", message, true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				translated, err := ut.T("required", fe.Field())
				if err != nil {
					return message
				}
				return translated
			},
		)
	}
}
