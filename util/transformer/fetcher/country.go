package fetcher

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
