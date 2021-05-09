package utils

import (
	"log"
	"os"
)

func AppendToFile(filePath string, content string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		LogError(err)
	}
}
