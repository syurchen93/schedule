package db

import (
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

func InitDbOrPanic() *gorm.DB {
	db, err = InitDB()
	if err != nil {
		panic(err)
	}
	return db
}

func InitDB() (*gorm.DB, error) {
	CloseDB()

	dsn := util.GetEnv("DB_DSN")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, err
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
		return nil, err
	}

	return db, nil
}

func CloseDB() {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	if err := sqlDB.Close(); err != nil {
		panic(err)
	}
}
