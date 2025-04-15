package pubsublib

import (
	"context"

	"cloud.google.com/go/pubsub"
)

const PROJECT_ID = "flatchecker"

func GetClient(ctx context.Context) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, PROJECT_ID)
}
