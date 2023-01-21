package main

import (
	"context"
	"github.com/qwddz/bitmexws/internal/upgrader"
	"github.com/qwddz/bitmexws/internal/wsclient"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	receiver := make(chan wsclient.WSMessage)

	cs, err := upgrader.New()
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
		cl, err := wsclient.New()
		if err != nil {
			log.Fatalln(err)
		}

		if err := cl.ServeWS(receiver); err != nil {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-quit
}
