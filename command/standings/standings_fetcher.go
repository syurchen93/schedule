package main

import (
	"github.com/syurchen93/api-football-client/client"
	"github.com/syurchen93/api-football-client/request/standings"
	response "github.com/syurchen93/api-football-client/response/standings"

	"gorm.io/gorm"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"schedule/db"
	model "schedule/model/league"
	"schedule/util"
)

var apiClient *client.Client
var dbGorm *gorm.DB

func main() {
	app := &cli.App{
		Name:  "fetch-standings",
		Usage: "Fetch and persist standings from API Football. Creating teams on the fly if they don't exist.",
		Action: func(*cli.Context) error {
			apiClient = client.NewClient(util.GetEnv("API_FOOTBALL_KEY"))
			db.Init()
			dbGorm = db.Db()

			fetchAndPersistStandings()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func fetchAndPersistStandings() {
	var competitions []model.Competition
	dbGorm.Where("enabled = ?", true).Find(&competitions)

	for _, competition := range competitions {
		standingsRequest := standings.Standings{
			League: int(competition.ID),
			Season: int(competition.CurrentSeason),
		}

		standingResponse, err := apiClient.DoRequest(standingsRequest)
		if err != nil {
			panic(err)
		}

		league := standingResponse[0].(response.Standings).League

		_ = league
	}
}
