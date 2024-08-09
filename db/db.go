package db

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"schedule/model"
	"schedule/model/bot"
	"schedule/model/league"
	"schedule/util"
)

var db *gorm.DB
var err error

func Init(ctx ...context.Context) {

	dsn := util.GetEnv("DB_DSN")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		panic("DB Connection Error")
	}

	if len(ctx) > 0 {
		db = db.WithContext(ctx[0])
	}

	err = db.AutoMigrate(
		&league.Country{},
		&league.Competition{},
		&league.Standing{},
		&league.Fixture{},

		&model.Team{},

		&bot.User{},
		&bot.FavTeam{},
		&bot.Alert{},
		&bot.UserShare{},
	)
	if err != nil {
		panic("DB Migration Error")
	}
}

func Db() *gorm.DB {
	return db
}
