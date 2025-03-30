package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CALALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	config := AppConfig{}

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	graphqlServer, err := NewGraphQLServer(config.AccountURL, config.CatalogURL, config.OrderURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	http.Handle("/graphql", handler.New(graphqlServer.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("playground", "/playground"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
