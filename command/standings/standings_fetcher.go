package main

import (
	"fmt"

	"github.com/syurchen93/api-football-client/client"
	"github.com/syurchen93/api-football-client/request/standings"
	response "github.com/syurchen93/api-football-client/response/standings"

	"log"
	"os"

	"gorm.io/gorm"

	"schedule/db"
	model "schedule/model/league"
	"schedule/util"
	"schedule/util/transformer"

	"github.com/urfave/cli/v2"
)

var apiClient *client.Client
var dbGorm *gorm.DB

func main() {
	app := &cli.App{
		Name:  "fetch-standings",
		Usage: "Fetch and persist standings from API Football. Creating teams on the fly if they don't exist.",
		Action: func(*cli.Context) error {
			apiClient = client.NewClient(util.GetEnv("API_FOOTBALL_KEY"), client.RateLimiterSettings{})
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
	var errorCount int
	var teamCreatedCount int
	var standingsCreatedCount int
	var competitionCount int
	var competitions []model.Competition
	dbGorm.Where("enabled = ? and no_standings = 0", true).Find(&competitions)

	for _, competition := range competitions {
		standingsRequest := standings.Standings{
			League: int(competition.ID),
			Season: int(competition.CurrentSeason),
		}

		standingResponse, err := apiClient.DoRequest(standingsRequest)
		if err != nil {
			panic(err)
		}

		var rankings []response.Ranking
		if len(standingResponse) > 0 {
			for _, rankingSlice := range standingResponse[0].(response.Standings).League.Standings {
				rankings = append(rankings, rankingSlice...)
			}
		}

		for _, ranking := range rankings {
			team := transformer.CreateTeamFromResponse(ranking.Team, competition.ID)
			result := dbGorm.Where("id = ?", team.ID).Assign(team).FirstOrCreate(&team)
			if result.Error != nil {
				errorCount++
			}
			if result.RowsAffected > 0 {
				teamCreatedCount++
			}
			standing := transformer.CreateStandingFromResponse(ranking, competition.ID)
			result = dbGorm.Where("team_id = ? AND competition_id = ?", standing.TeamID, standing.CompetitionID).Assign(standing).FirstOrCreate(&standing)
			if result.Error != nil {
				errorCount++
			}
			if result.RowsAffected > 0 {
				standingsCreatedCount++
			}
		}
		competitionCount++
	}
	fmt.Printf("Created %d teams, %d standings in %d competitions with %d errors\n", teamCreatedCount, standingsCreatedCount, competitionCount, errorCount)
}
