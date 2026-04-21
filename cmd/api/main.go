package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/oguzhantasimaz/Go-Clean-Architecture-Template/api/route"
	"github.com/oguzhantasimaz/Go-Clean-Architecture-Template/bootstrap"
	"github.com/oguzhantasimaz/Go-Clean-Architecture-Template/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.MySql
	defer app.CloseDBConnection()

	utils.MigrateDB(db)

	timeout := time.Duration(env.ContextTimeout) * time.Second

	r := mux.NewRouter()

	route.Setup(env, timeout, db, r)

	// --- FIX: Logic to handle Railway's dynamic port ---
	port := os.Getenv("PORT")
	address := env.ServerAddress // Defaults to your .env/config value (e.g., localhost:8080)

	// If we are in a cloud environment (like Railway), use the PORT env var
	// and listen on all interfaces by prepending just the colon ":"
	if port != "" {
		address = ":" + port
	}
	// ----------------------------------------------------

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Printf("Server starting on %s", address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err)
		}
	}()

	log.Info("server started")

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	srv.Shutdown(ctx)
	log.Info("shutting down")
	os.Exit(0)
}