package bitmex

type SubscriptionType string

const (
	TypeInstrument SubscriptionType = "instrument"
)

type SubscriptionAction string

const (
	ActionSubscribe   SubscriptionAction = "subscribe"
	ActionUnsubscribe SubscriptionAction = "unsubscribe"
	ActionUpdate      SubscriptionAction = "update"
)
