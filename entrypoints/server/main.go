package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegodesousas/clean-boilerplate-go/infra/monitor"

	"cloud.google.com/go/pubsub"
	"github.com/diegodesousas/clean-boilerplate-go/infra/config"
	"github.com/diegodesousas/clean-boilerplate-go/infra/database"
	"github.com/diegodesousas/clean-boilerplate-go/infra/http/handlers/healthcheck"
	"github.com/diegodesousas/clean-boilerplate-go/infra/http/server"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	config.Load()

	if err := database.InitPool(); err != nil {
		log.Fatal(err)
	}

	nrApp, err := monitor.New()
	if err != nil {
		log.Fatal(err)
	}

	// instantiate app dependencies

	// pubsub usage example
	pubsubClient, err := pubsub.NewClient(
		context.Background(),
		viper.GetString("PUBSUB_PROJECT_ID"),
		option.WithCredentialsFile(viper.GetString("PUBSUB_SERVICE_KEY_PATH")),
	)
	if err != nil {
		log.Printf("Could not connect to gcloud pub / sub; err: %s", err)
	}

	s := server.NewServer(
		server.WithNewRelicWrapper(nrApp),
		healthcheck.Routes(database.Pool()),
	)

	go func() {
		log.Printf("http server up and listen on port %s", viper.GetString("HTTP_PORT"))
		if err := s.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		case <-interrupt:
			log.Println("shutdown application")
			log.Println("shutdown http server")
			if err := s.Shutdown(context.Background()); err != nil {
				log.Printf("http: %s", err)
			}

			log.Println("close database connections")
			if err := database.ClosePool(); err != nil {
				log.Printf("http: %s", err)
			}

			log.Println("close pubsub connections")
			if err := pubsubClient.Close(); err != nil {
				log.Printf("queue: %s", err)
			}

			return
		}
	}
}
