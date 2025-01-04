package main

import (
	"backend/internal/infrastructure"
	"backend/internal/utils"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func init() {
	// Çalıştırma dizinini alın
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Log dosyasını çalıştırma dizininde oluştur
	logFilePath := filepath.Join(workingDir, "application.log")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	logrus.SetOutput(multiWriter)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("Application logging initialized")
}

func main() {
	//// Çalıştırma dizinini alın
	//workingDir, err := os.Getwd()
	//if err != nil {
	//	log.Fatalf("Failed to get working directory: %v", err)
	//}
	//
	//// Screenshots dizinini çalıştırma dizininde oluştur
	//dir := filepath.Join(workingDir, "screenshots")
	//if _, err := os.Stat(dir); os.IsNotExist(err) {
	//	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
	//		log.Fatalf("Failed to create screenshots directory: %v", err)
	//	}
	//	log.Println("Screenshots directory created.")
	//} else {
	//	log.Println("Screenshots directory already exists.")
	//}

	go utils.CleanUpOldFiles("screenshots", 10*time.Minute)
	app := infrastructure.SetupRouter()
	log.Println("Server is running on :8000")
	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
