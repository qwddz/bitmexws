package wsclient

import "time"

type WSMessage struct {
	Action    string
	Symbol    string    `json:"symbol"`
	FairPrice float64   `json:"fairPrice"`
	MarkPrice float64   `json:"markPrice"`
	Timestamp time.Time `json:"timestamp"`
}
