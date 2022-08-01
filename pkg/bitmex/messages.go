package bitmex

import "time"

type SubscriptionMessage struct {
	Action SubscriptionAction `json:"op"`
	Args   []SubscriptionType `json:"args"`
}

type ReceiveMessage struct {
	Table  string             `json:"table"`
	Action SubscriptionAction `json:"action"`
	Data   []MessageData      `json:"data"`
}

type MessageData struct {
	Symbol                string    `json:"symbol"`
	FairPrice             float64   `json:"fairPrice"`
	FairBasis             float64   `json:"fairBasis"`
	LastPriceProtected    float64   `json:"lastPriceProtected"`
	MarkPrice             float64   `json:"markPrice"`
	OpenValue             float64   `json:"openValue"`
	IndicativeSettlePrice float64   `json:"indicativeSettlePrice"`
	Timestamp             time.Time `json:"timestamp"`
}
