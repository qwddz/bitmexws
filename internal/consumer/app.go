package consumer

import (
	"context"
	"github.com/qwddz/bitmexws/internal/config"
	"github.com/qwddz/bitmexws/pkg/bitmex"
	"github.com/qwddz/bitmexws/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	config *config.Config
	log    logger.Logger

	wsc Upgrader
}

func NewConsumer() *Consumer {
	return &Consumer{
		config: config.NewConfig(),
		wsc:    bitmex.NewWSClient(),
		log:    logrus.New(),
	}
}

func (cns *Consumer) ListenReceiver(ctx context.Context, receiver chan []byte) error {
	if err := cns.wsc.Connect(ctx, cns.config.WSBitmex.URL); err != nil {
		cns.log.Errorln("error when tying to open websocket connection: ", err.Error())

		return err
	}

	if err := cns.wsc.Subscribe(bitmex.TypeInstrument); err != nil {
		cns.log.Errorln("error when tying to subscribe to websocket channel: ", err.Error())

		return err
	}

	var listen = true

	go func() {
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
		cns.log.Errorln("error when trying to unsubscribe from websocket channel", err.Error())

		return err
	}

	if err := cns.wsc.Close(); err != nil {
		cns.log.Errorln("error when tying to close websocket connection: ", err.Error())

		return err
	}

	return nil
}

func (cns *Consumer) setLogger() {
	cns.log = logrus.New()
}
