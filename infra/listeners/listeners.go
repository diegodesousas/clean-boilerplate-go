package listeners

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type Listener interface {
	Subscription() string
	Receive(ctx context.Context, message *pubsub.Message)
}
