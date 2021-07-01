package main

import (
	"context"
	"log"

	"github.com/diegodesousas/clean-boilerplate-go/infra/newrelic"

	"github.com/diegodesousas/clean-boilerplate-go/infra/config"
	"github.com/diegodesousas/clean-boilerplate-go/infra/database"
	"github.com/diegodesousas/clean-boilerplate-go/infra/http/handlers/healthcheck"
	"github.com/diegodesousas/clean-boilerplate-go/infra/http/server"
	"github.com/spf13/viper"
)

func main() {
	config.Load()

	if err := database.InitPool(); err != nil {
		log.Fatal(err)
	}

	// instantiate app dependencies

	s := server.NewServer(
		server.WithMonitorWrapper(newrelic.NewMonitorWrapper()),
		healthcheck.Routes(database.Pool()),
	)

	closeDbHandler := func(ctx context.Context) error {
		log.Println("close database connections")
		return database.ClosePool()
	}

	log.Printf("http server up and listen on port %s", viper.GetString("HTTP_PORT"))
	err := s.ListenAndServe(context.Background(), closeDbHandler)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("http server finished")
}
