package transformer

import (
	"schedule/model"
	leaguem "schedule/model/league"

	"github.com/syurchen93/api-football-client/response/leagues"
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

func CreateTeamFromResponse(response team.Team, competitionID uint) model.Team {
	return model.Team{
		ID:      response.ID,
		Name:    response.Name,
		Code:    response.Code,
		Country: response.Country,
		Logo:    response.Logo,
	}
}
