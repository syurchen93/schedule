package db

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/driver/mysql"

	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB
var err error

func Init() {
	err := godotenv.Load(".env")
	if err != nil{
	 panic(fmt.Sprintf("Error loading .env file: %s", err))
	}

	dsn := os.Getenv("DB_DSN")
	fmt.Println(dsn)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("DB Connection Error")
	}
}

func Db() *gorm.DB {
	return db 
}