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
	log.Info("--- STARTUP: Initializing Application ---")
	
	log.Info("Step 1: Running bootstrap.App()")
	app := bootstrap.App()
	log.Info("Step 2: Bootstrap successful")

	env := app.Env
	db := app.MySql
	defer app.CloseDBConnection()

	log.Info("Step 3: Starting Database Migration")
	utils.MigrateDB(db)
	log.Info("Step 4: Migration finished")

	if db == nil {
		log.Warn("database is not configured; the server will start, but DB-backed routes will return errors")
	}

	timeout := time.Duration(env.ContextTimeout) * time.Second

	r := mux.NewRouter()
	route.Setup(env, timeout, db, r)

	port := os.Getenv("PORT")
	address := env.ServerAddress
	if port != "" {
		address = ":" + port
	}

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Infof("--- SUCCESS: Server listening on %s ---", address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Server error: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	srv.Shutdown(ctx)
	log.Info("shutting down")
	os.Exit(0)
}
