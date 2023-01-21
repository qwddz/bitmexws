package upgrader

import (
	"context"
	"github.com/qwddz/bitmexws/pkg/bitmex"
)

type Upgrader interface {
	Connect(ctx context.Context, url string) error
	ReadMessage() ([]byte, error)
	Subscribe(topic bitmex.SubscriptionType) error
	Unsubscribe(topic bitmex.SubscriptionType) error
	Close() error
}
