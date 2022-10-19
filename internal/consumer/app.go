package consumer

import (
	"context"
	"github.com/qwddz/bitmexws/internal/config"
	"github.com/qwddz/bitmexws/pkg/bitmex"
)

type Consumer struct {
	config *config.Config

	wsc Upgrader
}

func NewConsumer() *Consumer {
	return &Consumer{
		config: config.NewConfig(),
		wsc:    bitmex.NewWSClient(),
	}
}

func (cns *Consumer) ListenReceiver(ctx context.Context, receiver chan []byte) error {
	if err := cns.wsc.Connect(ctx, cns.config.WSBitmex.URL); err != nil {
		return err
	}

	if err := cns.wsc.Subscribe(bitmex.TypeInstrument); err != nil {
		return err
	}

	go func() {
		var listen = true

		for listen {
			select {
			case <-ctx.Done():
				{
					listen = false

					break
				}
			default:
				{
					err := cns.wsc.ReadMessage(receiver)

					if err != nil {
						listen = false

						break
					}
				}
			}
		}
	}()

	return nil
}

func (cns *Consumer) Shutdown() error {
	if err := cns.wsc.Unsubscribe(bitmex.TypeInstrument); err != nil {
		return err
	}

	if err := cns.wsc.Close(); err != nil {
		return err
	}

	return nil
}
