package pubsublib

import (
	"context"

	"cloud.google.com/go/pubsub"
)

const PROJECT_ID = "flatchecker"

func GetClient(ctx context.Context) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, PROJECT_ID)
}

func GetTopic(ctx context.Context, client *pubsub.Client, topicID string) (*pubsub.Topic, error) {
	topic := client.Topic(topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return topic, nil
	}

	return client.CreateTopic(ctx, topicID)
}
