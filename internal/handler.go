package internal

import "github.com/gin-gonic/gin"

type Handler interface {
	Handle() gin.HandlerFunc
}
