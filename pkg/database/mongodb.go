package database

import (
	"context"
	"crypto/tls"
	"log"
	"net/url"
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

	if parsed, err := url.Parse(uri); err == nil && parsed.Host != "" {
		log.Printf("MongoDB host: %s", parsed.Host)
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
