package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"schedule/model/league"
	"schedule/util"
)

var db *gorm.DB
var err error

func Init() {

	dsn := util.GetEnv("DB_DSN")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		panic("DB Connection Error")
	}

	err = db.AutoMigrate(&league.Country{}, &league.Competition{})
	if err != nil {
		panic("DB Migration Error")
	}
}

func Db() *gorm.DB {
	return db
}
