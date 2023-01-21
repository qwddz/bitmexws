package wsclient

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/qwddz/bitmexws/internal/config"
	"github.com/qwddz/bitmexws/pkg/logger"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Client struct {
	config *config.Config
	router *gin.Engine
	ws     *websocket.Upgrader
	log    logger.Logger

	receiver <-chan WSMessage
}

func New() (*Client, error) {
	c := Client{
		config: config.NewConfig(),
		router: gin.Default(),
		log:    logrus.New(),
	}

	return &c, nil
}

func (cl *Client) ServeWS(receiver <-chan WSMessage) error {
	cl.setWSUpgrader()
	cl.setReceiverChan(receiver)
	cl.setRouter()

	cl.log.Infoln("setup log: application has been successfully started")

	return cl.router.Run(cl.config.AppConfig.BindAddr)
}

func (cl *Client) setRouter() {
	cl.log.Infoln("setup log: setup router")

	cl.router.GET("/ws", NewWSHandler(cl.ws, cl.receiver, cl.log).Handle())
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

func (cl *Client) setReceiverChan(receiver <-chan WSMessage) {
	cl.receiver = receiver
}
