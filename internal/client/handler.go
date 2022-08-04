package client

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/qwddz/bitmexws/pkg/bitmex"
	"log"
)

type Handler struct {
	receiver <-chan []byte

	ws *websocket.Upgrader
}

func NewHandler(ws *websocket.Upgrader, receiver <-chan []byte) *Handler {
	return &Handler{ws: ws, receiver: receiver}
}

func (h *Handler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := h.ws.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)

			return
		}

		defer func() {
			if err := ws.Close(); err != nil {
				log.Printf("ws close error: %s", err.Error())
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
					log.Printf("receive user message error: %s", err)
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
			case msg, closed := <-h.receiver:
				{
					if closed == false {
						listen = false

						break
					}

					message, err := h.formatMessage(msg)
					if err != nil {
						continue
					}

					if subscribed != bitmex.ActionSubscribe || h.symbolFiltered(request.Symbols, message.Symbol) == false {
						continue
					}

					if err := ws.WriteJSON(message); err != nil {
						listen = false

						break
					}
				}
			}
		}
	}
}

func (h *Handler) formatMessage(jsonMsg []byte) (*bitmex.Response, error) {
	var msg bitmex.ReceiveMessage

	if err := json.Unmarshal(jsonMsg, &msg); err != nil {
		return nil, err
	}

	if msg.Action != bitmex.ActionUpdate {
		return nil, errors.New("message is not supported")
	}

	data := msg.Data[0]

	return &bitmex.Response{
		Symbol:    data.Symbol,
		Price:     data.MarkPrice,
		Timestamp: data.Timestamp,
	}, nil
}

func (h *Handler) symbolFiltered(request []string, currentSymbol string) bool {
	if len(request) == 0 {
		return true
	}

	for _, rs := range request {
		if rs == currentSymbol {
			return true
		}
	}

	return false
}
