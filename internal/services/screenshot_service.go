package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func CaptureScreenshot(url string) (string, error) {
	fmt.Println("Ekran görüntüsü alımı başlatıldı: ", url)

	// Ekran görüntüsü dosya yolunu belirle
	fileName := generateFileName(url) + ".png"
	outputDir := "screenshots"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("Klasör oluşturulamadı: %w", err)
	}
	filePath := filepath.Join(outputDir, fileName)

	// Tarayıcı başlatma
	launchURL := launcher.New().
		Headless(true).  // Başsız mod
		NoSandbox(true). // Sandbox'u devre dışı bırak
		MustLaunch()

	log.Printf("Tarayıcı başlatıldı: %s", launchURL)

	browser := rod.New().ControlURL(launchURL)
	if err := browser.Connect(); err != nil {
		return "", fmt.Errorf("Tarayıcıya bağlanılamadı: %w", err)
	}
	defer browser.Close()

	// Sayfayı aç ve ekran görüntüsü al
	page := browser.MustPage(url).MustWaitLoad()
	screenshot := page.MustScreenshot()
	if err := os.WriteFile(filePath, screenshot, 0644); err != nil {
		return "", fmt.Errorf("Ekran görüntüsü kaydedilemedi: %w", err)
	}

	fmt.Println("Ekran görüntüsü başarıyla alındı: ", filePath)
	return fileName, nil
}

// generateFileName: URL'den hashlenmiş bir dosya adı oluşturur
func generateFileName(url string) string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
