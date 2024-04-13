package main

import (
	"net/http"
	"schedule/db"
	"schedule/util"

	"github.com/labstack/echo/v4"

	"fmt"
)

func main() {
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

	err = dbGorm.Ping()
	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(
		fmt.Sprintf(":%s", util.GetEnv("WEB_SERVER_PORT")),
	))
}
