package template

import (
	"fmt"
	model "schedule/model/bot"
	"schedule/tgbot/manager"
	"schedule/util"
	"strconv"
	"strings"

	"github.com/go-telegram/bot/models"
)

const (
	CbdSettingsCountryToggle     = "settings_country_toggle_"
	CbdSettingsCompetitionToggle = "settings_competition_toggle_"
	CbdSettingsCountry           = "settings_country"
	CbdSettingsCompetition       = "settings_competition"
	CbdSettingsAlert             = "settings_alert"
	CbdSettingsUserAlertOffset   = "settings_user_alert_offset"
	CbdSettingsUser              = "settings_user"
	CbdSettingsTimezone          = "settings_timezone"
	CbdSettingsTimezoneInput     = "settings_timezone_input"
	CbdSettingsTimezoneLocation  = "settings_timezone_location"
	CbdSettings                  = "settings"
	CbdSchedule                  = "schedule"
	CbdShowStandings             = "standings_"
	CbdSetLang                   = "set_lang_"
	CbdFixtureToggle             = "fixture_toggle_"

	keyboardButtonTextLength = 50
)

var TimeFormat = "Mon 2.01 15:04"

var LanguageSelectKeyboard = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "🇬🇧 English", CallbackData: CbdSetLang + "en"},
			{Text: "🇷🇺 Русский", CallbackData: CbdSetLang + "ru"},
			{Text: "🇩🇪 Deutsch", CallbackData: CbdSetLang + "de"},
		},
	},
}

var KeyboardSettingsGeneral = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "SettingsCountry", CallbackData: CbdSettingsCountry},
			{Text: "SettingsCompetition", CallbackData: CbdSettingsCompetition},
		},
		{
			{Text: "SettingsAlert", CallbackData: CbdSettingsAlert},
			{Text: "SettingsUser", CallbackData: CbdSettingsUser},
		},
		{
			ButtonSchedule,
		},
	},
}

var KeyboardSettingsUser = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "SettingsTimezone", CallbackData: CbdSettingsTimezone},
			{Text: "SettingsAlertOffset", CallbackData: CbdSettingsUserAlertOffset},
		},
		{
			{Text: "Back", CallbackData: CbdSettings},
		},
	},
}

var KeyboardBack = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "Back", CallbackData: CbdSettings},
		},
	},
}

var KeyboardToSchedule = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
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

var ButtonRefreshSchedule = models.InlineKeyboardButton{
	Text:         "RefreshSchedule",
	CallbackData: CbdSchedule,
}

func GetUserSettingsKeyboardForUser(user model.User) *models.InlineKeyboardMarkup {
	keyboard := KeyboardSettingsUser

	return TranslateKeyboardForUser(user, keyboard)
}

func AppendTranslatedButtonToKeyboard(keyboard *models.InlineKeyboardMarkup, button models.InlineKeyboardButton, user model.User) {
	AppendButtonToKeyboard(keyboard, translateButtonForUser(user, button))
}

func AppendButtonToKeyboard(keyboard *models.InlineKeyboardMarkup, button models.InlineKeyboardButton) {
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{button})
}

func GetCompetitionFixturesKeyboardForUser(user model.User, compView manager.CompetitionView) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{}

	if len(compView.Standings) > 0 {
		standingsButton := models.InlineKeyboardButton{
			Text:         "ShowStandings",
			CallbackData: fmt.Sprintf("%s%d", CbdShowStandings, compView.CompId),
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{translateButtonForUser(user, standingsButton)})
	}

	for _, fixture := range compView.Fixtures {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         generateFixtureButtonText(user, fixture),
				CallbackData: fmt.Sprintf("%s%d", CbdFixtureToggle, fixture.ID),
			},
		})
	}

	return keyboard
}

func ToggleFixtureOnCachedKeyboard(user model.User, fixture manager.FixtureView, originalKeyboard models.InlineKeyboardMarkup) *models.InlineKeyboardMarkup {
	keyboard := originalKeyboard
	for i, row := range keyboard.InlineKeyboard {
		for j, button := range row {
			if checkButtonBelongsToFixture(button, fixture) {
				newText := generateFixtureButtonText(user, fixture)
				if newText == button.Text {
					fixture.IsToggled = !fixture.IsToggled
					newText = generateFixtureButtonText(user, fixture)
				}
				keyboard.InlineKeyboard[i][j].Text = newText
				break
			}
		}
	}

	return &keyboard
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
	icon, score := generateIconAndScore(fixture, user)

	buttonText := fmt.Sprintf(
		"%s %s %s %s %s",
		icon,
		fixture.Date.Format(TimeFormat),
		fixture.HomeTeamName,
		score,
		fixture.AwayTeamName,
	)

	if len(buttonText) > keyboardButtonTextLength && fixture.AwayTeamCode != "" {
		buttonText = fmt.Sprintf(
			"%s %s %s %s %s",
			icon,
			fixture.Date.Format(TimeFormat),
			fixture.HomeTeamName,
			score,
			fixture.AwayTeamCode,
		)
	}
	if len(buttonText) > keyboardButtonTextLength && fixture.HomeTeamCode != "" {
		buttonText = fmt.Sprintf(
			"%s %s %s %s %s",
			icon,
			fixture.Date.Format(TimeFormat),
			fixture.HomeTeamCode,
			score,
			fixture.AwayTeamCode,
		)
	}

	return buttonText
}

func generateIconAndScore(fixture manager.FixtureView, user model.User) (string, string) {
	var icon string
	var toggleIcon string
	score := fixture.Score
	toggleScore := score

	if fixture.Status.IsFinished() {
		if user.EnableSpoilers {
			icon = "🙉"
			toggleIcon = "✔️"
			score = "? : ?"
		} else {
			icon = "✔️"
			toggleIcon = "🙉"
			toggleScore = "? : ?"
		}
	} else {
		if fixture.HasAlert {
			icon = "🔕"
			toggleIcon = "🔔"
		} else {
			icon = "🔔"
			toggleIcon = "🔕"
		}
	}
	if fixture.IsToggled {
		icon = toggleIcon
		score = toggleScore
	}

	return icon, score
}

func checkButtonBelongsToFixture(button models.InlineKeyboardButton, fixture manager.FixtureView) bool {
	return strings.HasSuffix(button.CallbackData, strconv.Itoa(fixture.ID))
}
