package main

import (
	"context"
	"database/sql"
	"flatchecker-scheduler/db"
	"flatchecker-scheduler/mapper"
	"flatchecker-scheduler/pubsublib"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
)

func main() {
	fmt.Println("starting")
	if len(os.Args) < 2 {
		fmt.Println("Usage: flatchecker-scheduler <config-filename>")
		os.Exit(1)
	}
	configFilename := os.Args[1]

	ctx := context.Background()
	pubsubClient, err := pubsublib.GetClient(ctx)
	handleError("error creating pubsub client", err)
	fmt.Println("successfully created pubsub client")

	config, err := ReadConfig(configFilename)
	handleError("error reading config", err)
	fmt.Println("successfully read config")

	dbConn, err := db.GetDB(config)
	handleError("error connecting to db", err)
	defer dbConn.Close()
	fmt.Println("successfully connected to database")

	handler := GetHandler(ctx, dbConn, pubsubClient)
	http.HandleFunc("/", handler)
	fmt.Println("starting server on :8080")
	handleError("error starting server", http.ListenAndServe(":8080", nil))
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
