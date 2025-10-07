package gen

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

const (
	schemaWithMessages = `
    directive @validate(rule: String!, message: String) on INPUT_FIELD_DEFINITION | ARGUMENT_DEFINITION

    input SimpleInput {
        name: String @validate(rule: "min=2,required", message: "name is required")
    }

    input OwnershipInput {
        userOwned: Boolean!
        legalName: String @validate(rule: "required_if=userOwned false")
    }

    input ProxyInput {
        proxy: Boolean @validate(rule: "required")
        target: String @validate(rule: "required_if=proxy true")
    }

    type Mutation {
        registerUser(input: OwnershipInput!): Boolean!
    }
`

	schemaWithoutMessage = `
    directive @validate(rule: String!, message: String) on INPUT_FIELD_DEFINITION

    input MinimalInput {
        field: String @validate(rule: "required")
    }
`
)

func TestToGoRuleParams(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "empty", input: "", expect: ""},
		{name: "unchanged-required", input: "required", expect: "required"},
		{name: "unchanged-min", input: "min=10", expect: "min=10"},
		{name: "pipe-mixed", input: "eqfield=confirmPassword|required_with=email", expect: "eqfield=ConfirmPassword|required_with=Email"},
		{name: "comma-mixed", input: "eqfield=confirmPassword,required_with=email", expect: "eqfield=ConfirmPassword,required_with=Email"},
		{name: "eqfield", input: "eqfield=questionId", expect: "eqfield=QuestionID"},
		{name: "eqcsfield", input: "eqcsfield=parent.child", expect: "eqcsfield=Parent.Child"},
		{name: "eqsfield", input: "eqsfield=parent.child", expect: "eqsfield=Parent.Child"},
		{name: "required_with", input: "required_with=email phone", expect: "required_with=Email Phone"},
		{name: "required_if", input: "required_if=otherField foo anotherField bar", expect: "required_if=OtherField foo AnotherField bar"},
		{name: "paired-odd", input: "required_if=field1 foo field2", expect: "required_if=field1 foo field2"},
		{name: "pipe-trims", input: " required_with=email | required_if=otherField foo", expect: "required_with=Email|required_if=OtherField foo"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expect, toGoRuleParams(tc.input))
		})
	}
}

func TestGetArgumentValueAsString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		arg         *ast.Argument
		expectValue string
		expectErr   string
	}{
		{name: "nil", arg: nil, expectErr: "argument is nil"},
		{name: "non-string", arg: &ast.Argument{Name: "rule", Value: &ast.Value{Kind: ast.IntValue, Raw: "5"}}, expectErr: "argument value is not a string"},
		{name: "empty", arg: &ast.Argument{Name: "rule", Value: &ast.Value{Kind: ast.StringValue, Raw: ""}}, expectErr: "argument value is an empty string"},
		{name: "valid", arg: &ast.Argument{Name: "rule", Value: &ast.Value{Kind: ast.StringValue, Raw: "required"}}, expectValue: "required"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			value, err := getArgumentValueAsString(tc.arg)
			if tc.expectErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectValue, value)
		})
	}
}

func TestPluginMutateSchema(t *testing.T) {
	t.Run("injects goTags and tracks marker types", func(t *testing.T) {
		schema := mustLoadSchema(t, schemaWithMessages)
		plugin := New().(*Plugin)

		require.NoError(t, plugin.MutateSchema(schema))

		field := schema.Types["SimpleInput"].Fields.ForName("name")
		require.NotNil(t, field)

		assert.Equal(t, "min=2,required", goTagValue(t, field, "validate"))
		assert.Equal(t, "name is required", goTagValue(t, field, "message"))

		assert.ElementsMatch(t, []string{"OwnershipInput", "ProxyInput", "SimpleInput"}, plugin.markerTypes.values())
	})

	t.Run("skips message directive when argument absent", func(t *testing.T) {
		schema := mustLoadSchema(t, schemaWithoutMessage)
		plugin := New().(*Plugin)

		require.NoError(t, plugin.MutateSchema(schema))

		field := schema.Types["MinimalInput"].Fields.ForName("field")
		require.NotNil(t, field)

		assert.Equal(t, "required", goTagValue(t, field, "validate"))
		assert.False(t, hasGoTag(field, "message"))

		assert.ElementsMatch(t, []string{"MinimalInput"}, plugin.markerTypes.values())
	})
}

func TestPluginMutateSchemaErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		schema string
		err    string
	}{
		{
			name: "directive on type",
			schema: `
                directive @validate(rule: String!) on INPUT_FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_OBJECT

                input BadInput @validate(rule: "required") {
                    field: String
                }
            `,
			err: "@validate may only be applied to input fields (found on BadInput)",
		},
		{
			name: "empty rule",
			schema: `
                directive @validate(rule: String!, message: String) on INPUT_FIELD_DEFINITION

                input BadInput {
                    name: String @validate(rule: "", message: "oops")
                }
            `,
			err: "@validate on BadInput.name requires a rule",
		},
		{
			name: "duplicate directive",
			schema: `
                directive @validate(rule: String!) on INPUT_FIELD_DEFINITION

                input BadInput {
                    name: String @validate(rule: "required") @validate(rule: "min=2")
                }
            `,
			err: "@validate may only be applied once per field (BadInput.name)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			schema, err := gqlparser.LoadSchema(&ast.Source{Name: "schema.graphql", Input: tc.schema})
			require.NoError(t, err)

			err = New().(*Plugin).MutateSchema(schema)
			require.Error(t, err)
			assert.Equal(t, tc.err, err.Error())
		})
	}
}

