package main

import (
	"net/http"
	"os"
	"schedule/db"

	"github.com/labstack/echo/v4"

	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil{
	 panic(fmt.Sprintf("Error loading .env file: %s", err))
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"hello": "world",
		})
	})

	db.Init()
	gorm := db.Db()

	dbGorm, err := gorm.DB()
	if err != nil {
		panic(err)
	}

	dbGorm.Ping()
	
	e.Logger.Fatal(e.Start(
		fmt.Sprintf(":%s", os.Getenv("WEB_SERVER_PORT")),
	))
}