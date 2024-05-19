package manager

import (
	"fmt"
	"schedule/model/bot"
	"schedule/model/league"
	"time"

	"github.com/syurchen93/api-football-client/common"
)

type CompetitionView struct {
	CompId      uint
	CompName    string
	CountryName string
	Fixtures    []FixtureView
	Standings   []StandingsData
}

type StandingsData struct {
	GroupName string
	Standings []StandingView
}

type StandingView struct {
	TeamName    string
	TeamCode    string
	Position    int
	Points      int
	GoalsDiff   int
	Form        string
	Played      int
	Won         int
	Drawn       int
	Lost        int
	Description string
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

func GetCompetitionViewsForUser(user *bot.User) []CompetitionView {
	var fixturesByComp []CompetitionView
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
			compView := CompetitionView{
				CompId:      fixture.CompetitionID,
				CompName:    fixture.Competition.Name,
				CountryName: fixture.Competition.Country.Name,
				Fixtures: []FixtureView{
					fixtureView,
				},
			}
			if !fixture.Competition.NoStandings {
				compView.Standings = GetCachedCompetitionStandings(fixture.CompetitionID)
			}
			fixturesByComp = append(fixturesByComp, compView)
		}
	}

	return fixturesByComp
}

func GetCompetitionFixturesAndToggleByFixtureId(user *bot.User, fixtureId int) CompetitionView {
	var wantedComp CompetitionView
	var wantedFixture FixtureView
	competitionFixtures := GetCompetitionViewsForUser(user)
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
	toggleFixtureViewAlertIfNeeded(user, wantedFixture)

	return wantedComp
}

func GetToggleFixtureViewByFixtureId(user *bot.User, fixtureId int) FixtureView {
	var fixture league.Fixture
	dbGorm.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Competition").
		Preload("Competition.Country").
		First(&fixture, fixtureId)

	if fixture.ID == 0 {
		panic(fmt.Errorf("fixture with id %d not found", fixtureId))
	}
	view := createFixtureView(fixture)
	view.IsToggled = true
	toggleFixtureViewAlertIfNeeded(user, view)

	return view
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

func toggleFixtureViewAlertIfNeeded(user *bot.User, fixture FixtureView) {
	if !fixture.Status.IsFinished() {
		createOrDeleteAlertForFixture(user, fixture.ID)
	}
}

func getCompetitionStandings(standingsData *[]StandingsData, competitionId uint) {
	var standings []league.Standing
	dbGorm.
		Table("standing as s").
		Select("s.*, team.code").
		Joins("join competition c on c.id = s.competition_id and c.current_season = s.season").
		Joins("join team on team.id = s.team_id").
		Joins("join ("+
			"select `rank`, `group`, max(id) as max_id "+
			"from standing "+
			"group by `rank`, `group`"+
			") as latest on latest.`rank` = s.`rank` and latest.`group` = s.`group` and latest.max_id = s.id").
		Preload("Team").
		Where("s.competition_id = ?", competitionId).
		Order("s.rank ASC").
		Find(&standings)

	groupedStandings := make(map[string][]StandingView)
	for _, standing := range standings {
		standingsView := StandingView{
			TeamName:    standing.Team.Name,
			TeamCode:    *standing.Team.Code,
			Position:    standing.Rank,
			Points:      standing.Points,
			GoalsDiff:   standing.GoalsDiff,
			Form:        standing.Form,
			Played:      standing.Played,
			Won:         standing.Won,
			Drawn:       standing.Drawn,
			Lost:        standing.Lost,
			Description: standing.Description,
		}

		groupedStandings[standing.Group] = append(groupedStandings[standing.Group], standingsView)
	}

	for groupName, groupStandings := range groupedStandings {
		standingData := StandingsData{
			GroupName: groupName,
			Standings: groupStandings,
		}

		*standingsData = append(*standingsData, standingData)
	}
}
