package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin/modelgen"
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

func loadPluginSchema(t *testing.T) *ast.Schema {
	t.Helper()

	schema, err := gqlparser.LoadSchema(&ast.Source{Input: pluginSchema, Name: "schema.graphql"})
	if err != nil {
		t.Fatalf("load schema: %v", err)
	}

	return schema
}

func TestModelPlugin_FieldTags(t *testing.T) {
	schema := loadPluginSchema(t)
	plugin := New().(*ModelPlugin)
	require.NoError(t, plugin.collectRules(schema))
	plugin.seedGoNames(schema)

	cases := []struct {
		name    string
		obj     string
		field   string
		initial string
		want    string
	}{
		{
			name:    "simple",
			obj:     "SimpleInput",
			field:   "name",
			initial: `json:"name"`,
			want:    `json:"name" validate:"min=2,required" message:"name is required"`,
		},
		{
			name:    "ownership",
			obj:     "OwnershipInput",
			field:   "legalName",
			initial: `json:"legalName"`,
			want:    `json:"legalName" validate:"required_if=UserOwned false"`,
		},
		{
			name:    "proxy",
			obj:     "ProxyInput",
			field:   "target",
			initial: `json:"target"`,
			want:    `json:"target" validate:"required_if=Proxy true"`,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s_%s", tc.obj, tc.field), func(t *testing.T) {
			field := schema.Types[tc.obj].Fields.ForName(tc.field)
			require.NotNil(t, field)

			goName := plugin.goNames[tc.obj][tc.field]
			require.NotEmpty(t, goName)

			out, err := plugin.fieldHook(schema.Types[tc.obj], field, &modelgen.Field{
				Name:   tc.field,
				GoName: goName,
				Tag:    tc.initial,
			})
			require.NoError(t, err)
			require.NotNil(t, out)
			assert.Equal(t, tc.want, out.Tag)
		})
	}
}

func TestModelPlugin_InlineDirectiveDisallowed(t *testing.T) {
	const invalidSchema = `
		directive @validate(rule: String!, message: String) on INPUT_FIELD_DEFINITION | ARGUMENT_DEFINITION

		input OwnershipInput {
			userOwned: Boolean!
			legalName: String @validate(rule: "required_if=userOwned false")
		}

		type Mutation {
			registerUser(input: OwnershipInput!): Boolean!
			test(name: String! @validate(rule: "required")): Boolean!
		}
	`

	schema, err := gqlparser.LoadSchema(&ast.Source{Input: invalidSchema, Name: "invalid.graphql"})
	if err != nil {
		t.Fatalf("load schema: %v", err)
	}
	mutation := schema.Types["Mutation"]
	if mutation == nil {
		t.Fatalf("mutation type missing")
	}
	arg := mutation.Fields.ForName("test").Arguments.ForName("name")
	if arg == nil {
		t.Fatalf("argument missing")
	}
	if len(arg.Directives) == 0 {
		t.Fatalf("expected directive on argument")
	}
	if arg.Directives[0].Name != directiveName {
		t.Fatalf("unexpected directive name: %s", arg.Directives[0].Name)
	}
	if err := ensureNoArgumentRules(mutation); err == nil {
		t.Fatalf("ensureNoArgumentRules should fail")
	}

	plugin := New().(*ModelPlugin)
	err = plugin.collectRules(schema)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not supported on argument")
}

func TestSetStructTag(t *testing.T) {
	tag := setStructTag(`json:"name"`, "validate", "required")
	assert.Equal(t, `json:"name" validate:"required"`, tag)

	updated := setStructTag(tag, "validate", "min=3")
	assert.Equal(t, `json:"name" validate:"min=3"`, updated)

	withMessage := setStructTag(updated, "message", "custom")
	assert.Equal(t, `json:"name" validate:"min=3" message:"custom"`, withMessage)
}

func TestConvertRule(t *testing.T) {
	p := &ModelPlugin{
		goNames: map[string]map[string]string{
			"SampleInput": {
				"field": "Field",
				"other": "Other",
			},
		},
	}

	assert.Equal(t, "required_without=Field", p.convertRule("SampleInput", "required_without=field"))
	assert.Equal(t, "unique=Field:Other", p.convertRule("SampleInput", "unique=field:other"))
	assert.Equal(t, "required_without=Field,Other", p.convertRule("SampleInput", "required_without=field,other"))
	assert.Equal(t, "required_without=Field Other", p.convertRule("SampleInput", "required_without=field other"))
	assert.Equal(t, "min=3", p.convertRule("SampleInput", "min=3"))
}

func TestRulesForField(t *testing.T) {
	p := &ModelPlugin{
		rules: map[string]map[string][]string{
			"Input": {
				"name": {" required ", "min=2", "required"},
			},
		},
	}

	rules := p.rulesForField("Input", "name")
	assert.Equal(t, []string{"min=2", "required"}, rules)
}

func TestCollectRulesUnexpectedArgument(t *testing.T) {
	const schemaStr = `
		directive @validate(rule: String!, message: String, extra: String) on INPUT_FIELD_DEFINITION

		input SampleInput {
			name: String @validate(rule: "required", extra: "nope")
		}
	`

	schema, err := gqlparser.LoadSchema(&ast.Source{Input: schemaStr, Name: "invalid.graphql"})
	require.NoError(t, err)

	plugin := New().(*ModelPlugin)
	err = plugin.collectRules(schema)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only supports the 'rule' and 'message' arguments")
}

