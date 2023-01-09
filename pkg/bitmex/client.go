package bitmex

import (
	"context"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	conn *websocket.Conn
}

func NewWSClient() *WSClient {
	return &WSClient{}
}

func (wsc *WSClient) Connect(ctx context.Context, url string) error {
	if wsc.conn != nil {
		return nil
	}

	ws, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return err
	}

	wsc.conn = ws

	return nil
}

func (wsc *WSClient) ReadMessage() ([]byte, error) {
	_, msg, err := wsc.conn.ReadMessage()

	return msg, err
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
