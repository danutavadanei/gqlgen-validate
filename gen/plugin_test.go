package gen

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

const pluginSchema = `
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

func TestNormalizeDependableRule(t *testing.T) {
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
			assert.Equal(t, tc.expect, toGoRuleParams(tc.input))
		})
	}
}

func TestGetArgumentValueAsString(t *testing.T) {
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

func TestMutateSchemaInjectsGoTagDirectives(t *testing.T) {
	schema := mustLoadSchema(t, pluginSchema)
	p := New().(*Plugin)
	err := p.MutateSchema(schema)
	require.NoError(t, err)

	assert.NotNil(t, schema.Directives[goTagDirectiveName])

	field := schema.Types["SimpleInput"].Fields.ForName("name")
	require.NotNil(t, field)

	var validateTag, messageTag *ast.Directive
	for _, directive := range field.Directives {
		if directive.Name != goTagDirectiveName {
			continue
		}
		key := directive.Arguments.ForName("key")
		value := directive.Arguments.ForName("value")
		require.NotNil(t, key)
		require.NotNil(t, value)
		switch key.Value.Raw {
		case "validate":
			validateTag = directive
		case "message":
			messageTag = directive
		}
	}

	require.NotNil(t, validateTag, "expected validate goTag directive")
	assert.Equal(t, "min=2,required", validateTag.Arguments.ForName("value").Value.Raw)
	require.NotNil(t, messageTag, "expected message goTag directive")
	assert.Equal(t, "name is required", messageTag.Arguments.ForName("value").Value.Raw)

	for _, typeName := range []string{"SimpleInput", "OwnershipInput", "ProxyInput"} {
		assert.Contains(t, p.markerTypes, typeName)
	}
}

func TestMutateSchemaErrors(t *testing.T) {
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
			schema, err := gqlparser.LoadSchema(&ast.Source{Name: "schema.graphql", Input: tc.schema})
			require.NoError(t, err)

			err = New().(*Plugin).MutateSchema(schema)
			require.Error(t, err)
			assert.Equal(t, tc.err, err.Error())
		})
	}
}

func TestMutateConfigAddsDirectiveDefinitions(t *testing.T) {
	cfg := &config.Config{Directives: map[string]config.DirectiveConfig{}}
	err := New().(*Plugin).MutateConfig(cfg)
	require.NoError(t, err)

	assert.True(t, cfg.Directives[goTagDirectiveName].SkipRuntime)
	assert.True(t, cfg.Directives[directiveName].SkipRuntime)
}

func TestMutateConfigRespectsExistingDefinitions(t *testing.T) {
	cfg := &config.Config{Directives: map[string]config.DirectiveConfig{
		goTagDirectiveName: {SkipRuntime: false},
		directiveName:      {SkipRuntime: false},
	}}

	err := New().(*Plugin).MutateConfig(cfg)
	require.NoError(t, err)

	assert.False(t, cfg.Directives[goTagDirectiveName].SkipRuntime)
	assert.False(t, cfg.Directives[directiveName].SkipRuntime)
}

func TestGenerateCodeRemovesStaleFile(t *testing.T) {
	p := &Plugin{markerTypes: make(set)}
	tmpDir := t.TempDir()
	modelPath := filepath.Join(tmpDir, "models_gen.go")
	createConfigPackage(t, modelPath)

	filename := filepath.Join(tmpDir, "validatable_gen.go")
	require.NoError(t, os.WriteFile(filename, []byte("stale"), 0o600))

	data := &codegen.Data{Config: newCodegenConfig(t, modelPath)}

	require.NoError(t, p.GenerateCode(data))
	_, err := os.Stat(filename)
	assert.True(t, errors.Is(err, fs.ErrNotExist))
}

func TestGenerateCodeWritesMarkers(t *testing.T) {
	p := &Plugin{markerTypes: set{
		"OwnershipInput": {},
		"ProxyInput":     {},
	}}

	tmpDir := t.TempDir()
	modelPath := filepath.Join(tmpDir, "models_gen.go")
	createConfigPackage(t, modelPath)

	data := &codegen.Data{Config: newCodegenConfig(t, modelPath)}

	require.NoError(t, p.GenerateCode(data))

	content, err := os.ReadFile(filepath.Join(tmpDir, "validatable_gen.go"))
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "func (OwnershipInput) IsValidatable() {}")
	assert.Contains(t, output, "func (ProxyInput) IsValidatable() {}")
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
