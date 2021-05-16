package queue

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type Publisher interface {
	Publisher(ctx context.Context, topic string, data []byte)
}

type DefaultPublisher struct {
	client *pubsub.Client
}

func NewDefaultPublisher(client *pubsub.Client) DefaultPublisher {
	return DefaultPublisher{client: client}
}

func (p DefaultPublisher) Publisher(ctx context.Context, topic string, data []byte) {
	message := &pubsub.Message{
		Data: data,
	}

	p.client.Topic(topic).Publish(ctx, message)
}
