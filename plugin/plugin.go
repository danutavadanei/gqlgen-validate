package plugin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/vektah/gqlparser/v2/ast"
)

var renderTemplate = templates.Render

// directiveName identifies the GraphQL directive handled by this package.
const directiveName = "validate"

// ModelPlugin augments the default gqlgen model generator by injecting
// `validator:"..."` struct tags onto generated input models based on the schema
// directives.
type ModelPlugin struct {
	*modelgen.Plugin

	rules    map[string]map[string][]string
	goNames  map[string]map[string]string
	messages map[string]map[string]string
	mutate   func(*config.Config) error
}

// New constructs a gqlgen plugin that replaces the built-in
// modelgen plugin. It preserves the default behaviour while adding support for
// the @validator directive.
func New() plugin.Plugin {
	base := modelgen.New().(*modelgen.Plugin)
	p := &ModelPlugin{
		Plugin:   base,
		rules:    map[string]map[string][]string{},
		goNames:  map[string]map[string]string{},
		messages: map[string]map[string]string{},
		mutate:   base.MutateConfig,
	}
	base.FieldHook = p.fieldHook

	return p
}

func (p *ModelPlugin) Name() string {
	// Match the default model generator name so we can replace it cleanly.
	return "modelgen"
}

func (p *ModelPlugin) MutateConfig(cfg *config.Config) error {
	if cfg == nil || cfg.Schema == nil {
		return errors.New("nil schema passed to validator model plugin")
	}

	err := p.collectRules(cfg.Schema)
	if err != nil {
		return err
	}

	p.seedGoNames(cfg.Schema)

	if p.mutate != nil {
		if err = p.mutate(cfg); err != nil {
			return err
		}
	}

	err = p.renderMarkers(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (p *ModelPlugin) collectRules(schema *ast.Schema) error {
	p.rules = map[string]map[string][]string{}
	p.messages = map[string]map[string]string{}

	for typeName, def := range schema.Types {
		if def == nil || def.Kind != ast.InputObject || strings.HasPrefix(typeName, "__") {
			continue
		}

		if dir := def.Directives.ForName(directiveName); dir != nil {
			return fmt.Errorf("@%s may only be applied to input fields (found on %s)", directiveName, def.Name)
		}

		for _, field := range def.Fields {
			if field == nil || len(field.Directives) == 0 {
				continue
			}

			count := 0
			for _, dir := range field.Directives {
				if dir == nil || dir.Name != directiveName {
					continue
				}

				count++
				if count > 1 {
					return fmt.Errorf("@%s may only be applied once per field (%s.%s)", directiveName, def.Name, field.Name)
				}

				rule := argumentValue(dir, "rule")
				if rule == "" {
					return fmt.Errorf("@%s on %s.%s requires a rule", directiveName, def.Name, field.Name)
				}

				message := argumentValue(dir, "message")

				for _, arg := range dir.Arguments {
					if arg.Name != "rule" && arg.Name != "message" {
						return fmt.Errorf("@%s only supports the 'rule' and 'message' arguments (unexpected %q on %s.%s)", directiveName, arg.Name, def.Name, field.Name)
					}
				}

				addRule(p.rules, def.Name, field.Name, rule)
				if message != "" {
					addMessage(p.messages, def.Name, field.Name, message)
				}
			}
		}
	}

	if err := ensureNoArgumentRules(schema.Mutation); err != nil {
		return err
	}
	if err := ensureNoArgumentRules(schema.Query); err != nil {
		return err
	}
	if err := ensureNoArgumentRules(schema.Subscription); err != nil {
		return err
	}

	return nil
}

func (p *ModelPlugin) seedGoNames(schema *ast.Schema) {
	p.goNames = map[string]map[string]string{}

	for typeName, def := range schema.Types {
		if def == nil || def.Kind != ast.InputObject || strings.HasPrefix(typeName, "__") {
			continue
		}

		mapping := make(map[string]string, len(def.Fields))
		for _, field := range def.Fields {
			mapping[field.Name] = templates.ToGo(field.Name)
		}

		p.goNames[typeName] = mapping
	}
}

func (p *ModelPlugin) fieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	mf, err := modelgen.DefaultFieldMutateHook(td, fd, f)
	if err != nil || mf == nil {
		return mf, err
	}

	if _, ok := p.goNames[td.Name]; !ok {
		p.goNames[td.Name] = map[string]string{}
	}

	p.goNames[td.Name][fd.Name] = mf.GoName

	rules := p.rulesForField(td.Name, fd.Name)
	if len(rules) == 0 {
		return mf, nil
	}

	converted := make([]string, 0, len(rules))
	for _, rule := range rules {
		converted = append(converted, p.convertRule(td.Name, rule))
	}

	mf.Tag = setStructTag(mf.Tag, "validate", strings.Join(converted, ","))

	if msg := p.messageForField(td.Name, fd.Name); msg != "" {
		mf.Tag = setStructTag(mf.Tag, "message", msg)
	}

	return mf, nil
}

func (p *ModelPlugin) rulesForField(typeName, fieldName string) []string {
	fields := p.rules[typeName]
	if len(fields) == 0 {
		return nil
	}

	rules := fields[fieldName]
	if len(rules) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(rules))

	out := make([]string, 0, len(rules))
	for _, rule := range rules {
		trimmed := strings.TrimSpace(rule)
		if trimmed == "" {
			continue
		}

		if _, ok := seen[trimmed]; ok {
			continue
		}

		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}

	sort.Strings(out)

	return out
}

