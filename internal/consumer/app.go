package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/qwddz/bitmexws/internal/config"
	"github.com/qwddz/bitmexws/internal/message"
	"github.com/qwddz/bitmexws/internal/statistics"
	"github.com/qwddz/bitmexws/pkg/bitmex"
	"github.com/qwddz/bitmexws/pkg/logger"
	"github.com/qwddz/bitmexws/pkg/store"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	config *config.Config
	log    logger.Logger
	stats  *statistics.StatRepo

	wsc Upgrader
}

func NewConsumer() (*Consumer, error) {
	c := Consumer{
		config: config.NewConfig(),
		wsc:    bitmex.NewWSClient(),
		log:    logrus.New(),
	}

	if err := c.configureStore(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (cns *Consumer) ListenReceiver(ctx context.Context, receiver chan message.WSMessage) error {
	if err := cns.wsc.Connect(ctx, cns.config.WSBitmex.URL); err != nil {
		cns.log.Errorln("error when tying to open websocket connection: ", err.Error())

		return err
	}

	if err := cns.wsc.Subscribe(bitmex.TypeInstrument); err != nil {
		cns.log.Errorln("error when tying to subscribe to websocket channel: ", err.Error())

		return err
	}

	listen := true

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
					msg, err := cns.wsc.ReadMessage()

					if err != nil {
						listen = false

						break
					}

					decmsg, err := cns.decodeMessage(msg)
					if err != nil {
						continue
					}

					if decmsg.Action != string(bitmex.ActionUpdate) {
						continue
					}

					if err := cns.saveStatistics(ctx, decmsg); err != nil {
						cns.log.Errorln("error when trying to saving statistics data", err.Error())
					}

					receiver <- decmsg
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

func (cns *Consumer) configureStore() error {
	conf := store.Config{
		Host: store.Host{
			Master: cns.config.DB.Host.Master,
			Slave:  cns.config.DB.Host.Slave,
		},
		Name:     cns.config.DB.Name,
		User:     cns.config.DB.User,
		Password: cns.config.DB.Password,
	}

	st, err := store.New(conf)
	if err != nil {
		return err
	}

	if err := st.ForceMasterConnection().Ping(); err != nil {
		return err
	}

	cns.stats = statistics.NewStat(st.SlaveConnection())

	return nil
}

func (cns *Consumer) decodeMessage(bmsg []byte) (message.WSMessage, error) {
	var msg bitmex.ReceiveMessage
	var resmsg message.WSMessage

	if err := json.Unmarshal(bmsg, &msg); err != nil {
		return resmsg, err
	}

	if len(msg.Data) == 0 {
		return resmsg, errors.New("message is not supported")
	}

	data := msg.Data[0]

	resmsg.Action = string(msg.Action)
	resmsg.Symbol = data.Symbol
	resmsg.MarkPrice = data.MarkPrice
	resmsg.FairPrice = data.FairPrice
	resmsg.Timestamp = data.Timestamp

	return resmsg, nil
}

func (cns *Consumer) saveStatistics(ctx context.Context, msg message.WSMessage) error {
	if err := cns.stats.Save(ctx, msg.Symbol, msg.MarkPrice); err != nil {
		return err
	}

	return nil
}
