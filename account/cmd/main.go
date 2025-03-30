package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"

	account "github.com/unawaretub86/graph-qrcp-go-ecommerce/account"
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

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(attempt int) (err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}

		return err
	})

	defer r.Close()
	log.Println("listening on port 8080")

	s := account.NewAccountService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
