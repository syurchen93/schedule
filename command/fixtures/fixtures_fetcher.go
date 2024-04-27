package main

import (
	"fmt"
	"time"

	"github.com/syurchen93/api-football-client/client"
	"github.com/syurchen93/api-football-client/request/fixture"
	response "github.com/syurchen93/api-football-client/response/fixtures"

	"log"
	"os"
	"strconv"

	"gorm.io/gorm"

	"schedule/db"
	model "schedule/model/league"
	"schedule/util"
	"schedule/util/transformer"

	"github.com/urfave/cli/v2"
)

const DEFAULT_DAYS = 14

var apiClient *client.Client
var dbGorm *gorm.DB

func main() {
	app := &cli.App{
		Name: "fetch-fixtures",
		Usage: fmt.Sprintf(
			"Fetch and persist fixtures from API. 1st arguemnt is the number of days before today, 2nd argument is the number of days after today. Default is %d days before and %d after.",
			DEFAULT_DAYS,
			DEFAULT_DAYS,
		),
		Action: func(context *cli.Context) error {
			daysBefore, err := strconv.Atoi(context.Args().Get(0))
			if err != nil {
				daysBefore = DEFAULT_DAYS
			}
			daysAfter, err := strconv.Atoi(context.Args().Get(1))
			if err != nil {
				daysAfter = DEFAULT_DAYS
			}

			apiClient = client.NewClient(util.GetEnv("API_FOOTBALL_KEY"), client.RateLimiterSettings{})
			db.Init()
			dbGorm = db.Db()

			fetchAndPersistFixtures(daysBefore, daysAfter)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func fetchAndPersistFixtures(daysBefore, daysAfter int) {
	var errorCount int
	var fixturesCreatedCount int
	var competitionCount int
	var competitions []model.Competition
	now := time.Now()
	dbGorm.Where("enabled = ?", true).Find(&competitions)

	for _, competition := range competitions {
		fixturesRequest := fixture.Fixture{
			League: int(competition.ID),
			Season: int(competition.CurrentSeason),
			From:   now.AddDate(0, 0, -daysBefore),
			To:     now.AddDate(0, 0, daysAfter),
		}

		fixturesResponse, err := apiClient.DoRequest(fixturesRequest)
		if err != nil {
			panic(err)
		}

		for _, fixtureResponse := range fixturesResponse {
			fixture := transformer.CreateFixtureFromResponse(fixtureResponse.(response.Fixture))
			result := dbGorm.Where("id = ?", fixture.ID).Assign(fixture).FirstOrCreate(&fixture)
			if result.Error != nil {
				errorCount++
			}
			if result.RowsAffected > 0 {
				fixturesCreatedCount++
			}
		}
		competitionCount++
	}
	fmt.Printf(
		"Created %d fixtures %d days before and %d after in %d competitions with %d errors\n",
		fixturesCreatedCount,
		daysBefore,
		daysAfter,
		competitionCount,
		errorCount,
	)
}
