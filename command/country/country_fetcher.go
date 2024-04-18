package main

import (
	"github.com/syurchen93/api-football-client/client"
	"github.com/syurchen93/api-football-client/request/league"
	"github.com/syurchen93/api-football-client/response/leagues"

	"fmt"
	"gorm.io/gorm"

	"schedule/db"
	"schedule/util"
	"schedule/util/transformer"
)

var apiClient *client.Client
var dbGorm *gorm.DB

func main() {
	apiClient = client.NewClient(util.GetEnv("API_FOOTBALL_KEY"))
	db.Init()
	dbGorm = db.Db()

	fetchAndPersistCountries()
}

func fetchAndPersistCountries() {
	var errorCount int
	var createdCount int

	getCountriesRequest := league.Country{}

	countryResponses, err := apiClient.DoRequest(getCountriesRequest)
	if err != nil {
		panic(err)
	}

	for _, countryResponse := range countryResponses {
		country := transformer.CreateCountryFromResponse(countryResponse.(leagues.Country))
		result := dbGorm.Where("name = ?", country.Name).Assign(country).FirstOrCreate(&country)
		if result.Error != nil {
			errorCount++
		}
		if result.RowsAffected > 0 {
			createdCount++
		}
	}

	fmt.Printf("Countries Created: %d, Error: %d\n", createdCount, errorCount)
}
