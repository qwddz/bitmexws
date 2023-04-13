package internal

import "github.com/gin-gonic/gin"

type Router interface {
	RegisterRouter(r *gin.RouterGroup)
}
