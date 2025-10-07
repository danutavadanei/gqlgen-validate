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
	"strings"

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

//go:embed markers.gotpl
var markersTemplate string

var (
	_ plugin.CodeGenerator = &Plugin{}
	_ plugin.ConfigMutator = &Plugin{}
	_ plugin.SchemaMutator = &Plugin{}
)

var (
	crossFieldRules = set{
		"eqfield": {}, "nefield": {}, "gtfield": {}, "gtefield": {}, "ltfield": {}, "ltefield": {},
	}

	crossFieldRelativeRules = set{
		"eqcsfield": {}, "necsfield": {}, "gtcsfield": {}, "gtecsfield": {}, "ltcsfield": {}, "ltecsfield": {},
		"eqsfield": {}, "nesfield": {}, "gtsfield": {}, "gtesfield": {}, "ltsfield": {}, "ltesfield": {},
		"fieldcontains": {}, "fieldexcludes": {}, "containsfield": {}, "excludesfield": {},
	}

	multiFieldRules = set{
		"required_with": {}, "required_with_all": {}, "required_without": {}, "required_without_all": {},
		"excluded_with": {}, "excluded_with_all": {}, "excluded_without": {}, "excluded_without_all": {},
	}

	pairedFieldRules = set{
		"required_if": {}, "required_unless": {}, "excluded_if": {}, "excluded_unless": {},
		"skip_unless": {},
	}
)

var toGo = templates.ToGo

// set is a simple string set.
type set map[string]struct{}

func (s set) add(name string) {
	s[name] = struct{}{}
}

func (s set) contains(name string) bool {
	_, ok := s[name]
	return ok
}

func (s set) values() []string {
	return slices.Collect(maps.Keys(s))
}

// Plugin is a gqlgen plugin that wires validation rules into generated models.
type Plugin struct {
	markerTypes set
}

// New constructs the plugin instance.
func New() plugin.Plugin {
	return &Plugin{
		markerTypes: make(set),
	}
}

// Name implements plugin.Plugin.
func (p *Plugin) Name() string { return "gqlgen-validate" }

// MutateSchema ensures directives exist and rewrites fields with validation metadata.
func (p *Plugin) MutateSchema(schema *ast.Schema) error {
	// Ensure goTag directive exists (used to inject struct tags).
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

			rule, err := getArgumentValueAsString(validate.Arguments.ForName("rule"))
			if err != nil {
				return fmt.Errorf("@%s on %s.%s requires a rule", directiveName, def.Name, field.Name)
			}

			field.Directives = append(field.Directives, newGoTagDirective("validate", toGoRuleParams(rule)))

			if message, err := getArgumentValueAsString(validate.Arguments.ForName("message")); err == nil {
				field.Directives = append(field.Directives, newGoTagDirective("message", message))
			}
		}

		if hasValidateDirectives {
			p.markerTypes.add(typeName)
		}
	}

	return nil
}

// MutateConfig registers the directives so gqlgen does not expect runtime handlers.
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

// GenerateCode emits a small file that marks the validated input types.
func (p *Plugin) GenerateCode(cfg *codegen.Data) error {
	types := p.markerTypes.values()
	sort.Strings(types)

	filename := filepath.Join(filepath.Dir(cfg.Config.Model.Filename), "validatable_gen.go")
	if len(types) == 0 {
		_ = os.Remove(filename)
		return nil
	}

	data := struct {
		Types []string
	}{Types: types}

	return templates.Render(templates.Options{
		PackageName:     cfg.Config.Model.Package,
		Filename:        filename,
		Template:        markersTemplate,
		Data:            data,
		Packages:        cfg.Config.Packages,
		GeneratedHeader: true,
	})
}

func getArgumentValueAsString(arg *ast.Argument) (string, error) {
	var (
		v  string
		ok bool
	)

	if arg == nil {
		return "", errors.New("argument is nil")
	}

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

func toGoRuleParams(rule string) string {
	var out strings.Builder
	i := 0
	n := len(rule)

	for i < n {
		start := i
		for i < n && rule[i] != ',' && rule[i] != '|' {
			i++
		}

		segment := strings.TrimSpace(rule[start:i])
		if segment != "" {
			if eq := strings.IndexByte(segment, '='); eq >= 0 {
				name := segment[:eq]
				params := segment[eq+1:]
				segment = name + "=" + transformRuleParams(name, params)
			}
			out.WriteString(segment)
		}

		if i < n {
			out.WriteByte(rule[i])
			i++
		}
	}

	return out.String()
}

func transformRuleParams(name, params string) string {
	switch {
	case crossFieldRules.contains(name):
		return toGo(params)
	case crossFieldRelativeRules.contains(name):
		return toGoBySeparator(params, ".")
	case multiFieldRules.contains(name):
		return toGoBySeparator(params, " ")
	case pairedFieldRules.contains(name):
		return toGoPairs(params)
	default:
		return params
	}
}

func toGoBySeparator(value, separator string) string {
	segments := strings.Split(value, separator)
	for i, segment := range segments {
		segments[i] = toGo(segment)
	}
	return strings.Join(segments, separator)
}

func toGoPairs(value string) string {
	fields := strings.Fields(value)
	if len(fields)%2 != 0 {
		return value
	}

	for i := 0; i < len(fields); i += 2 {
		fields[i] = toGo(fields[i])
	}
	return strings.Join(fields, " ")
}
