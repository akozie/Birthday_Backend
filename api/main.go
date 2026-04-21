package main

import (
	"log"
	"net/http"
	"os"
    "github.com/go-chi/cors"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	
	// MAKE SURE TO UPDATE THESE PATHS TO YOUR MODULE NAME
	httpdelivery "github.com/akozie/babe-25th-backend/internal/delivery/http"
	"github.com/akozie/babe-25th-backend/internal/repository/mongodb"
	"github.com/akozie/babe-25th-backend/internal/usecase"
	"github.com/akozie/babe-25th-backend/pkg/database"
	"github.com/akozie/babe-25th-backend/pkg/media"
)


func main() {
	// 1. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Initialize Infrastructure
	mongoClient := database.NewMongoClient(os.Getenv("MONGO_URI"))
	db := mongoClient.Database("birthday")
	
	cloudinaryService, err := media.NewCloudinaryService(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// 3. Wire Dependencies
	memRepo := mongodb.NewMongoMemoryRepository(db)
	memUsecase := usecase.NewMemoryUsecase(memRepo, cloudinaryService)
	memHandler := &httpdelivery.MemoryHandler{Usecase: memUsecase}

	guestRepo := mongodb.NewMongoGuestbookRepository(db)
	guestUsecase := usecase.NewGuestbookUsecase(guestRepo)
	guestHandler := &httpdelivery.GuestbookHandler{Usecase: guestUsecase}

	msgRepo := mongodb.NewMessageRepository(db)
	msgUsecase := usecase.NewMessageUsecase(msgRepo)
	messageHandler := &httpdelivery.MessageHandler{Usecase: msgUsecase}
	
	// 4. Setup Routes
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
    // Update this line to allow both 3000 and 3001, or just 3001 if that's your new port
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
		// ADD THESE TWO LINES:
    	r.Get("/messages", messageHandler.GetAll) // Assuming you have a messageHandler
    	r.Post("/messages", messageHandler.Create)
	})

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s...", port)
	http.ListenAndServe(":"+port, r)
}