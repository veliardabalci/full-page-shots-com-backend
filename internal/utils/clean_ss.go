package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// CleanUpOldFiles belirli bir klasördeki dosyaları kontrol eder ve belirli bir yaşı aşanları siler
func CleanUpOldFiles(dir string, maxAge time.Duration) {
	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Printf("Failed to read directory %s: %v", dir, err)
			time.Sleep(1 * time.Minute) // Hata olduğunda bir dakika bekle
			continue
		}

		now := time.Now()
		for _, file := range files {
			if !file.Type().IsRegular() {
				continue
			}

			filePath := filepath.Join(dir, file.Name())
			info, err := os.Stat(filePath)
			if err != nil {
				log.Printf("Failed to stat file %s: %v", filePath, err)
				continue
			}

			// Dosyanın yaşı kontrol edilir
			if now.Sub(info.ModTime()) > maxAge {
				if err := os.Remove(filePath); err != nil {
					log.Printf("Failed to delete file %s: %v", filePath, err)
				} else {
					log.Printf("Deleted old file: %s", filePath)
				}
			}
		}

		time.Sleep(1 * time.Minute) // Kontrol süresini ayarla (ör. 1 dakika)
	}
}
