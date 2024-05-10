package template

import (
	"fmt"
	model "schedule/model/bot"
	"schedule/tgbot/manager"
	"schedule/util"
	"strings"

	"github.com/go-telegram/bot/models"
)

const (
	CbdSettingsCountryToggle     = "settings_country_toggle_"
	CbdSettingsCompetitionToggle = "settings_competition_toggle_"
	CbdSettingsCountry           = "settings_country"
	CbdSettingsCompetition       = "settings_competition"
	CbdSettingsAlert             = "settings_alert"
	CbdSettings                  = "settings"
	CbdSchedule                  = "schedule"
	CbdSetLang                   = "set_lang_"
	CbdToggleAlert               = "alert_toggle_"

	keyboardButtonTextLength = 30
)

var TimeFormat = "Mon 2.01 15:04"

var LanguageSelectKeyboard = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "🇬🇧 English", CallbackData: CbdSetLang + "_en"},
			{Text: "🇷🇺 Русский", CallbackData: CbdSetLang + "set_lang_ru"},
			{Text: "🇩🇪 Deutsch", CallbackData: CbdSetLang + "set_lang_de"},
		},
	},
}

var KeyboardSettingsGeneral = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "SettingsCountry", CallbackData: CbdSettingsCountry},
			{Text: "SettingsCompetition", CallbackData: CbdSettingsCompetition},
			{Text: "SettingsAlert", CallbackData: CbdSettingsAlert},
		},
		{
			ButtonSchedule,
		},
	},
}

var ButtonSettings = models.InlineKeyboardButton{
	Text:         "ToSettings",
	CallbackData: CbdSettings,
}

var ButtonSchedule = models.InlineKeyboardButton{
	Text:         "ToSchedule",
	CallbackData: CbdSchedule,
}

func GetFixturesKeyboardForUser(user model.User, fixtures []manager.FixtureView) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{}
	for _, fixture := range fixtures {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         generateFixtureButtonText(user, fixture),
				CallbackData: fmt.Sprintf("%s%d", CbdToggleAlert, fixture.ID),
			},
		})
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
		translateButtonForUser(user, ButtonSettings),
	})

	return keyboard
}

func GetUserCompetitonSettingsKyboard(user *model.User, compSettings []manager.CompetitionSettings) *models.InlineKeyboardMarkup {
	var toggleTextKey string
	keyboard := &models.InlineKeyboardMarkup{}
	for _, comp := range compSettings {
		if comp.UserDisabled {
			toggleTextKey = "Enable"
		} else {
			toggleTextKey = "Disable"
		}

		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text: fmt.Sprintf(
					"%s %s",
					util.Translate(user.Locale, toggleTextKey),
					comp.Name,
				),
				CallbackData: fmt.Sprintf("%s%d", CbdSettingsCompetitionToggle, comp.ID),
			},
		})
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
		{
			Text:         util.Translate(user.Locale, "Back"),
			CallbackData: CbdSettings,
		},
	})

	return keyboard
}

func GetUserCountrySettingsKyboard(user *model.User, countrySettings []manager.CountrySettings) *models.InlineKeyboardMarkup {
	var toggleTextKey string
	keyboard := &models.InlineKeyboardMarkup{}
	for _, country := range countrySettings {
		if country.UserDisabled {
			toggleTextKey = "Enable"
		} else {
			toggleTextKey = "Disable"
		}

		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text: fmt.Sprintf(
					"%s %s %s",
					util.Translate(user.Locale, toggleTextKey),
					country.Emoji,
					country.Name,
				),
				CallbackData: fmt.Sprintf("%s%d", CbdSettingsCountryToggle, country.ID),
			},
		})
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
		{
			Text:         util.Translate(user.Locale, "Back"),
			CallbackData: "settings",
		},
	})

	return keyboard
}

func TranslateKeyboardForUser(user model.User, keyboard *models.InlineKeyboardMarkup) *models.InlineKeyboardMarkup {
	userKeyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: make([][]models.InlineKeyboardButton, len(keyboard.InlineKeyboard)),
	}

	for i, row := range keyboard.InlineKeyboard {
		userKeyboard.InlineKeyboard[i] = make([]models.InlineKeyboardButton, len(row))
		for j, button := range row {
			userKeyboard.InlineKeyboard[i][j] = translateButtonForUser(user, button)
		}
	}

	return userKeyboard
}

func GetLanguageSelectKeyboardForUser(user model.User) *models.InlineKeyboardMarkup {
	userKeyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: make([][]models.InlineKeyboardButton, len(LanguageSelectKeyboard.InlineKeyboard)),
	}

	for i, row := range LanguageSelectKeyboard.InlineKeyboard {
		userKeyboard.InlineKeyboard[i] = make([]models.InlineKeyboardButton, len(row))
		copy(userKeyboard.InlineKeyboard[i], row)
	}

	for i, button := range userKeyboard.InlineKeyboard[0] {
		if strings.HasSuffix(button.CallbackData, user.Locale) {
			userKeyboard.InlineKeyboard[0] = remove(userKeyboard.InlineKeyboard[0], i)
			break
		}
	}
	userKeyboard.InlineKeyboard = append(userKeyboard.InlineKeyboard,
		[]models.InlineKeyboardButton{translateButtonForUser(user, ButtonSettings)},
		[]models.InlineKeyboardButton{translateButtonForUser(user, ButtonSchedule)},
	)

	return userKeyboard
}

func translateButtonForUser(user model.User, button models.InlineKeyboardButton) models.InlineKeyboardButton {
	button.Text = util.Translate(user.Locale, button.Text)

	return button
}

func remove(slice []models.InlineKeyboardButton, i int) []models.InlineKeyboardButton {
	slice[i] = slice[len(slice)-1]
	slice[len(slice)-1] = models.InlineKeyboardButton{}
	slice = slice[:len(slice)-1]

	return slice
}

func generateFixtureButtonText(user model.User, fixture manager.FixtureView) string {
	score := fixture.Score
	if fixture.Status.IsFinished() && !user.EnableSpoilers {
		score = "🙈 : 🙈"
	}

	buttonText := fmt.Sprintf(
		"%s %s : %s %s",
		fixture.Date.Format(TimeFormat),
		fixture.HomeTeamName,
		score,
		fixture.AwayTeamName,
	)

	if len(buttonText) > keyboardButtonTextLength && fixture.HomeTeamCode != "" && fixture.AwayTeamCode != "" {
		buttonText = fmt.Sprintf(
			"%s %s %s %s",
			fixture.Date.Format(TimeFormat),
			fixture.HomeTeamCode,
			score,
			fixture.AwayTeamCode,
		)
	}

	return buttonText
}
