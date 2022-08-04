package consumer

import (
	"context"
	"github.com/qwddz/bitmexws/internal/config"
	"github.com/qwddz/bitmexws/pkg/bitmex"
	"log"
)

type Consumer struct {
	config *config.Config

	bmWS Upgrader
}

func NewConsumer() *Consumer {
	return &Consumer{
		config: config.NewConfig(),
		bmWS:   bitmex.NewWSClient(),
	}
}

func (cns *Consumer) ListenReceiver(ctx context.Context, receiver chan []byte) {
	cns.bmWS.Connect(ctx, cns.config.WSBitmex.URL)

	if err := cns.bmWS.Subscribe(bitmex.TypeInstrument); err != nil {
		log.Fatal(err)
	}

	go func() {
		var listen = true

		for listen {
			select {
			case <-ctx.Done():
				{
					listen = false

					log.Println("receive message error: existing from app")
					break
				}
			default:
				{
					err := cns.bmWS.ReadMessage(receiver)

					if err != nil {
						listen = false

						log.Printf("receive message error: %s", err.Error())
						break
					}
				}
			}
		}
	}()
}

func (cns *Consumer) Shutdown() error {
	log.Println("unsubscribing from binmex...")

	if err := cns.bmWS.Unsubscribe(bitmex.TypeInstrument); err != nil {
		return err
	}

	log.Println("close connection from binmex...")

	if err := cns.bmWS.Close(); err != nil {
		return err
	}

	return nil
}
