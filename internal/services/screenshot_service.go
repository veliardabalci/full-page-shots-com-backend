package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

func CaptureScreenshot(url string) (string, error) {
	// Ekran görüntüsü dosya yolunu belirle
	fileName := generateFileName(url) + ".png"
	outputDir := "screenshots"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create screenshots directory: %w", err)
	}
	filePath := filepath.Join(outputDir, fileName)

	// Chrome context oluştur
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Timeout ekle
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Ekran görüntüsü al
	var buf []byte
	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.FullScreenshot(&buf, 90),
	})
	if err != nil {
		return "", fmt.Errorf("failed to capture screenshot: %w", err)
	}

	// Görüntüyü dosyaya kaydet
	if err := os.WriteFile(filePath, buf, 0644); err != nil {
		return "", fmt.Errorf("failed to save screenshot: %w", err)
	}
	filePath = fileName
	return filePath, nil
}

// generateFileName: URL'den hashlenmiş bir dosya adı oluşturur
func generateFileName(url string) string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
