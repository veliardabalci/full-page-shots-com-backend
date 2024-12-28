package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateFileName(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}
