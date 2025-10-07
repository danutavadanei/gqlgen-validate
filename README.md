# gqlgen-validate 

`gqlgen-validate` adds first-class validation support for gqlgen projects by
teaching the code generator about a `@validate` directive. The plugin injects
`validate:"..."` struct tags into the generated Go models so that you can rely on
[`go-playground/validator`](https://github.com/go-playground/validator) without manually editing generated code.

## Schema usage

Declare the directive once in your schema:

```graphql
"""Input validation directive (e.g., @validate(rule: "required,min=8"))."""
directive @validate(rule: String!, message: String) on INPUT_FIELD_DEFINITION | ARGUMENT_DEFINITION
```

Attach rules to specific input fields. Field names in the `rule` string
use the GraphQL casing – the plugin automatically maps them to the Go
struct field names produced by gqlgen.

Place `@validate` directly on the input fields (nested fields are fine). Set the
mandatory `rule` argument to any `go-playground/validator` expression and, when
needed, add a custom `message` to override the default runtime error text. Field
names in the `rule` string use the GraphQL casing - the plugin automatically
maps them to the Go struct field names produced by gqlgen.

During generation the plugin drops a tiny `IsValidatable` method next to every
validated input type. The runtime middleware only inspects values that
implement the `Validatable` interface, so unrelated arguments are never
touched.

```graphql
input AccountMetadata {
  bic: String @validate(rule: "len=11", message: "BIC must be exactly 11 chars")
  iban: String @validate(rule: "required_without=bic len=24")
}
```

The rules above result in generated structs that look like:

```go
type AccountMetadata struct {
    Bic  *string `json:"bic" validate:"len=11" message:"BIC must be exactly 11 chars"`
    Iban *string `json:"iban" validate:"len=24,required_without=Bic"`
}
```

## Integrating with gqlgen

To use a plugin during code generation, you need to create a new entry point.
Please refer to the [example generator](/example/cmd/gqlgen/main.go) for implementation details.

Running `go run cmd/gqlgen` will now inject the  appropriate `validate:"..."`
tags wherever your schema uses `@validate`.

## Runtime validation helper

Once the models carry validation tags you just need to wire up the runtime middleware.

```go
package main

import (
    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/danutavadanei/gqlgen-validate/runtime"
	
    "github.com/your/app/graph"
    "github.com/your/app/graph/generated"
)

func main() {
	resolver := &graph.Resolver{}
	cfg := generated.Config{Resolvers: resolver}
	
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))
	srv.AroundFields(runtime.Middleware())
}
```

The field middleware ensures every resolver argument is validated before your
business logic runs. Custom `message` tags added by the plugin automatically
override the default validator error text. The runtime middleware returns
GraphQL errors that point at the offending fields (e.g. `input.bic`).

## Example project

A runnable gqlgen server that uses the plugin lives in [example](/example)

## Design considerations

- **Middleware-first validation:** running validation in `AroundFields`
  guarantees it executes after gqlgen unmarshals inputs and before business
  logic runs. This yields consistent error formatting, avoids per-resolver
  boilerplate, and keeps validation isolated from transport-specific code.
- **Current limitations:** middleware triggers only for arguments that
  implement the generated `Validatable` marker. Scalars or primitives still need
  resolver-level checks.

## Open Questions

These are areas I am still exploring - not final decisions or guaranteed features.  

1. Add configuration knobs for global validator options (e.g., custom
   `RegisterValidation` hooks, locale-aware tag-name functions).
2. Improve error reporting with optional translation layers and richer
   extensions payloads.
3. Support additional schema shapes such as interface inputs or directive-level
   opt-outs without breaking existing tags.

Have other ideas you would like to see? Open an issue or a PR, and I'm happy to discuss!
