package main

import (
	"github.com/syurchen93/api-football-client/client"
	"github.com/syurchen93/api-football-client/request/league"
	"github.com/syurchen93/api-football-client/response/leagues"

	"fmt"
	"log"
	"os"

	"gorm.io/gorm"

	"schedule/db"
	model "schedule/model/league"
	"schedule/tgbot/manager"
	"schedule/util"
	"schedule/util/transformer"

	"github.com/urfave/cli/v2"
)

var apiClient *client.Client
var dbGorm *gorm.DB

func main() {
	app := &cli.App{
		Name:  "fetch-competitions",
		Usage: "Fetch and persist competitions for enabled countries from API Football",
		Action: func(*cli.Context) error {
			apiClient = client.NewClient(util.GetEnv("API_FOOTBALL_KEY"), client.RateLimiterSettings{})
			dbGorm = db.InitDbOrPanic()

			fetchAndPersistCompetitions()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func fetchAndPersistCompetitions() {
	var errorCount int
	var createdCount int

	var countries []model.Country
	dbGorm.Where("enabled = ?", 1).Find(&countries)

	fmt.Printf("Countries Enabled: %d\n", len(countries))

	// Update World to be the first country to be processed to not overwrite others
	for i, country := range countries {
		if country.Name == manager.WorldName {
			countries[0], countries[i] = countries[i], countries[0]
			break
		}
	}

	for _, country := range countries {
		getCompetitionsRequest := league.League{CountryCode: country.Code}

		competitionResponses, err := apiClient.DoRequest(getCompetitionsRequest)
		if err != nil {
			fmt.Println(err)
		}

		for _, competitionResponse := range competitionResponses {
			competition := transformer.CreateCompetitionFromResponse(competitionResponse.(leagues.LeagueData), country.ID)
			result := dbGorm.Where("id = ?", competition.ID).Assign(competition).FirstOrCreate(&competition)
			if result.Error != nil {
				errorCount++
			}
			if result.RowsAffected > 0 {
				createdCount++
			}
		}
		fmt.Printf("%s competitions Created: %d, Error: %d\n", country.Name, createdCount, errorCount)
		errorCount, createdCount = 0, 0
	}
}
