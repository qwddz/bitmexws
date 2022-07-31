package client

import (
	"github.com/gin-gonic/gin"
	"github.com/qwddz/bitmexws/internal/config"
)

type Client struct {
	config *config.Config
	router *gin.Engine
}

func NewClient() *Client {
	return &Client{
		config: config.NewConfig(),
	}
}

func (cl *Client) ServeHTTP() error {
	cl.configureRouter()

	return nil
}

func (cl *Client) configureRouter() {
	v1 := cl.router.Group("/")
	v1.GET("/ws")
}
