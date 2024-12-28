package main

import (
	"backend/internal/infrastructure"
	"backend/internal/utils"
	"log"
	"time"
)

func main() {
	// Temizleme işlemi için bir goroutine başlat
	go utils.CleanUpOldFiles("screenshots", 10*time.Minute)

	app := infrastructure.SetupRouter()

	log.Println("Server is running on :8000")
	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
