package client

import (
	"github.com/gin-gonic/gin"
	"github.com/qwddz/bitmexws/internal/statistics"
	"net/http"
)

const DefaultLimit = 100

type Stats struct {
	store *statistics.StatRepo
}

type request struct {
	Limit  int    `json:"limit" form:"limit" binding:"required,min=1,max=1000"`
	LastID int    `json:"last_id" form:"last_id" binding:"required"`
	Symbol string `json:"symbol" form:"symbol"`
}

func NewStatsHandler(stats *statistics.StatRepo) *Stats {
	return &Stats{store: stats}
}

func (h *Stats) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := new(request)
		if err := c.ShouldBind(r); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})

			return
		}

		stat, err := h.store.Find(c, r.LastID, h.prepareLimit(r.Limit), r.Symbol)
		if err != nil {
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