func TestCollectRulesMessageStored(t *testing.T) {
	const schemaStr = `
		directive @validate(rule: String!, message: String) on INPUT_FIELD_DEFINITION

		input SampleInput {
			name: String @validate(rule: "required", message: " custom message ")
		}
	`

	schema, err := gqlparser.LoadSchema(&ast.Source{Input: schemaStr, Name: "schema.graphql"})
	require.NoError(t, err)

	plugin := New().(*ModelPlugin)
	require.NoError(t, plugin.collectRules(schema))

	assert.Equal(t, "custom message", plugin.messages["SampleInput"]["name"])
}

func TestCollectRulesDuplicateDirective(t *testing.T) {
	const schemaStr = `
		directive @validate(rule: String!, message: String) on INPUT_FIELD_DEFINITION

		input SampleInput {
			name: String @validate(rule: "required") @validate(rule: "min=2")
		}
	`

	schema, err := gqlparser.LoadSchema(&ast.Source{Input: schemaStr, Name: "invalid.graphql"})
	require.NoError(t, err)

	plugin := New().(*ModelPlugin)
	err = plugin.collectRules(schema)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "may only be applied once per field")
}

func TestModelPlugin_Name(t *testing.T) {
	plugin := New().(*ModelPlugin)
	assert.Equal(t, "modelgen", plugin.Name())
}

func TestModelPlugin_MutateConfigNil(t *testing.T) {
	plugin := New().(*ModelPlugin)
	assert.Error(t, plugin.MutateConfig(nil))
}

func TestRenderMarkersRemovesFile(t *testing.T) {
	p := &ModelPlugin{rules: map[string]map[string][]string{}}

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "models_gen.go")
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "validatable_gen.go"), []byte("test"), 0o644))

	cfg := &config.Config{
		Model: config.PackageConfig{Filename: filename},
	}

	assert.NoError(t, p.renderMarkers(cfg))
	assert.NoFileExists(t, filepath.Join(tmpDir, "validatable_gen.go"))
}

func TestModelPlugin_MutateConfigSuccess(t *testing.T) {
	const schemaStr = `
		directive @validate(rule: String!, message: String) on INPUT_FIELD_DEFINITION

		input SampleInput {
			name: String @validate(rule: "required", message: "hi")
		}

		type Query {
			lookup(input: SampleInput!): Boolean!
		}
	`

	schema, err := gqlparser.LoadSchema(&ast.Source{Input: schemaStr, Name: "schema.graphql"})
	require.NoError(t, err)

	originalRender := renderTemplate
	defer func() { renderTemplate = originalRender }()

	var rendered templates.Options
	renderTemplate = func(opts templates.Options) error {
		rendered = opts
		return nil
	}

	plugin := New().(*ModelPlugin)
	plugin.mutate = func(*config.Config) error { return nil }

	tmpDir := t.TempDir()
	cfg := &config.Config{
		Schema: schema,
		Model: config.PackageConfig{
			Filename: filepath.Join(tmpDir, "models_gen.go"),
			Package:  "graph",
		},
		Models: config.TypeMap{},
	}

	require.NoError(t, plugin.MutateConfig(cfg))
	assert.Equal(t, "graph", rendered.PackageName)
	assert.Equal(t, filepath.Join(tmpDir, "validatable_gen.go"), rendered.Filename)
	assert.NotEmpty(t, rendered.Data)
}

func TestArgumentValueRaw(t *testing.T) {
	dir := &ast.Directive{
		Arguments: ast.ArgumentList{
			&ast.Argument{
				Name:  "rule",
				Value: &ast.Value{Raw: "required"},
			},
		},
	}

	assert.Equal(t, "required", argumentValue(dir, "rule"))

	dir.Arguments[0].Value = &ast.Value{
		Kind: ast.ListValue,
		Children: ast.ChildValueList{
			{Value: &ast.Value{Raw: "1", Kind: ast.IntValue}},
		},
	}
	assert.Equal(t, "[1]", argumentValue(dir, "rule"))
}

func TestModelPlugin_MutateConfigPropagatesError(t *testing.T) {
	const schemaStr = `
		type Query {
			ping: Boolean!
		}
	`

	schema, err := gqlparser.LoadSchema(&ast.Source{Input: schemaStr, Name: "schema.graphql"})
	require.NoError(t, err)

	plugin := New().(*ModelPlugin)
	plugin.mutate = func(*config.Config) error { return fmt.Errorf("boom") }

	cfg := &config.Config{
		Schema: schema,
		Model:  config.PackageConfig{Filename: filepath.Join(t.TempDir(), "models_gen.go")},
		Models: config.TypeMap{},
	}

	err = plugin.MutateConfig(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestParseStructTag(t *testing.T) {
	parts := parseStructTag(`json:"name" validate:"required" message:"hi"`)
	assert.Equal(t, []tagPart{{key: "json", value: "name"}, {key: "validate", value: "required"}, {key: "message", value: "hi"}}, parts)

	assert.Nil(t, parseStructTag(""))
	assert.Empty(t, parseStructTag("incomplete"))
}
