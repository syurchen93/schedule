package main

import (
	"fmt"

	"github.com/syurchen93/api-football-client/client"
	"github.com/syurchen93/api-football-client/request/team"
	response "github.com/syurchen93/api-football-client/response/team"

	"log"
	"os"
	"time"

	"gorm.io/gorm"

	"schedule/db"
	"schedule/model/league"
	"schedule/util"
	"schedule/util/transformer"

	"github.com/urfave/cli/v2"
)

var apiClient *client.Client
var dbGorm *gorm.DB

func main() {
	app := &cli.App{
		Name:  "fetch-teams",
		Usage: "Fetch and persist missing team data.",
		Action: func(*cli.Context) error {
			apiClient = client.NewClient(util.GetEnv("API_FOOTBALL_KEY"), client.RateLimiterSettings{})
			db.Init()
			dbGorm = db.Db()

			fetchAndPersistTeamsWithMissingData()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func fetchAndPersistTeamsWithMissingData() {
	var errorCount int
	var teamCreatedCount int
	var competitionCount int

	for {
		competition := fetchCompetitionWithMostTeamsWithMissingData()
		if competition == nil {
			break
		}
		competitionCount++
		teams, err := apiClient.DoRequest(team.Team{
			League: int(competition.ID),
			Season: int(competition.CurrentSeason),
		})
		if err != nil {
			errorCount++
			log.Printf("Error fetching teams for competition %s: %s", competition.Name, err)
			break
		}
		for _, team := range teams {
			teamModel := transformer.CreateTeamFromTeamInformation(team.(response.Information))
			teamModel.TriedToFetchCodeAt = time.Now()
			result := dbGorm.Where("id = ?", teamModel.ID).Assign(teamModel).FirstOrCreate(&teamModel)
			if result.Error != nil {
				errorCount++
				continue
			}
			if result.RowsAffected == 0 {
				teamCreatedCount++
			}
		}
	}
	fmt.Printf("Competitions processed: %d, Teams created: %d, Errors: %d\n", competitionCount, teamCreatedCount, errorCount)
}

func fetchCompetitionWithMostTeamsWithMissingData() *league.Competition {
	var competition league.Competition

	err := dbGorm.Model(&competition).
		Select("competition.*, COUNT(DISTINCT team.id) as team_count").
		Joins("JOIN fixture ON fixture.competition_id = competition.id").
		Joins("JOIN team ON team.id = fixture.home_team_id OR team.id = fixture.away_team_id").
		Where("(team.code = '' OR team.country = '') AND team.tried_to_fetch_code_at IS NULL").
		Group("competition.id").
		Order("team_count DESC").
		First(&competition).Error

	if err != nil {
		return nil
	}

	return &competition
}
