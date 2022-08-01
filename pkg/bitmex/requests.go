package bitmex

type Request struct {
	Action  SubscriptionAction `json:"action"`
	Symbols []string           `json:"symbols"`
}
