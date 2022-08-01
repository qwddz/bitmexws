package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/qwddz/bitmexws/internal/client"
	"github.com/qwddz/bitmexws/internal/consumer"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no dot env file")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	app := consumer.NewConsumer()

	receiver := app.ListenReceiver(ctx)

	defer func() {
		cancel()

		if err := app.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := client.NewClient().ServeHTTP(receiver); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
}
