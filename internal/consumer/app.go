package consumer

import (
	"context"
	"encoding/json"
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

func (cns *Consumer) ListenReceiver(ctx context.Context) <-chan bitmex.ReceiveMessage {
	cns.bmWS.Connect(ctx, cns.config.WSBitmex.URL)

	if err := cns.bmWS.Subscribe(bitmex.TypeInstrument); err != nil {
		log.Fatal(err)
	}

	receiver := make(chan []byte)

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
						log.Printf("receive message error: %s", err.Error())
					}
				}
			}
		}

		close(receiver)
	}()

	upgrader := make(chan bitmex.ReceiveMessage)

	go func() {
		var listen = true

		for listen {
			select {
			case <-ctx.Done():
				{
					listen = false

					break
				}
			case msg, closed := <-receiver:
				{
					if closed == false {
						listen = false

						break
					}

					var message bitmex.ReceiveMessage
					if err := json.Unmarshal(msg, &message); err != nil {
						log.Printf("reading message error: %s", err.Error())
					}

					if message.Action != bitmex.ActionUpdate {
						continue
					}

					upgrader <- message
				}
			}
		}

		close(upgrader)
	}()

	return upgrader
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