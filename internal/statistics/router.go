package statistics

import (
	"github.com/gin-gonic/gin"
	"github.com/qwddz/bitmexws/internal"
	"github.com/qwddz/bitmexws/pkg/logger"
	"github.com/qwddz/bitmexws/pkg/store"
)

type Router struct {
	infoHandler internal.Handler
}

func NewRouter(store *store.Store, log logger.Logger) *Router {
	return &Router{
		infoHandler: newHandler(store, log),
	}
}

func (r *Router) RegisterRouter(rg *gin.RouterGroup) {
	rg.GET("/", r.infoHandler.Handle())
}
