package services

import (
	"backend/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

// CaptureScreenshot takes a screenshot of the given URL and saves it as a PNG file.
func CaptureScreenshot(request models.Request, url string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"url":         url,
		"device_type": request.DeviceType,
		"width":       request.Width,
		"height":      request.Height,
		"format":      request.Format,
	}).Info("Received screenshot capture request")

	// Generate a unique filename for the screenshot
	fileName := generateFileName(url) + ".png"
	outputDir := "screenshots"

	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, fileName)

	// Set ChromeDP options for a headless browser
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,              // Run in headless mode
		chromedp.DisableGPU,            // Disable GPU for more stability
		chromedp.NoFirstRun,            // Skip the first run tasks
		chromedp.NoDefaultBrowserCheck, // Skip default browser check
	)

	// Create a Chrome context
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(logrus.Infof), // Log for debugging purposes
	)
	defer cancel()

	// Set a timeout for the operation
	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var buf []byte

	// Select device emulation based on deviceType
	var emulation chromedp.Action
	switch request.DeviceType {
	case "mobile":
		emulation = chromedp.Emulate(device.IPhone13Pro) // Mobile device
	case "tablet":
		emulation = chromedp.Emulate(device.IPadPro11) // Tablet device
	case "desktop":
		emulation = chromedp.EmulateViewport(request.Width, request.Height) // Desktop with custom viewport
	default:
		return "", fmt.Errorf("invalid deviceType: %s", request.DeviceType)
	}

	// Run tasks to navigate to the URL, emulate device, and capture a full-page screenshot
	err := chromedp.Run(timeoutCtx, chromedp.Tasks{
		emulation,                          // Emulate the selected device
		chromedp.Navigate(url),             // Navigate to the URL
		chromedp.WaitReady("body"),         // Wait for the <body> element to be ready
		chromedp.FullScreenshot(&buf, 100), // Capture the screenshot
	})
	if err != nil {
		return "", fmt.Errorf("failed to capture screenshot: %w", err)
	}

	// Write the screenshot to a file
	if err = os.WriteFile(filePath, buf, 0644); err != nil {
		return "", fmt.Errorf("failed to save screenshot: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"file_name": fileName,
		"url":       url,
	}).Info("Screenshot successfully saved")

	if request.Format == "PDF" {
		pdfFileName, err := convertPNGToPDF(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to convert PNG to PDF: %w", err)
		}
		return filepath.Base(pdfFileName), nil
	} else if request.Format == "JPG" {
		jpgFileName, err := convertPNGToJPG(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to convert PNG to JPG: %w", err)
		}
		return filepath.Base(jpgFileName), nil
	}

	return fileName, nil
}

// generateFileName creates a unique file name based on the URL and the current timestamp.
func generateFileName(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url + time.Now().String()))
	return hex.EncodeToString(hasher.Sum(nil))[:20]
}

// convertPNGToPDF converts a PNG image file to a PDF file.
func convertPNGToPDF(pngPath string) (string, error) {
	file, err := os.Open(pngPath)
	if err != nil {
		return "", fmt.Errorf("failed to open PNG file: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode PNG file: %w", err)
	}

	imgWidth := float64(img.Bounds().Dx()) / 72.0 // Convert pixels to inches
	imgHeight := float64(img.Bounds().Dy()) / 72.0

	pdf := gofpdf.New("P", "in", "Letter", "")
	pdf.AddPageFormat("P", gofpdf.SizeType{Wd: imgWidth, Ht: imgHeight})

	file.Seek(0, 0) // Reset the file pointer to the beginning

	imgPath := filepath.Base(pngPath)
	if pdf.RegisterImageOptionsReader(imgPath, gofpdf.ImageOptions{ImageType: "PNG"}, file) == nil {
		return "", fmt.Errorf("failed to register image")
	}

	pdfFileName := pngPath[:len(pngPath)-4] + ".pdf"
	if err := pdf.OutputFileAndClose(pdfFileName); err != nil {
		return "", fmt.Errorf("failed to save PDF file: %w", err)
	}

	return filepath.Base(pdfFileName), nil
}

// convertPNGToJPG converts a PNG image file to a JPG file.
func convertPNGToJPG(pngPath string) (string, error) {
	file, err := os.Open(pngPath)
	if err != nil {
		return "", fmt.Errorf("failed to open PNG file: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode PNG file: %w", err)
	}

	jpgPath := pngPath[:len(pngPath)-4] + ".jpg"
	outFile, err := os.Create(jpgPath)
	if err != nil {
		return "", fmt.Errorf("failed to create JPG file: %w", err)
	}
	defer outFile.Close()

	// Encode the image to JPG format
	err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 80})
	if err != nil {
		return "", fmt.Errorf("failed to encode JPG file: %w", err)
	}

	return filepath.Base(jpgPath), nil
}
