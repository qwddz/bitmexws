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

func (cl *Client) ServeHTTP(receiver <-chan []byte) error {
	cl.configureWS()
	cl.configureReceiverChan(receiver)
	cl.configureHandlers()

	cl.configureRouter()

	return cl.router.Run(cl.config.AppConfig.BindAddr)
}

func (cl *Client) configureRouter() {
	v1 := cl.router.Group("/")
	v1.GET("/ws", cl.handler.Handle())
}

func (cl *Client) configureWS() {
	var upgrader = websocket.Upgrader{}

	if cl.config.AppConfig.Debug {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	cl.ws = &upgrader
}

func (cl *Client) configureReceiverChan(receiver <-chan []byte) {
	cl.receiver = receiver
}

func (cl *Client) configureHandlers() {
	cl.handler = NewHandler(cl.ws, cl.receiver)
}
