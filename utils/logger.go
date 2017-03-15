package utils

import (
	"log"
)

func Debug(v ...interface{}) {
	log.Printf("[DEBUG] %s\n", v)
}

func Error(v ...interface{}) {
	log.Printf("[DEBUG] %s\n", v)
}
