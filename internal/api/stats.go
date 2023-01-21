package api

import (
	"github.com/gin-gonic/gin"
	"github.com/qwddz/bitmexws/internal/statistics"
	"github.com/qwddz/bitmexws/pkg/logger"
	"github.com/qwddz/bitmexws/pkg/store"
	"net/http"
)

const DefaultLimit = 100

type Stats struct {
	store *store.Store
	log   logger.Logger
}

type request struct {
	Limit  int    `json:"limit" form:"limit" binding:"required,min=1,max=1000"`
	LastID int    `json:"last_id" form:"last_id" binding:"required"`
	Symbol string `json:"symbol" form:"symbol"`
}

func NewStatsHandler(store *store.Store, log logger.Logger) *Stats {
	return &Stats{store: store, log: log}
}

func (h *Stats) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := new(request)
		if err := c.ShouldBind(r); err != nil {
			h.log.Errorln("statistics:", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})

			return
		}

		stat, err := statistics.New(h.store.SlaveConnection()).Find(c, r.LastID, h.prepareLimit(r.Limit), r.Symbol)
		if err != nil {
			h.log.Errorln("statistics:", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})

			return
		}

		c.AsciiJSON(http.StatusOK, gin.H{
			"msg":  "Success",
			"data": stat,
		})
	}
}

func (h *Stats) prepareLimit(limit int) int {
	if limit == 0 {
		limit = DefaultLimit
	}

	return limit
}
