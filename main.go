package main

import (
	"context"
	"github.com/qwddz/bitmexws/internal/client"
	"github.com/qwddz/bitmexws/internal/consumer"
	"log"
	"os"
	"os/signal"
	"syscall"
)

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
		if err := client.NewClient().ServeTCP(receiver); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-quit
}