func (p *ModelPlugin) convertRule(typeName, rule string) string {
	mapping := p.goNames[typeName]
	if len(mapping) == 0 {
		return rule
	}

	var (
		builder strings.Builder
		token   strings.Builder
	)

	convertFollowing := false
	convertToken := false

	flush := func() {
		if token.Len() == 0 {
			return
		}

		word := token.String()
		if convertToken {
			if replacement, ok := mapping[word]; ok {
				builder.WriteString(replacement)
				token.Reset()

				convertToken = false

				return
			}
		}

		builder.WriteString(word)
		token.Reset()

		convertToken = false
	}

	for _, r := range rule {
		if isIdentifierRune(r) {
			if token.Len() == 0 {
				convertToken = convertFollowing
			}

			token.WriteRune(r)
		} else {
			flush()

			switch r {
			case '=', ':':
				convertFollowing = true
			case ',', ' ', '\t':
				// continue converting additional parameters in a rule list
				// without resetting the mode
			default:
				convertFollowing = false
			}

			builder.WriteRune(r)
		}
	}

	flush()

	return builder.String()
}

func addRule(target map[string]map[string][]string, typeName, fieldName, rule string) {
	if _, ok := target[typeName]; !ok {
		target[typeName] = map[string][]string{}
	}

	target[typeName][fieldName] = append(target[typeName][fieldName], strings.TrimSpace(rule))
}

func addMessage(target map[string]map[string]string, typeName, fieldName, message string) {
	if _, ok := target[typeName]; !ok {
		target[typeName] = map[string]string{}
	}

	target[typeName][fieldName] = strings.TrimSpace(message)
}

func (p *ModelPlugin) messageForField(typeName, fieldName string) string {
	fields := p.messages[typeName]
	if len(fields) == 0 {
		return ""
	}

	return fields[fieldName]
}

func argumentValue(dir *ast.Directive, name string) string {
	if dir == nil {
		return ""
	}

	arg := dir.Arguments.ForName(name)
	if arg == nil || arg.Value == nil {
		return ""
	}

	if arg.Value.Raw != "" {
		return arg.Value.Raw
	}

	return arg.Value.String()
}

func isIdentifierRune(r rune) bool {
	switch {
	case r == '_', r == '.', r == '-', r == '$':
		return true
	case r >= '0' && r <= '9':
		return true
	case r >= 'a' && r <= 'z':
		return true
	case r >= 'A' && r <= 'Z':
		return true
	default:
		return false
	}
}

type tagPart struct {
	key   string
	value string
}

func setStructTag(tagString, key, value string) string {
	parts := parseStructTag(tagString)
	updated := false

	for i := range parts {
		if parts[i].key == key {
			parts[i].value = value
			updated = true

			break
		}
	}

	if !updated {
		parts = append(parts, tagPart{key: key, value: value})
	}

	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		if part.key == "" {
			continue
		}

		segments = append(segments, fmt.Sprintf(`%s:"%s"`, part.key, part.value))
	}

	return strings.TrimSpace(strings.Join(segments, " "))
}

const markerTemplate = `{{- if .Types }}
{{ range .Types }}
func ({{ . }}) IsValidatable() {}
{{ end }}
{{- end }}`

func (p *ModelPlugin) renderMarkers(cfg *config.Config) error {
	types := make([]string, 0)

	for typeName, fieldRules := range p.rules {
		if len(fieldRules) == 0 {
			continue
		}

		types = append(types, templates.ToGo(typeName))
	}

	sort.Strings(types)

	filename := filepath.Join(filepath.Dir(cfg.Model.Filename), "validatable_gen.go")
	if len(types) == 0 {
		_ = os.Remove(filename)

		return nil
	}

	data := struct {
		Types []string
	}{Types: types}

	return renderTemplate(templates.Options{
		PackageName:     cfg.Model.Package,
		Filename:        filename,
		Template:        markerTemplate,
		Data:            data,
		Packages:        cfg.Packages,
		GeneratedHeader: true,
	})
}

func parseStructTag(tag string) []tagPart {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}

	rawParts := splitTagParts(tag)

	parts := make([]tagPart, 0, len(rawParts))
	for _, raw := range rawParts {
		if raw == "" {
			continue
		}

		idx := strings.Index(raw, ":")
		if idx <= 0 {
			continue
		}

		key := strings.TrimSpace(raw[:idx])
		value := strings.Trim(raw[idx+1:], "\"")
		parts = append(parts, tagPart{key: key, value: value})
	}

	return parts
}

func splitTagParts(tag string) []string {
	var (
		parts []string
		buf   strings.Builder
		inStr bool
	)

	for _, r := range tag {
		switch r {
		case '"':
			inStr = !inStr

			buf.WriteRune(r)
		case ' ':
			if inStr {
				buf.WriteRune(r)
			} else if buf.Len() > 0 {
				parts = append(parts, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteRune(r)
		}
	}

	if buf.Len() > 0 {
		parts = append(parts, buf.String())
	}

	return parts
}

func ensureNoArgumentRules(def *ast.Definition) error {
	if def == nil {
		return nil
	}
	for _, field := range def.Fields {
		if field == nil {
			continue
		}
		for _, arg := range field.Arguments {
			if arg == nil {
				continue
			}
			if arg.Directives.ForName(directiveName) != nil {
				return fmt.Errorf("@%s is not supported on argument %s.%s(%s)", directiveName, def.Name, field.Name, arg.Name)
			}
		}
	}
	return nil
}
