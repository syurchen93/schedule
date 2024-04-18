package main

import (
	"github.com/syurchen93/api-football-client/client"
	"github.com/syurchen93/api-football-client/request/league"
	"github.com/syurchen93/api-football-client/response/leagues"

	"fmt"
	"gorm.io/gorm"

	"schedule/db"
	model "schedule/model/league"
	"schedule/util"
	"schedule/util/transformer"
)

var apiClient *client.Client
var dbGorm *gorm.DB

func main() {
	apiClient = client.NewClient(util.GetEnv("API_FOOTBALL_KEY"))
	db.Init()
	dbGorm = db.Db()

	fetchAndPersistCompetitions()
}

func fetchAndPersistCompetitions() {
	var errorCount int
	var createdCount int

	var countries []model.Country
	dbGorm.Where("enabled = ?", 1).Find(&countries)

	fmt.Printf("Countries Enabled: %d\n", len(countries))

	for _, country := range countries {
		getCompetitionsRequest := league.League{CountryCode: country.Code}

		competitionResponses, err := apiClient.DoRequest(getCompetitionsRequest)
		if err != nil {
			fmt.Println(err)
		}

		for _, competitionResponse := range competitionResponses {
			competition := transformer.CreateCompetitionFromResponse(competitionResponse.(leagues.LeagueData), country.ID)
			result := dbGorm.Where(model.Competition{ID: competition.ID}).Assign(competition).FirstOrCreate(&competition)
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
