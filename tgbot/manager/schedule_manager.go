package manager

import (
	"fmt"
	"schedule/model/bot"
	"schedule/model/league"
	"sort"
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
	TeamId      int
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
	ID            int
	HomeTeamName  string
	HomeTeamCode  string
	IsHomeUserFav bool
	AwayTeamName  string
	AwayTeamCode  string
	IsAwayUserFav bool
	Date          time.Time
	Score         string
	Status        common.FixtureStatus
	HasAlert      bool
	IsToggled     bool
}

func (s StandingView) GetTeamNameWithCode() string {
	if s.TeamCode != "" {
		return fmt.Sprintf("%s (%s)", s.TeamName, s.TeamCode)
	}
	return s.TeamName
}

const DefaultDaysInFuture = 7
const DefaultDaysInPast = 7

func GetCompetitionViewsForUser(user *bot.User) []CompetitionView {
	var fixturesByComp []CompetitionView
	fixtures := getHydratedFixturesForUser(user)
	for _, fixture := range fixtures {
		var compFound bool

		fixtureView := createFixtureView(fixture, user)
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
	var userFavTeamIDs []int

	dbGorm.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Competition").
		Preload("Competition.Country").
		First(&fixture, fixtureId)

	if fixture.ID == 0 {
		panic(fmt.Errorf("fixture with id %d not found", fixtureId))
	}

	view := createFixtureView(fixture, user)
	dbGorm.Table("fav_team").Select("team_id").Where("user_id = ?", user.ID).Scan(&userFavTeamIDs)

	view.IsHomeUserFav = contains(userFavTeamIDs, fixture.HomeTeam.ID)
	view.IsAwayUserFav = contains(userFavTeamIDs, fixture.AwayTeam.ID)
	view.IsToggled = true
	toggleFixtureViewAlertIfNeeded(user, view)

	return view
}

func CreateCompetitionFixtureViewFromAlers(alerts []bot.Alert) []CompetitionView {
	var compViews []CompetitionView

	for _, alert := range alerts {
		fixture := alert.Fixture
		fixture.HasUserAlert = true
		var compFound bool
		fixtureView := createFixtureView(fixture, &alert.User)
		for i, comp := range compViews {
			if comp.CompId == fixture.CompetitionID {
				compViews[i].Fixtures = append(compViews[i].Fixtures, fixtureView)
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
			compViews = append(compViews, compView)
		}
	}

	return compViews
}

func createFixtureView(fixture league.Fixture, user *bot.User) FixtureView {
	var homeTeamCode, awayTeamCode *string

	userTime, err := time.LoadLocation(user.Timezone)
	if err != nil {
		userTime = time.UTC
	}
	if fixture.HomeTeam.Code != nil {
		homeTeamCode = fixture.HomeTeam.Code
	}
	if fixture.AwayTeam.Code != nil {
		awayTeamCode = fixture.AwayTeam.Code
	}

	return FixtureView{
		ID:            fixture.ID,
		HomeTeamName:  fixture.HomeTeam.Name,
		HomeTeamCode:  *homeTeamCode,
		IsHomeUserFav: fixture.HomeTeam.IsUserFav,
		AwayTeamName:  fixture.AwayTeam.Name,
		AwayTeamCode:  *awayTeamCode,
		IsAwayUserFav: fixture.AwayTeam.IsUserFav,
		Date:          fixture.Date.In(userTime),
		Score:         generateScoreString(fixture),
		Status:        fixture.Status,
		HasAlert:      fixture.HasUserAlert,
	}
}

func fetchUpToDateCompetitionStandings(competitionId uint) []league.Standing {
	var standings []league.Standing
	dbGorm.
		Table("standing as s").
		Joins("join competition c on c.id = s.competition_id and c.current_season = s.season").
		Joins("join team on team.id = s.team_id").
		Preload("Team").
		Where("s.competition_id = ?", competitionId).
		Order("s.rank ASC").
		Find(&standings)

	return standings
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
	var results []struct {
		league.Fixture
		HomeIsUserFav bool `gorm:"column:home_is_user_fav"`
		AwayIsUserFav bool `gorm:"column:away_is_user_fav"`
	}

	query := dbGorm.Table("fixture").Joins("left join competition on competition.id = fixture.competition_id").
		Joins("left join country on country.id = competition.country_id").
		Joins("left join alert on alert.fixture_id = fixture.id AND alert.user_id = ?", user.ID).
		Joins("left join fav_team as home_fav on home_fav.team_id = fixture.home_team_id AND home_fav.user_id = ?", user.ID).
		Joins("left join fav_team as away_fav on away_fav.team_id = fixture.away_team_id AND away_fav.user_id = ?", user.ID).
		Where("fixture.date > ?", time.Now().AddDate(0, 0, -DefaultDaysInPast)).
		Where("fixture.date < ?", time.Now().AddDate(0, 0, DefaultDaysInFuture)).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Competition").
		Preload("Competition.Country").
		Select("fixture.*, CASE WHEN alert.id IS NOT NULL THEN true ELSE false END AS has_user_alert, CASE WHEN home_fav.id IS NOT NULL THEN true ELSE false END AS home_is_user_fav, CASE WHEN away_fav.id IS NOT NULL THEN true ELSE false END AS away_is_user_fav").
		Order("fixture.date ASC")

	if len(user.GetDisabledCountries()) > 0 {
		query = query.Not("country_id", user.GetDisabledCountries())
	}

	if len(user.GetDisabledCompetitions()) > 0 {
		query = query.Not("competition_id", user.GetDisabledCompetitions())
	}

	query.Find(&results)

	for _, result := range results {
		result.HomeTeam.IsUserFav = result.HomeIsUserFav
		result.AwayTeam.IsUserFav = result.AwayIsUserFav
		fixtures = append(fixtures, result.Fixture)
	}

	return fixtures
}

func toggleFixtureViewAlertIfNeeded(user *bot.User, fixture FixtureView) {
	if !fixture.Status.IsFinished() {
		createOrDeleteAlertForFixture(user, fixture.ID)
	}
}

func buildStandingDatas(standings []league.Standing) []StandingsData {
	var standingsData []StandingsData
	groupedStandings := make(map[string][]StandingView)
	for _, standing := range standings {
		standingsView := StandingView{
			TeamId:      standing.Team.ID,
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

	var groupNames []string
	for groupName := range groupedStandings {
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)

	for _, groupName := range groupNames {
		groupStandings := groupedStandings[groupName]

		sort.Slice(groupStandings, func(i, j int) bool {
			return groupStandings[i].Position < groupStandings[j].Position
		})

		standingData := StandingsData{
			GroupName: groupName,
			Standings: groupStandings,
		}

		standingsData = append(standingsData, standingData)
	}

	return standingsData
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
