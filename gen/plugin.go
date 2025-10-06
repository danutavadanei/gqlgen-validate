package gen

import (
	_ "embed"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
)

const (
	// directiveName identifies the GraphQL directive handled by this package.
	directiveName = "validate"

	// goTagDirectiveName identifies the gqlgen directive used to inject struct tags.
	goTagDirectiveName = "goTag"
)

var (
	renderTemplate = templates.Render

	//go:embed markers.gotpl
	markersTemplate string
)

type Plugin struct {
	markerTypes map[string]struct{}
}

var (
	_ plugin.CodeGenerator = &Plugin{}
	_ plugin.ConfigMutator = &Plugin{}
	_ plugin.SchemaMutator = &Plugin{}
)

func New() plugin.Plugin {
	return &Plugin{
		markerTypes: make(map[string]struct{}),
	}
}

func (p *Plugin) Name() string {
	return "gqlgen-validate"
}

func (p *Plugin) MutateSchema(schema *ast.Schema) error {
	if _, ok := schema.Directives[goTagDirectiveName]; !ok {
		schema.Directives[goTagDirectiveName] = &ast.DirectiveDefinition{
			Name: goTagDirectiveName,
		}
	}

	for typeName, def := range schema.Types {
		if def.Kind != ast.InputObject {
			continue
		}

		if d := def.Directives.ForName(directiveName); d != nil {
			return fmt.Errorf("@%s may only be applied to input fields (found on %s)", directiveName, def.Name)
		}

		hasValidateDirectives := false

		for _, field := range def.Fields {
			validateDirectives := field.Directives.ForNames(directiveName)

			if len(validateDirectives) == 0 {
				continue
			}

			if len(validateDirectives) > 1 {
				return fmt.Errorf("@%s may only be applied once per field (%s.%s)", directiveName, def.Name, field.Name)
			}

			validate := validateDirectives[0]
			hasValidateDirectives = true

			if rule, err := getArgumentValueAsString(validate.Arguments.ForName("rule")); err == nil {
				field.Directives = append(field.Directives, newGoTagDirective("validate", convertRule(rule)))
			} else {
				return fmt.Errorf("@%s on %s.%s requires a rule", directiveName, def.Name, field.Name)
			}

			if message, err := getArgumentValueAsString(validate.Arguments.ForName("message")); err == nil {
				field.Directives = append(field.Directives, newGoTagDirective("message", message))
			}
		}

		if hasValidateDirectives {
			p.markerTypes[typeName] = struct{}{}
		}
	}

	return nil
}

func (p *Plugin) MutateConfig(cfg *config.Config) error {
	if _, ok := cfg.Directives[goTagDirectiveName]; !ok {
		cfg.Directives[goTagDirectiveName] = config.DirectiveConfig{
			SkipRuntime: true,
		}
	}

	if _, ok := cfg.Directives[directiveName]; !ok {
		cfg.Directives[directiveName] = config.DirectiveConfig{
			SkipRuntime: true,
		}
	}

	return nil
}

func (p *Plugin) GenerateCode(cfg *codegen.Data) error {
	types := slices.Collect(maps.Keys(p.markerTypes))
	sort.Strings(types)

	filename := filepath.Join(filepath.Dir(cfg.Config.Model.Filename), "validatable_gen.go")
	if len(types) == 0 {
		_ = os.Remove(filename)

		return nil
	}

	data := struct {
		Types []string
	}{Types: types}

	return renderTemplate(templates.Options{
		PackageName:     cfg.Config.Model.Package,
		Filename:        filename,
		Template:        markersTemplate,
		Data:            data,
		Packages:        cfg.Config.Packages,
		GeneratedHeader: true,
	})
}

func convertRule(rule string) string {
	return rule
}

func getArgumentValueAsString(arg *ast.Argument) (string, error) {
	var (
		v  string
		ok bool
	)

	if k, err := arg.Value.Value(nil); err == nil {
		if v, ok = k.(string); !ok {
			return "", errors.New("argument value is not a string")
		}
	}

	if v == "" {
		return "", errors.New("argument value is an empty string")
	}

	return v, nil
}

func newGoTagDirective(key, value string) *ast.Directive {
	return &ast.Directive{
		Name: "goTag",
		Arguments: ast.ArgumentList{
			&ast.Argument{
				Name: "key",
				Value: &ast.Value{
					Raw:  key,
					Kind: ast.StringValue,
				},
			},
			&ast.Argument{
				Name: "value",
				Value: &ast.Value{
					Raw:  value,
					Kind: ast.StringValue,
				},
			},
		},
	}
}
