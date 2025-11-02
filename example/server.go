package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server listen: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("listening on http://localhost:8080")

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("http server shutdown: %v", err)
	}

	log.Println("server shut down gracefully")
}
