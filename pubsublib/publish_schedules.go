package pubsublib

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

const SCHEDULE_TOPIC_NAME = "flatchecker_scheduled_searches"

type ScheduledSearchesMessage struct {
	ScheduleId int
	SearchId   int
}

// Publisher is an interface for publishing schedules
type Publisher interface {
	PublishSchedules(ctx context.Context, schedules []ScheduledSearchesMessage) error
}

// GCPPublisher implements Publisher using Google Cloud Pub/Sub
type GCPPublisher struct {
	client *pubsub.Client
}

// NewPublisher creates a new Publisher
func NewPublisher(client *pubsub.Client) Publisher {
	return &GCPPublisher{client: client}
}

// PublishSchedules publishes schedules to Pub/Sub
func (p *GCPPublisher) PublishSchedules(ctx context.Context, schedules []ScheduledSearchesMessage) error {
	topic, err := GetTopic(ctx, p.client, SCHEDULE_TOPIC_NAME)
	if err != nil {
		return fmt.Errorf("error getting topic: %v", err)
	}

	results := make([]*pubsub.PublishResult, len(schedules))
	for i, schedule := range schedules {
		msgBytes, err := json.Marshal(schedule)
		if err != nil {
			return fmt.Errorf("error marshaling schedule: %v", err)
		}

		msg := &pubsub.Message{
			Data: msgBytes,
		}
		result := topic.Publish(ctx, msg)
		results[i] = result
	}

	for _, result := range results {
		_, err = result.Get(ctx)
		if err != nil {
			return fmt.Errorf("error publishing schedule: %v", err)
		}
	}

	return nil
}

// PublishSchedules is a convenience function that uses the client directly
func PublishSchedules(ctx context.Context, schedules []ScheduledSearchesMessage, client *pubsub.Client) error {
	publisher := NewPublisher(client)
	return publisher.PublishSchedules(ctx, schedules)
}