func TestPluginMutateConfig(t *testing.T) {
	t.Run("adds directive definitions", func(t *testing.T) {
		cfg := &config.Config{Directives: map[string]config.DirectiveConfig{}}

		err := New().(*Plugin).MutateConfig(cfg)
		require.NoError(t, err)

		assert.True(t, cfg.Directives[goTagDirectiveName].SkipRuntime)
		assert.True(t, cfg.Directives[directiveName].SkipRuntime)
	})

	t.Run("respects existing definitions", func(t *testing.T) {
		cfg := &config.Config{Directives: map[string]config.DirectiveConfig{
			goTagDirectiveName: {SkipRuntime: false},
			directiveName:      {SkipRuntime: false},
		}}

		err := New().(*Plugin).MutateConfig(cfg)
		require.NoError(t, err)

		assert.False(t, cfg.Directives[goTagDirectiveName].SkipRuntime)
		assert.False(t, cfg.Directives[directiveName].SkipRuntime)
	})
}

func TestPluginGenerateCode(t *testing.T) {
	t.Run("removes stale file when no marker types", func(t *testing.T) {
		plugin := &Plugin{markerTypes: make(set)}
		tmpDir := t.TempDir()
		modelPath := filepath.Join(tmpDir, "models_gen.go")
		createConfigPackage(t, modelPath)

		filename := filepath.Join(tmpDir, "validatable_gen.go")
		require.NoError(t, os.WriteFile(filename, []byte("stale"), 0o600))

		data := &codegen.Data{Config: newCodegenConfig(t, modelPath)}

		require.NoError(t, plugin.GenerateCode(data))
		_, err := os.Stat(filename)
		assert.True(t, errors.Is(err, fs.ErrNotExist))
	})

	t.Run("writes sorted marker functions", func(t *testing.T) {
		plugin := &Plugin{markerTypes: set{
			"ProxyInput":   {},
			"AlphaInput":   {},
			"MonitorInput": {},
		}}

		tmpDir := t.TempDir()
		modelPath := filepath.Join(tmpDir, "models_gen.go")
		createConfigPackage(t, modelPath)

		data := &codegen.Data{Config: newCodegenConfig(t, modelPath)}

		require.NoError(t, plugin.GenerateCode(data))

		content, err := os.ReadFile(filepath.Join(tmpDir, "validatable_gen.go"))
		require.NoError(t, err)

		output := string(content)
		alphaIdx := strings.Index(output, "func (AlphaInput) IsValidatable()")
		monitorIdx := strings.Index(output, "func (MonitorInput) IsValidatable()")
		proxyIdx := strings.Index(output, "func (ProxyInput) IsValidatable()")

		require.NotEqual(t, -1, alphaIdx)
		require.NotEqual(t, -1, monitorIdx)
		require.NotEqual(t, -1, proxyIdx)
		assert.Less(t, alphaIdx, monitorIdx)
		assert.Less(t, monitorIdx, proxyIdx)
	})
}

func goTagValue(t *testing.T, field *ast.FieldDefinition, key string) string {
	t.Helper()

	for _, directive := range field.Directives {
		if directive.Name != goTagDirectiveName {
			continue
		}

		keyArg := directive.Arguments.ForName("key")
		valueArg := directive.Arguments.ForName("value")
		require.NotNil(t, keyArg)
		require.NotNil(t, valueArg)

		if keyArg.Value.Raw == key {
			return valueArg.Value.Raw
		}
	}

	t.Fatalf("goTag %q not found", key)
	return ""
}

func hasGoTag(field *ast.FieldDefinition, key string) bool {
	for _, directive := range field.Directives {
		if directive.Name != goTagDirectiveName {
			continue
		}

		keyArg := directive.Arguments.ForName("key")
		if keyArg != nil && keyArg.Value.Raw == key {
			return true
		}
	}
	return false
}

func mustLoadSchema(t *testing.T, input string) *ast.Schema {
	t.Helper()

	schema, err := gqlparser.LoadSchema(&ast.Source{Name: "schema.graphql", Input: input})
	require.NoError(t, err)
	return schema
}

func createConfigPackage(t *testing.T, modelPath string) {
	t.Helper()

	dir := filepath.Dir(modelPath)
	require.NoError(t, os.MkdirAll(dir, 0o755))
}

func newCodegenConfig(t *testing.T, modelPath string) *config.Config {
	t.Helper()

	cfg := &config.Config{
		Model: config.PackageConfig{
			Filename: modelPath,
			Package:  "model",
		},
		Directives: map[string]config.DirectiveConfig{},
	}

	ensurePackagesInitialized(t, cfg)
	return cfg
}

func ensurePackagesInitialized(t *testing.T, cfg *config.Config) {
	t.Helper()

	value := reflect.ValueOf(cfg).Elem().FieldByName("Packages")
	if !value.IsNil() {
		return
	}

	value.Set(reflect.New(value.Type().Elem()))
}
