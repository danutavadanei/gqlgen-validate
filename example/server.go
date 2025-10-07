package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/danutavadanei/gqlgen-validate/example/graph"
	"github.com/danutavadanei/gqlgen-validate/example/graph/generated"
	"github.com/danutavadanei/gqlgen-validate/runtime"
)

func main() {
	resolver := &graph.Resolver{}
	cfg := generated.Config{Resolvers: resolver}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))
	srv.AroundFields(runtime.Middleware())

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
