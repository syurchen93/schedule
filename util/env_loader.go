package util

import (
	"os"
	"github.com/joho/godotenv"
	"fmt"
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