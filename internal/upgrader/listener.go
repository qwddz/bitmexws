package upgrader

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/qwddz/bitmexws/internal/config"
	"github.com/qwddz/bitmexws/internal/statistics"
	"github.com/qwddz/bitmexws/internal/wsclient"
	"github.com/qwddz/bitmexws/pkg/bitmex"
	"github.com/qwddz/bitmexws/pkg/logger"
	"github.com/qwddz/bitmexws/pkg/store"
	"github.com/sirupsen/logrus"
)

var (
	ErrMessageUnsupported = errors.New("message is not supported")
)

type Listener struct {
	config *config.Config
	log    logger.Logger
	store  *store.Store

	wsc Upgrader
}

func New() (*Listener, error) {
	c := Listener{
		config: config.NewConfig(),
		wsc:    bitmex.NewWSClient(),
		log:    logrus.New(),
	}

	if err := c.configureStore(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (l *Listener) Connect(ctx context.Context) error {
	if err := l.wsc.Connect(ctx, l.config.WSBitmex.URL); err != nil {
		l.log.Errorln("error when tying to open websocket connection:", err.Error())

		return err
	}

	if err := l.wsc.Subscribe(bitmex.TypeInstrument); err != nil {
		l.log.Errorln("error when tying to subscribe to websocket channel:", err.Error())

		return err
	}

	return nil
}

func (l *Listener) Listen(ctx context.Context, receiver chan wsclient.WSMessage) {
	listen := true

	for listen {
		select {
		case <-ctx.Done():
			{
				listen = false

				break
			}
		default:
			{
				msg, err := l.wsc.ReadMessage()

				if err != nil {
					listen = false

					break
				}

				decmsg, err := l.decodeMessage(msg)
				if err != nil {
					l.log.Errorln("error when trying to decode message:", err.Error())

					continue
				}

				if err := l.saveStatistics(ctx, decmsg); err != nil {
					l.log.Errorln("error when trying to saving statistics data:", err.Error())
				}

				receiver <- decmsg
			}
		}
	}
}

func (l *Listener) Shutdown() error {
	if err := l.wsc.Unsubscribe(bitmex.TypeInstrument); err != nil {
		l.log.Errorln("error when trying to unsubscribe from websocket channel:", err.Error())

		return err
	}

	if err := l.wsc.Close(); err != nil {
		l.log.Errorln("error when tying to close websocket connection:", err.Error())

		return err
	}

	return nil
}

func (l *Listener) configureStore() error {
	conf := store.Config{
		Host: store.Host{
			Master: l.config.DB.Host.Master,
			Slave:  l.config.DB.Host.Slave,
		},
		Name:     l.config.DB.Name,
		User:     l.config.DB.User,
		Password: l.config.DB.Password,
	}

	st, err := store.New(conf)
	if err != nil {
		return err
	}

	l.store = st

	return nil
}

func (l *Listener) decodeMessage(bmsg []byte) (wsclient.WSMessage, error) {
	var msg bitmex.ReceiveMessage
	var resmsg wsclient.WSMessage

	if err := json.Unmarshal(bmsg, &msg); err != nil {
		return resmsg, err
	}

	if len(msg.Data) == 0 {
		return resmsg, ErrMessageUnsupported
	}

	data := msg.Data[0]

	resmsg.Action = string(msg.Action)
	resmsg.Symbol = data.Symbol
	resmsg.MarkPrice = data.MarkPrice
	resmsg.FairPrice = data.FairPrice
	resmsg.Timestamp = data.Timestamp

	return resmsg, nil
}

func (l *Listener) saveStatistics(ctx context.Context, msg wsclient.WSMessage) error {
	if err := statistics.New(l.store.ForceMasterConnection()).Save(ctx, msg.Symbol, msg.MarkPrice); err != nil {
		return err
	}

	return nil
}
