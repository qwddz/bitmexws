package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

		closedWS := make(chan int)

		go func() {
			for {
				if _, _, err := ws.ReadMessage(); err != nil {
					closedWS <- 1

					break
				}
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
			case <-closedWS:
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

					pm, err := websocket.NewPreparedMessage(websocket.TextMessage, msg)
					if err != nil {
						listen = false

						break
					}

					if err := ws.WritePreparedMessage(pm); err != nil {
						listen = false
						break
					}
				}
			}
		}
	}
}
