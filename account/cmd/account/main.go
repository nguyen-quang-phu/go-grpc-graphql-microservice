package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nguyen-quang-phu/go-grpc-graphql-microservice/account"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		_, err := account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}

		defer r.Close()
		return nil
	})
	log.Println("listening on :8080")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
