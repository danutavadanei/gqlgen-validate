package main

import (
	"log"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"

	"github.com/danutavadanei/gqlgen-validate/plugin"
)

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		log.Fatal(err)
	}

	if err = api.Generate(cfg, api.ReplacePlugin(plugin.New())); err != nil {
		log.Fatal(err)
	}
}
