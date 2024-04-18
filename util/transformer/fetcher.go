package transformer

import (
	"github.com/syurchen93/api-football-client/response/leagues"
	model "schedule/model/league"
)

func CreateCountryFromResponse(response leagues.Country) model.Country {
	return model.Country{
		Name: response.Name,
		Code: response.Code,
		Flag: response.Flag,
	}
}

func CreateCompetitionFromResponse(response leagues.LeagueData, countryID uint) model.Competition {
	return model.Competition{
		CountryID: countryID,
		ID:        uint(response.League.ID),
		Name:      response.League.Name,
		Type:      response.League.Type,
		Logo:      response.League.Logo,
	}
}
