package main

import (
	"context"
	"database/sql"
	"flatchecker-scheduler/db"
	"flatchecker-scheduler/mapper"
	"flatchecker-scheduler/pubsublib"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
)

func main() {
	fmt.Println("starting")
	ctx := context.Background()
	pubsubClient, err := pubsublib.GetClient(ctx)
	handleError("error creating pubsub client", err)
	fmt.Println("successfully created pubsub client")

	config, err := ReadConfig("db_credentials.txt")
	handleError("error reading config", err)
	fmt.Println("successfully read config")

	dbConn, err := db.GetDB(config)
	handleError("error connecting to db", err)
	defer dbConn.Close()
	fmt.Println("successfully connected to database")

	for {
		err = readAndPublishSchedules(ctx, dbConn, pubsubClient)
		if err != nil {
			fmt.Println("Error reading schedules", err)
		}
		time.Sleep(1 * time.Second)
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

func ReadConfig(filename string) (map[string]string, error) {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string)

	lines := strings.Split(string(rawData), "\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		out[words[0]] = strings.TrimSpace(words[1])
	}

	return out, nil
}

func handleError(errorMessage string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", errorMessage, err))
	}
}
