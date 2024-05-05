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
	"England":     "\U0001F3F4",
	"Germany":     "\U0001F1E9\U0001F1EA",
	"Spain":       "\U0001F1EA\U0001F1F8",
	"France":      "\U0001F1EB\U0001F1F7",
	"Italy":       "\U0001F1EE\U0001F1F9",
	"UEFA & FIFA": "\U0001F30D",
}

func ToggleUserCountrySettings(user *bot.User, countryID int) {
	disabledCountryIds := user.GetDisabledCountries()
	if sliceContains(disabledCountryIds, countryID) {
		user.SetDisabledCountries(removeElement(disabledCountryIds, countryID))
	} else {
		user.SetDisabledCountries(append(disabledCountryIds, countryID))
	}

	dbGorm.Save(user)
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
			UserDisabled: sliceContains(user.GetDisabledCountries(), int(country.ID)),
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
			UserDisabled: sliceContains(user.GetDisabledCompetitions(), int(competition.ID)),
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

func removeElement(slice []int, element int) []int {
	for i, value := range slice {
		if value == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
