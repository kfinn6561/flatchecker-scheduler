package main

import (
	"context"
	"database/sql"
	"flatchecker-scheduler/db"
	"flatchecker-scheduler/pubsublib"
	"fmt"
	"net/http"
	"os"
	"strings"

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

	config, err := InitializeConfig(configFilename)
	handleError("error reading config", err)

	pubsubClient, err := InitializePubSub(ctx)
	handleError("error creating pubsub client", err)

	dbConn, err := InitializeDB(config)
	handleError("error connecting to db", err)
	defer dbConn.Close()

	handleError("error starting server", StartServer(ctx, dbConn, pubsubClient))
}

// InitializeConfig reads and parses the configuration file
func InitializeConfig(filename string) (map[string]string, error) {
	config, err := ReadConfig(filename)
	if err != nil {
		return nil, err
	}
	fmt.Println("successfully read config")
	return config, nil
}

// InitializePubSub creates and returns a Pub/Sub client
func InitializePubSub(ctx context.Context) (*pubsub.Client, error) {
	client, err := pubsublib.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("successfully created pubsub client")
	return client, nil
}

// InitializeDB establishes and returns a database connection
func InitializeDB(config map[string]string) (*sql.DB, error) {
	dbConn, err := db.GetDB(config)
	if err != nil {
		return nil, err
	}
	fmt.Println("successfully connected to database")
	return dbConn, nil
}

// StartServer starts the HTTP server and blocks until it terminates
func StartServer(ctx context.Context, dbConn *sql.DB, pubsubClient *pubsub.Client) error {
	handler := GetHandler(ctx, dbConn, pubsubClient)
	http.HandleFunc("/", handler)
	fmt.Println("starting server on :8080")
	return http.ListenAndServe(":8080", nil)
}

func ReadConfig(filename string) (map[string]string, error) {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string)

	lines := strings.Split(string(rawData), "\n")
	for _, line := range lines {
		// Skip empty lines
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		words := strings.Split(line, " ")
		if len(words) < 2 {
			continue // Skip malformed lines
		}
		out[words[0]] = strings.TrimSpace(words[1])
	}

	return out, nil
}

func handleError(errorMessage string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", errorMessage, err))
	}
}
