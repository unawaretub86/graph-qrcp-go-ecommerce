package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"github.com/unawaretub86/graph-qrcp-go-ecommerce/catalog"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to process env vars: %v", err)
	}

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(attempt int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}

		return err
	})

	defer r.Close()

	s := catalog.NewService(r)
	log.Println("listening on port 8080")
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
