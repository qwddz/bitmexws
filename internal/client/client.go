package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/qwddz/bitmexws/internal/config"
	"net/http"
)

type Client struct {
	config *config.Config
	router *gin.Engine
	ws     *websocket.Upgrader

	receiver <-chan []byte

	handler *Handler
}

func NewClient() *Client {
	return &Client{
		config: config.NewConfig(),
		router: gin.Default(),
	}
}

func (cl *Client) ServeTCP(receiver <-chan []byte) error {
	cl.setWSUpgrader()
	cl.setReceiverChan(receiver)
	cl.setRouter()

	return cl.router.Run(cl.config.AppConfig.BindAddr)
}

func (cl *Client) setRouter() {
	cl.router.GET("/ws", NewHandler(cl.ws, cl.receiver).Handle())
}

func (cl *Client) setWSUpgrader() {
	var u = websocket.Upgrader{}

	if cl.config.AppConfig.Debug {
		u.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	cl.ws = &u
}

func (cl *Client) setReceiverChan(receiver <-chan []byte) {
	cl.receiver = receiver
}
