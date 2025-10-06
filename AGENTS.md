---
description: Core project guidelines for the gqlgen-validate codebase. Apply these rules when working on any code, documentation, or configuration files within the project.
alwaysApply: true
inclusion: always
---

# gqlgen-validate Structure and Overview

This document helps AI coding assistants reason about the gqlgen-validate project.

## Project Overview

gqlgen-validate ships two tightly coupled pieces:
- a gqlgen plugin that injects `validate:"…"` and `message:"…"` tags into generated input structs according to `@validate` directives in the schema.
- a runtime directive plus middleware that run `go-playground/validator` against resolver arguments and translate validation failures into GraphQL-friendly errors.

The repository targets gqlgen. Tests cover the plugin, runtime behaviour, and reflective helpers that bridge GraphQL inputs with validator expectations.

## Directory Structure

```
gqlgen-validate/
├── gen/                 # Code generation plugin; handles @validate rules and struct tag updates.
├── runtime/             # Directive + middleware wiring; hosts validation logic and error translation utilities.
├── example/             # Minimal gqlgen app proving end-to-end integration of the plugin and middleware.
│   ├── cmd/gqlgen/      # go run entrypoint that loads gqlgen config and injects the custom plugin.
│   ├── graph/           # Generated schema, models with injected tags, and resolvers for the demo server.
│   └── server.go        # gqlgen server setup registering the directive handler and middleware.
├── README.md            # High-level usage, integration steps, and design notes.
```

## Core Components

**Plugin (`plugin`)**
- Implements the `CodeGenerator`, `ConfigMutator` and `SchemaMutator` gqlgen interfaces. It hooks into gqlgen's codegen lifecycle to modify the schema and generated models.
- Accepts `@validate(rule: String!, message: String)` on input object fields only; rejects directive misuse (arguments, object-level application, duplicates).
- Normalises rule strings and maps GraphQL field casing onto generated Go fields.
- Persists optional custom messages alongside validator rules and emits `validatable_gen.go` marker interfaces for rule-bearing types, pruning the file when no rules remain.

**Runtime (`runtime`)**
- `Middleware()` walks resolver arguments, validating any value implementing `Validatable`, dereferencing pointers and traversing collections.
- Converts validator errors into `gqlerror.Error` with `BAD_USER_INPUT`, preserving GraphQL path segments, indices, and custom `message:"…"` tags.

**Example Application (`example`)**
- Demonstrates schema usage (`graph/schema.gql`) with array validation (`dive`) and cross-field rules.
- Generated models (`graph/model/models_gen.go`) expose injected tags; `validatable_gen.go` carries `IsValidatable` helpers consumed by the middleware.
- `server.go` registers the directive and middleware; provides an end-to-end harness for manual testing.

## Documentation Resources
- Consult `README.md` for integration steps, schema directives, and configuration hints.
- The example app's `README.md` walks through regenerating code via the wrapper in `example/cmd/gqlgen`.

## Coding Guidelines
- Write idiomatic Go that compiles with Go 1.22+; always run `gofmt`/`goimports` on modified files.
- Maintain clear boundaries between code generation logic (`gen`) and runtime behaviour (`runtime`); avoid cross-package dependencies beyond shared interfaces (`Validatable`).
- Uphold existing testing patterns: table-driven tests with descriptive names, helper builders when constructing schemas or validation payloads.
- Document non-obvious reflection or tag-manipulation logic with concise comments; prefer self-explanatory code otherwise.
- Keep generated artifacts (`*_gen.go`) untouched except through the code generation flow; modify templates or plugin logic instead.

## Testing & Tooling
- Execute the full test suite before finishing changes:
  ```bash
  go test ./...
  ```
- Lint with golangci-lint to enforce style and catch static issues:
  ```bash
  golangci-lint run
  ```
- When sandboxed environments block Go's default cache path, override it (`GOCACHE=$(pwd)/.gocache go test ./...`). Clean up temporary caches afterwards.
- Add or update unit tests alongside functional changes in `gen` or `runtime`. Generated files under `example/graph` should be regenerated via the wrapper when schema changes.

## Git & Workflow
- Work directly on the current branch; do not create new branches from this repository snapshot.
- Do not commit or amend history-leave commits to the repository maintainers.
- Preserve any user-authored changes already present in the worktree; never revert unrelated modifications.
- Before handoff, ensure tests and linting have been executed or clearly state why they were skipped.
