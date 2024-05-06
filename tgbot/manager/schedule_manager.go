package manager

import (
	"fmt"
	"schedule/model/bot"
	"schedule/model/league"
	"time"

	"github.com/syurchen93/api-football-client/common"
)

type CompetitionFixtures struct {
	CompId      uint
	CompName    string
	CountryName string
	Fixtures    []FixtureView
}

type FixtureView struct {
	ID           int
	HomeTeamName string
	AwayTeamName string
	Date         time.Time
	Score        string
	HasAlert     bool
}

func GetCompetitionFixturesForUser(user *bot.User) []CompetitionFixtures {
	var fixturesByComp []CompetitionFixtures
	fixtures := getHydratedFixturesForUser(user)
	for _, fixture := range fixtures {
		var compFound bool

		fixtureView := createFixtureView(fixture)
		for i, comp := range fixturesByComp {
			if comp.CompId == fixture.CompetitionID {
				fixturesByComp[i].Fixtures = append(fixturesByComp[i].Fixtures, fixtureView)
				compFound = true
				break
			}
		}
		if !compFound {
			fixturesByComp = append(fixturesByComp, CompetitionFixtures{
				CompId:      fixture.CompetitionID,
				CompName:    fixture.Competition.Name,
				CountryName: fixture.Competition.Country.Name,
				Fixtures: []FixtureView{
					fixtureView,
				},
			})
		}
	}

	return fixturesByComp
}

func createFixtureView(fixture league.Fixture) FixtureView {
	return FixtureView{
		ID:           fixture.ID,
		HomeTeamName: fixture.HomeTeam.Name,
		AwayTeamName: fixture.AwayTeam.Name,
		Date:         fixture.Date,
		Score:        generateScoreString(fixture),
	}
}

func generateScoreString(fixture league.Fixture) string {
	var scoreString string
	if !fixture.Status.IsFinished() {
		return scoreString
	}

	scoreString = fmt.Sprintf("%d : %d", fixture.GoalsHome, fixture.GoalsAway)

	if fixture.Status == common.Finished {
		return scoreString
	}

	if fixture.Status == common.FinishedAfterExtra {
		return scoreString + " (ET)"
	}

	if fixture.Status == common.FinishedAfterPenalty {
		return fmt.Sprintf("%s P(%d : %d)", scoreString, fixture.PenaltyHome, fixture.PenaltyAway)
	}

	return scoreString
}

func getHydratedFixturesForUser(user *bot.User) []league.Fixture {
	var fixtures []league.Fixture
	dbGorm.Joins("left join competition on competition.id = fixture.competition_id").
		Joins("left join country on country.id = competition.country_id").
		Where("fixture.date > ?", time.Now()).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Competition").
		Preload("Competition.Country")

	if len(user.GetDisabledCountries()) > 0 {
		dbGorm = dbGorm.Not("country.id", user.GetDisabledCountries())
	}

	if len(user.GetDisabledCompetitions()) > 0 {
		dbGorm = dbGorm.Not("competition.id", user.GetDisabledCompetitions())
	}

	dbGorm.Find(&fixtures)
	return fixtures
}
