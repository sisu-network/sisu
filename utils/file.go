package utils

import (
	"os"

	"github.com/sisu-network/lib/log"
)

func AppendToFile(filePath string, content string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		log.Error(err)
	}
}

func IsFileExisted(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}
