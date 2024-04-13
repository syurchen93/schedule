package util

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var isLoaded bool = false

func GetEnv(key string) string {
	if !isLoaded {
		LoadEnv()
	}
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %s", err))
	}
	isLoaded = true
}
