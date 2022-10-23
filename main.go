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
	receiver := make(chan []byte)

	cs := consumer.NewConsumer()

	if err := cs.ListenReceiver(ctx, receiver); err != nil {
		close(receiver)

		log.Fatal(err)
	}

	defer func() {
		cancel()

		if err := cs.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := client.NewClient().ServeTCP(receiver); err != nil {
			close(receiver)

			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-quit
}
