package main

import (
	"context"
	"flatchecker-scheduler/db"
	"flatchecker-scheduler/pubsublib"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()
	pubsubClient, err := pubsublib.GetClient(ctx)
	handleError("error creating pubsub client", err)

	topic, err := pubsubClient.CreateTopic(ctx, "test_topic")
	handleError("error creating topic", err)

	res := topic.Publish(ctx, &pubsub.Message{
		Data: []byte("hello world"),
	})

	msgID, err := res.Get(ctx)
	handleError("error sending message", err)
	fmt.Println(msgID)

	config, err := ReadConfig("C:\\Users\\kiera\\flatchecker\\flatchecker-database\\setup\\db_credentials.txt")
	handleError("error reading config", err)
	fmt.Println(config)

	db, err := db.GetDB(config)
	handleError("error connecting to db", err)
	defer db.Close()

	stmt, err := db.Prepare("INSERT into User (UserName, Email) VALUES (?, ?)")
	handleError("error preparing statement", err)
	defer stmt.Close()
	//stmt.Exec("tony", "tony@tonymail.com")

	fmt.Println("Hello, World!")
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
		out[words[0]] = words[1]
	}

	return out, nil
}

func handleError(errorMessage string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", errorMessage, err))
	}
}
