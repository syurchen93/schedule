package manager

import (
	"schedule/model/bot"
	"schedule/model/league"
)

type CountrySettings struct {
	ID           uint
	Name         string
	Emoji        string
	UserDisabled bool
}

type CompetitionSettings struct {
	ID           uint
	Name         string
	UserDisabled bool
}

var CountryEmojiMap = map[string]string{
	"England": "🏴󠁧󠁢󠁥󠁮󠁧󠁿",
	"Germany": "🇩🇪",
	"Spain":   "🇪🇸",
	"France":  "🇫🇷",
	"Italy":   "🇮🇹",
	"World":   "🌍",
}

func GetUserCountrySettings(user *bot.User) []CountrySettings {
	var countrySettings []CountrySettings

	var countries []league.Country
	dbGorm.Where("enabled = ?", 1).Find(&countries)
	for _, country := range countries {
		countrySettings = append(countrySettings, CountrySettings{
			ID:           country.ID,
			Name:         country.Name,
			Emoji:        CountryEmojiMap[country.Name],
			UserDisabled: sliceContains(user.DisabledCountries, int(country.ID)),
		})
	}

	return countrySettings
}

func GetUserCountryCompetitionSettings(user *bot.User, countryID uint) []CompetitionSettings {
	var competitionSettings []CompetitionSettings

	var competitions []league.Competition
	dbGorm.Where("country_id = ?", countryID).Find(&competitions)
	for _, competition := range competitions {
		competitionSettings = append(competitionSettings, CompetitionSettings{
			ID:           competition.ID,
			Name:         competition.Name,
			UserDisabled: sliceContains(user.DisabledCompetitions, int(competition.ID)),
		})
	}

	return competitionSettings
}

func sliceContains(slice []int, element int) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}
