package pubsublib

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

const SCHEDULE_TOPIC_NAME = "flatchecker_scheduled_searches"

type Schedule struct {
	ScheduleId int
	SearchId   int
}

func PublishSchedules(ctx context.Context, schedules []Schedule, client *pubsub.Client) error {
	topic, err := GetTopic(ctx, client, SCHEDULE_TOPIC_NAME)
	if err != nil {
		return err
	}

	results := make([]*pubsub.PublishResult, len(schedules))
	for i, schedule := range schedules {
		msgBytes, err := json.Marshal(schedule)
		if err != nil {
			return err
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
			return err
		}
	}

	return nil
}
