package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nguyen-quang-phu/go-grpc-graphql-microservice/catalog"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		_, err := catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}

		defer r.Close()
		return nil
	})
	log.Println("listening on :8080")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
