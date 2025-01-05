package main

import (
	"backend/internal/infrastructure"
	"backend/internal/utils"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"time"
)

func init() {
	// Log dosyasını bulunduğu dizine kaydet
	logFilePath := "application.log"
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Log dosyasını ve konsolu loglama için kullan
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	logrus.SetOutput(multiWriter)

	// Log formatını ayarla
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("Application logging initialized")
}

func main() {
	// "screenshots" klasörünü bulunduğu dizinde oluştur
	dir := "screenshots"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create screenshots directory: %v", err)
		}
		log.Println("Screenshots directory created.")
	} else {
		log.Println("Screenshots directory already exists.")
	}

	// Eski dosyaları temizlemek için goroutine başlat
	go utils.CleanUpOldFiles(dir, 10*time.Minute)

	// Uygulama router'ını başlat
	app := infrastructure.SetupRouter()

	log.Println("Server is running on :8000")
	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
