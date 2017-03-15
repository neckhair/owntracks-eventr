package utils

import (
	"log"
)

func Debug(text string) {
	log.Printf("[DEBUG] %s\n", text)
}

func Error(text string) {
	log.Printf("[ERROR] %s\n", text)
}
