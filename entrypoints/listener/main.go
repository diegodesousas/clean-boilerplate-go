package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
	"github.com/diegodesousas/clean-boilerplate-go/infra/config"
	"github.com/diegodesousas/clean-boilerplate-go/infra/listeners"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func main() {
	log.Println("start listener")

	config.Load()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()

	client, err := pubsub.NewClient(
		ctx,
		viper.GetString("PUBSUB_PROJECT_ID"),
		option.WithCredentialsFile(viper.GetString("PUBSUB_SERVICE_KEY_PATH")),
	)
	if err != nil {
		log.Fatal(err)
	}

	chanErrors := make(chan error, 1)

	for _, listener := range []listeners.Listener{} {
		listener := listener

		go func() {
			err := client.
				Subscription(listener.Subscription()).
				Receive(ctx, listener.Receive)
			if err != nil {
				chanErrors <- err
			}
		}()
	}

	for {
		select {
		case err := <-chanErrors:
			log.Println(err)
			interrupt <- syscall.SIGTERM

		case <-interrupt:
			log.Println("shutdown listener")
			err := client.Close()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
}
