package bitmex

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
)

type WSClient struct {
	conn *websocket.Conn
}

func NewWSClient() *WSClient {
	return &WSClient{}
}

func (wsc *WSClient) Connect(ctx context.Context, url string) {
	if wsc.conn != nil {
		return
	}

	ws, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		log.Printf("cannot connect to websocket: %s with %s", err.Error(), url)

		return
	}

	wsc.conn = ws
}

func (wsc *WSClient) ReadMessage(receiver chan []byte) error {
	_, msg, err := wsc.conn.ReadMessage()

	if err != nil {
		return err
	}

	receiver <- msg

	return nil
}

func (wsc *WSClient) Subscribe(topic SubscriptionType) error {
	return wsc.sendMessage(ActionSubscribe, topic)
}

func (wsc *WSClient) Unsubscribe(topic SubscriptionType) error {
	return wsc.sendMessage(ActionUnsubscribe, topic)
}

func (wsc *WSClient) Close() error {
	return wsc.conn.Close()
}

func (wsc *WSClient) sendMessage(action SubscriptionAction, topic SubscriptionType) error {
	return wsc.conn.WriteJSON(SubscriptionMessage{
		Action: action,
		Args:   []SubscriptionType{topic},
	})
}
