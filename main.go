package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
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
	ctx := context.Background()

	ctxWithCancel, cancel := context.WithCancel(ctx)

	app := consumer.NewConsumer()

	defer func() {
		cancel()

		if err := app.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	receiver := app.ListenReceiver(ctxWithCancel)

	go func() {
		var listen = true

		for listen {
			select {
			case <-ctxWithCancel.Done():
				{
					listen = false
				}
			case msg, closed := <-receiver:
				{
					if closed == false {
						listen = false
					}

					fmt.Println(msg.Data)
				}
			}
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
}
