package db

import (
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"

	"os"
)

var db *gorm.DB
var err error

func Init() {

	dsn := os.Getenv("DB_DSN")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		panic("DB Connection Error")
	}
}

func Db() *gorm.DB {
	return db 
}