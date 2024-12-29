package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
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

	// Chrome seçeneklerini belirle
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium"),             // Chrome'un tam yolu
		chromedp.Headless,                                  // Başsız mod
		chromedp.DisableGPU,                                // GPU'yu devre dışı bırak
		chromedp.NoSandbox,                                 // Sandbox'u devre dışı bırak
		chromedp.Flag("disable-software-rasterizer", true), // Yazılım rasterizer'ı devre dışı bırak
	)

	// Chrome context oluştur
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf), // Hata ayıklama için loglama
	)
	defer cancel()

	// Timeout ekle
	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second) // Timeout artırıldı
	defer cancel()

	// Ekran görüntüsü al
	var buf []byte
	err := chromedp.Run(timeoutCtx, chromedp.Tasks{
		chromedp.Navigate(url),            // URL'ye git
		chromedp.FullScreenshot(&buf, 90), // Ekran görüntüsü al
	})
	if err != nil {
		return "", fmt.Errorf("Ekran görüntüsü alınamadı: %w", err)
	}

	// Görüntüyü dosyaya kaydet
	if err := os.WriteFile(filePath, buf, 0644); err != nil {
		return "", fmt.Errorf("Ekran görüntüsü kaydedilemedi: %w", err)
	}

	fmt.Println("Ekran görüntüsü başarıyla alındı: ", filePath)
	return fileName, nil
}

// generateFileName: URL'den hashlenmiş bir dosya adı oluşturur
func generateFileName(url string) string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
