package transformer

import (
	"schedule/model"
	leaguem "schedule/model/league"

	"github.com/syurchen93/api-football-client/response/fixtures"
	"github.com/syurchen93/api-football-client/response/leagues"
	"github.com/syurchen93/api-football-client/response/standings"
	"github.com/syurchen93/api-football-client/response/team"
)

func CreateCountryFromResponse(response leagues.Country) leaguem.Country {
	return leaguem.Country{
		Name: response.Name,
		Code: response.Code,
		Flag: response.Flag,
	}
}

func CreateCompetitionFromResponse(response leagues.LeagueData, countryID uint) leaguem.Competition {
	var currentSeason leagues.Season
	for _, season := range response.Seasons {
		if season.Current {
			currentSeason = season
			break
		}
	}
	if currentSeason.Year == 0 {
		currentSeason = response.Seasons[len(response.Seasons)-1]
	}

	return leaguem.Competition{
		CountryID:     countryID,
		ID:            uint(response.League.ID),
		Name:          response.League.Name,
		Type:          response.League.Type,
		Logo:          response.League.Logo,
		CurrentSeason: uint(currentSeason.Year),
	}
}

func CreateTeamFromResponse(response team.Team) model.Team {
	return model.Team{
		ID:      response.ID,
		Name:    response.Name,
		Code:    &response.Code,
		Country: &response.Country,
		Logo:    response.Logo,
	}
}

func CreateTeamFromFixtureResponse(response fixtures.Team) model.Team {
	return model.Team{
		ID:   response.ID,
		Name: response.Name,
		Logo: response.Logo,
	}
}

func CreateStandingFromResponse(response standings.Ranking, competitionID uint) leaguem.Standing {
	return leaguem.Standing{
		Rank:          response.Rank,
		TeamID:        uint(response.Team.ID),
		CompetitionID: competitionID,
		Points:        response.Points,
		GoalsDiff:     response.GoalsDiff,
		Group:         response.Group,
		Form:          response.Form,
		Status:        response.Status,
		Description:   response.Description,
		Played:        response.All.Played,
		Won:           response.All.Win,
		Drawn:         response.All.Draw,
		Lost:          response.All.Lose,
		GoalsFor:      response.All.Goals.For,
		GoalsAgainst:  response.All.Goals.Against,
		UpdatedApi:    response.Updated,
	}
}

func CreateFixtureFromResponse(response fixtures.Fixture) leaguem.Fixture {
	fixtureModel := leaguem.Fixture{
		ID:            response.Fixture.ID,
		CompetitionID: uint(response.League.ID),
		HomeTeamID:    uint(response.Teams.Home.ID),
		AwayTeamID:    uint(response.Teams.Away.ID),
		Status:        response.Fixture.Status.Value,
		GoalsHome:     response.Goals.Home,
		GoalsAway:     response.Goals.Away,
		PenaltyHome:   response.Score.Penalty.Home,
		PenaltyAway:   response.Score.Penalty.Away,
		Date:          response.Fixture.Date,
	}

	return fixtureModel
}
