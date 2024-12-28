package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

// CaptureScreenshot belirli bir URL'nin ekran görüntüsünü alır
func CaptureScreenshot(url string) (string, error) {
	// Ekran görüntüsü dosya yolunu belirle
	fileName := generateFileName(url) + ".png"
	outputDir := "screenshots"

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create screenshots directory: %w", err)
	}
	filePath := filepath.Join(outputDir, fileName)

	// Tarayıcı binary'sini belirtin (gerektiğinde özelleştirin)
	chromePath := "/usr/bin/google-chrome"
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath), // Tarayıcı binary yolu
		chromedp.Headless,             // Headless mod
		chromedp.DisableGPU,           // GPU'yu devre dışı bırak
		chromedp.NoSandbox,            // Sandbox devre dışı (sunucularda gerekebilir)
	)

	// Context oluştur ve timeout ekle
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Ekran görüntüsü al
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.FullScreenshot(&buf, 90),
	}); err != nil {
		return "", fmt.Errorf("failed to capture screenshot: %w", err)
	}

	// Görüntüyü dosyaya kaydet
	if err := os.WriteFile(filePath, buf, 0644); err != nil {
		return "", fmt.Errorf("failed to save screenshot: %w", err)
	}

	fmt.Printf("Screenshot captured successfully: %s\n", filePath)
	return fileName, nil
}

// generateFileName: URL'den hashlenmiş bir dosya adı oluşturur
func generateFileName(url string) string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
