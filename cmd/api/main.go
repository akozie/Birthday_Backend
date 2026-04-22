package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"

	httpdelivery "github.com/akozie/babe-25th-backend/internal/delivery/http"
	"github.com/akozie/babe-25th-backend/bootstrap"
	"github.com/akozie/babe-25th-backend/internal/repository/mongodb"
	"github.com/akozie/babe-25th-backend/internal/usecase"
	"github.com/akozie/babe-25th-backend/pkg/database"
	"github.com/akozie/babe-25th-backend/pkg/media"
)

func main() {
	log.Info("--- STARTUP: Initializing Application ---")
	app := bootstrap.App()
	env := app.Env

	mongoURI := env.MongoURI
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGO_URI")
	}
	if mongoURI == "" {
		log.Fatal("MONGO_URI is missing from environment variables")
	}

	cloudinaryURL := env.CloudinaryURL
	if cloudinaryURL == "" {
		cloudinaryURL = os.Getenv("CLOUDINARY_URL")
	}
	if cloudinaryURL == "" {
		log.Fatal("CLOUDINARY_URL is missing from environment variables")
	}

	timeout := time.Duration(env.ContextTimeout) * time.Second
	if timeout == 0 {
		timeout = 2 * time.Second
	}

	mongoClient := database.NewMongoClient(mongoURI)
	if mongoClient == nil {
		log.Fatal("MongoDB client could not be initialized")
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Errorf("error disconnecting from MongoDB: %v", err)
		}
	}()

	db := mongoClient.Database("birthday")
	cloudinaryService, err := media.NewCloudinaryService(cloudinaryURL)
	if err != nil {
		log.Fatal(err)
	}

	memRepo := mongodb.NewMongoMemoryRepository(db)
	memUsecase := usecase.NewMemoryUsecase(memRepo, cloudinaryService, timeout)
	memHandler := &httpdelivery.MemoryHandler{Usecase: memUsecase}

	guestRepo := mongodb.NewMongoGuestbookRepository(db)
	guestUsecase := usecase.NewGuestbookUsecase(guestRepo, timeout)
	guestHandler := &httpdelivery.GuestbookHandler{Usecase: guestUsecase}

	msgRepo := mongodb.NewMessageRepository(db)
	msgUsecase := usecase.NewMessageUsecase(msgRepo, timeout)
	messageHandler := &httpdelivery.MessageHandler{Usecase: msgUsecase}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/memories", memHandler.Create)
		r.Get("/memories", memHandler.GetAll)
		r.Post("/guestbook", guestHandler.Create)
		r.Get("/guestbook", guestHandler.GetAll)
		r.Get("/messages", messageHandler.GetAll)
		r.Post("/messages", messageHandler.Create)
	})

	port := os.Getenv("PORT")
	address := env.ServerAddress
	if port != "" {
		address = ":" + port
	}
	if address == "" {
		address = ":8080"
	}

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
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
	_ = srv.Shutdown(ctx)
	log.Info("shutting down")
	os.Exit(0)
}
