package main

import (
	"github.com/syurchen93/api-football-client/client"
	"github.com/syurchen93/api-football-client/request/league"

	"schedule/util"
)

func main() {
	apiClient := client.NewClient(util.GetEnv("API_FOOTBALL_KEY"))

	getCountriesRequest := league.Country{}

	apiClient.DoRequest(getCountriesRequest)
}
