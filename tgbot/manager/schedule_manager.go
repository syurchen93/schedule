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
	HomeTeamCode string
	AwayTeamName string
	AwayTeamCode string
	Date         time.Time
	Score        string
	Status       common.FixtureStatus
	HasAlert     bool
	IsToggled    bool
}

const DefaultDaysInFuture = 7
const DefaultDaysInPast = 7

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

func GetCompetitionFixturesAndToggleByFixtureId(user *bot.User, fixtureId int) CompetitionFixtures {
	var wantedComp CompetitionFixtures
	var wantedFixture FixtureView
	competitionFixtures := GetCompetitionFixturesForUser(user)
	for _, comp := range competitionFixtures {
		for i, fixture := range comp.Fixtures {
			if fixture.ID == fixtureId {
				wantedFixture = fixture
				wantedComp = comp
				wantedComp.Fixtures[i].IsToggled = true
				break
			}
		}
	}
	if wantedFixture.ID == 0 {
		panic(fmt.Errorf("fixture with id %d not found", fixtureId))
	}
	if !wantedFixture.Status.IsFinished() {
		createOrDeleteAlertForFixture(user, wantedFixture.ID)
	}

	return wantedComp
}

func createFixtureView(fixture league.Fixture) FixtureView {
	return FixtureView{
		ID:           fixture.ID,
		HomeTeamName: fixture.HomeTeam.Name,
		AwayTeamName: fixture.AwayTeam.Name,
		Date:         fixture.Date,
		Score:        generateScoreString(fixture),
		Status:       fixture.Status,
		HasAlert:     fixture.HasUserAlert,
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
	subquery := dbGorm.Table("alert").Select("1").Where("alert.fixture_id = fixture.id AND alert.user_id = ?", user.ID)
	query := dbGorm.Joins("left join competition on competition.id = fixture.competition_id").
		Joins("left join country on country.id = competition.country_id").
		Where("fixture.date > ?", time.Now().AddDate(0, 0, -DefaultDaysInPast)).
		Where("fixture.date < ?", time.Now().AddDate(0, 0, DefaultDaysInFuture)).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Competition").
		Preload("Competition.Country").
		Select("fixture.*, CASE WHEN EXISTS(?) THEN true ELSE false END AS has_user_alert", subquery).
		Order("fixture.date ASC")
	if len(user.GetDisabledCountries()) > 0 {
		query = query.Not("country_id", user.GetDisabledCountries())
	}

	if len(user.GetDisabledCompetitions()) > 0 {
		query = query.Not("competition_id", user.GetDisabledCompetitions())
	}

	query.Find(&fixtures)
	return fixtures
}
