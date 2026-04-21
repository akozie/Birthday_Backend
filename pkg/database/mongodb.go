package database

import (
	"context"
	"crypto/tls"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient establishes a connection to MongoDB
func NewMongoClient(uri string) *mongo.Client {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		log.Println("MONGO_URI is empty")
		return nil
	}

	if parsed, err := url.Parse(uri); err == nil {
		logMongoURIInfo(parsed)
	}

	// Keep the startup timeout bounded, but don't crash the process if pinging
	// Atlas fails. A failed ping usually means network, TLS, or IP-allowlist
	// issues, and we want those to show up in logs without entering a restart loop.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(15*time.Second).
		SetServerSelectionTimeout(15*time.Second).
		SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to create MongoDB client: %v", err)
		return nil
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer pingCancel()

	// Ping the database to verify the connection is alive. If this fails, we
	// keep the client so the app can continue starting and surface the error.
	err = client.Ping(pingCtx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return client
	}

	log.Println("🚀 Successfully connected to MongoDB!")
	return client
}

func logMongoURIInfo(parsed *url.URL) {
	if parsed == nil {
		return
	}

	log.Printf("MongoDB URI scheme: %s", parsed.Scheme)
	if parsed.Host != "" {
		log.Printf("MongoDB host: %s", parsed.Host)
	}
	if parsed.Path != "" && parsed.Path != "/" {
		log.Printf("MongoDB database path: %s", strings.TrimPrefix(parsed.Path, "/"))
	}

	query := parsed.Query()
	if len(query) == 0 {
		log.Println("MongoDB URI query: none")
		return
	}

	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	safeValues := make([]string, 0, len(keys))
	for _, key := range keys {
		switch key {
		case "authMechanism", "authSource", "appName", "retryWrites", "tls", "w":
			safeValues = append(safeValues, key+"="+strings.Join(query[key], ","))
		default:
			safeValues = append(safeValues, key+"=[redacted]")
		}
	}

	log.Printf("MongoDB URI query params: %s", strings.Join(safeValues, ", "))
}
