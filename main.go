package main

import (
	"context"
	"github.com/qwddz/bitmexws/internal/client"
	"github.com/qwddz/bitmexws/internal/consumer"
	"github.com/qwddz/bitmexws/internal/message"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	receiver := make(chan message.WSMessage)

	cs, err := consumer.NewConsumer()
	if err != nil {
		log.Fatalln(err)
	}

	if err := cs.ListenReceiver(ctx, receiver); err != nil {
		log.Fatalln(err)
	}

	defer func() {
		close(receiver)
		cancel()

		if err := cs.Shutdown(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		cl, err := client.NewClient()
		if err != nil {
			log.Fatalln(err)
		}

		if err := cl.ServeTCP(receiver); err != nil {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-quit
}
