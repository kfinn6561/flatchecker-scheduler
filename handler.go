package main

import (
	"context"
	"database/sql"
	"flatchecker-scheduler/mapper"
	"flatchecker-scheduler/pubsublib"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
)

func GetHandler(ctx context.Context, dbConn *sql.DB, pubsubClient *pubsub.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received request")
		err := readAndPublishSchedules(ctx, dbConn, pubsubClient)
		if err != nil {
			fmt.Println("error reading and publishing schedules:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func readAndPublishSchedules(ctx context.Context, dbConn *sql.DB, pubsubClient *pubsub.Client) error {
	schedules, err := GetAndUpdateSchedules(dbConn)
	if err != nil {
		return err
	}

	pubsubSchedules := mapper.MapSchedules(schedules)
	return pubsublib.PublishSchedules(ctx, pubsubSchedules, pubsubClient)
}
