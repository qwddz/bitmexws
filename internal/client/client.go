package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/qwddz/bitmexws/internal/config"
	"github.com/qwddz/bitmexws/internal/message"
	"github.com/qwddz/bitmexws/internal/statistics"
	"github.com/qwddz/bitmexws/pkg/logger"
	"github.com/qwddz/bitmexws/pkg/store"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Client struct {
	config *config.Config
	router *gin.Engine
	ws     *websocket.Upgrader
	log    logger.Logger
	stats  *statistics.StatRepo

	receiver <-chan message.WSMessage
}

func NewClient() (*Client, error) {
	c := Client{
		config: config.NewConfig(),
		router: gin.Default(),
		log:    logrus.New(),
	}

	if err := c.configureStore(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (cl *Client) ServeTCP(receiver <-chan message.WSMessage) error {
	cl.setWSUpgrader()
	cl.setReceiverChan(receiver)
	cl.setRouter()

	cl.log.Infoln("setup log: application has been successfully started")

	return cl.router.Run(cl.config.AppConfig.BindAddr)
}

func (cl *Client) setRouter() {
	cl.log.Infoln("setup log: setup router")

	cl.router.GET("/ws", NewWSHandler(cl.ws, cl.receiver, cl.log).Handle())
	cl.router.GET("/statistics", NewStatsHandler(cl.stats).Handle())
}

func (cl *Client) setWSUpgrader() {
	cl.log.Infoln("setup log: setup websocket upgrader")

	var u = websocket.Upgrader{}

	if cl.config.AppConfig.Debug {
		u.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	cl.ws = &u
}

func (cl *Client) setReceiverChan(receiver <-chan message.WSMessage) {
	cl.receiver = receiver
}

func (cl *Client) configureStore() error {
	cl.log.Infoln("setup log: setup client db connection")

	st, err := store.New(cl.config)
	if err != nil {
		return err
	}

	cl.log.Infoln("setup log: ping client db connection")

	if err := st.ForceMasterConnection().Ping(); err != nil {
		return err
	}

	cl.stats = statistics.NewStat(st.ForceMasterConnection())

	return nil
}
