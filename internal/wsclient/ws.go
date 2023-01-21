package wsclient

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/qwddz/bitmexws/pkg/bitmex"
	"github.com/qwddz/bitmexws/pkg/logger"
)

type WS struct {
	receiver <-chan WSMessage
	log      logger.Logger

	ws *websocket.Upgrader
}

func NewWSHandler(ws *websocket.Upgrader, receiver <-chan WSMessage, log logger.Logger) *WS {
	return &WS{ws: ws, receiver: receiver, log: log}
}

func (h *WS) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := h.ws.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			h.log.Errorln("user upgrader creation error", err.Error())

			return
		}

		defer func() {
			if err := ws.Close(); err != nil {
				h.log.Errorln("user ws channel close error", err.Error())
			}
		}()

		exit := make(chan int)

		var request bitmex.Request

		subscribed := bitmex.ActionUnsubscribe

		go func() {
			for {
				_, msg, err := ws.ReadMessage()
				if err != nil {
					exit <- 1

					break
				}

				if len(msg) == 0 {
					continue
				}

				if err := json.Unmarshal(msg, &request); err != nil {
					h.log.Errorln("client ws: parse received user message error", err.Error())
				}

				subscribed = request.Action
			}
		}()

		listen := true

		for listen {
			select {
			case <-c.Done():
				{
					listen = false

					break
				}
			case <-exit:
				{
					listen = false

					break
				}
			case msg, open := <-h.receiver:
				{
					if open == false {
						listen = false

						break
					}

					if subscribed != bitmex.ActionSubscribe || h.symbolFiltered(request.Symbols, msg.Symbol) == false {
						continue
					}

					if err := ws.WriteJSON(msg); err != nil {
						listen = false

						break
					}
				}
			}
		}
	}
}

func (h *WS) symbolFiltered(request []string, code string) bool {
	if len(request) == 0 {
		return true
	}

	for _, rc := range request {
		if rc == code {
			return true
		}
	}

	return false
}
