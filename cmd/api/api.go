package main

import (
	"context"
	"github.com/qwddz/bitmexws/internal/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	a, err := api.New()
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		if err := a.ServeTCP(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API start error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API server...")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Shutdown(); err != nil {
		log.Fatal("API server forced to shutdown:", err)
	}

	log.Println("API server exiting...")
}
