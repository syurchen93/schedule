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
	CbdSettingsFavTeam           = "settings_fav_team"
	CbdSettingsFavTeamRemove     = "settings_fav_team_remove_"
	CbdSettingsFavTeamAddStart   = "settings_fav_team_add_start"
	CbdSettingsTimezone          = "settings_timezone"
	CbdSettingsTimezoneInput     = "settings_timezone_input"
	CbdSettingsTimezoneLocation  = "settings_timezone_location"
	CbdSettingsShareManage       = "settings_share_manage"
	CbdSettingsShareRemove       = "settings_share_remove"
	CbdSettings                  = "settings"
	CbdSchedule                  = "schedule"
	CbdShowStandings             = "standings_"
	CbdSetLang                   = "set_lang_"
	CbdFixtureToggle             = "fixture_toggle_"

	keyboardButtonTextLength = 50
	keyboardMaxButtonCount   = 19
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
			{Text: "SettingsFavTeam", CallbackData: CbdSettingsFavTeam},
			{Text: "SettingsShare", CallbackData: CbdSettingsShareManage},
		},
		{
			ButtonBack,
		},
	},
}

var KeyboardBack = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			ButtonBack,
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

var ButtonBack = models.InlineKeyboardButton{
	Text:         "Back",
	CallbackData: CbdSettings,
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

func GetShareAlertsKeyboardForUser(userShares []model.UserShare, user model.User) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{}

	for _, userShare := range userShares {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("❌ %s", userShare.SourceUser.Username),
				CallbackData: fmt.Sprintf("%s%d", CbdSettingsShareRemove, userShare.ID),
			},
		})
	}

	AppendTranslatedButtonToKeyboard(keyboard, ButtonBack, user)

	return keyboard
}

func RemoveFavTeamFromCachedKeyboard(favTeamId int, originalKeyboard models.InlineKeyboardMarkup) *models.InlineKeyboardMarkup {
	keyboard := originalKeyboard
	for i, row := range keyboard.InlineKeyboard {
		for j, button := range row {
			if strings.HasSuffix(button.CallbackData, strconv.Itoa(favTeamId)) {
				keyboard.InlineKeyboard[i] = removeButtonFromBlock(keyboard.InlineKeyboard[i], j)
				break
			}
		}
	}

	return &keyboard
}

func GetFavTeamKeyboardForUser(favTeams []model.FavTeam, user model.User) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{}

	for _, favTeam := range favTeams {
		AppendButtonToKeyboard(
			keyboard,
			models.InlineKeyboardButton{
				Text:         fmt.Sprintf("%s %s", "❌", favTeam.Team.Name),
				CallbackData: fmt.Sprintf("%s%d", CbdSettingsFavTeamRemove, favTeam.Team.ID),
			},
		)
	}

	buttonAdd := models.InlineKeyboardButton{
		Text:         "FavTeamAddStart",
		CallbackData: CbdSettingsFavTeamAddStart,
	}

	AppendTranslatedButtonToKeyboard(keyboard, buttonAdd, user)
	AppendTranslatedButtonToKeyboard(keyboard, ButtonBack, user)

	return keyboard
}

func GetUserSettingsKeyboardForUser(user model.User) *models.InlineKeyboardMarkup {
	keyboard := KeyboardSettingsUser

	return TranslateKeyboardForUser(user, keyboard)
}

func AppendTranslatedButtonToKeyboard(keyboard *models.InlineKeyboardMarkup, button models.InlineKeyboardButton, user model.User, preserveButtons ...int) {
	AppendButtonToKeyboard(keyboard, translateButtonForUser(user, button), preserveButtons...)
}

func AppendButtonToKeyboard(keyboard *models.InlineKeyboardMarkup, button models.InlineKeyboardButton, preserveButtons ...int) {
	if len(keyboard.InlineKeyboard) > keyboardMaxButtonCount {
		keyboard.InlineKeyboard = removeButtonBlock(
			keyboard.InlineKeyboard,
			findSmallestMissingInt(preserveButtons),
		)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{button})
}

func GetCompetitionFixturesKeyboardForUser(user model.User, compView manager.CompetitionView) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{}

	if len(compView.Standings) > 0 {
		standingsButton := models.InlineKeyboardButton{
			Text:         "ShowStandings",
			CallbackData: fmt.Sprintf("%s%d", CbdShowStandings, compView.CompId),
		}
		AppendTranslatedButtonToKeyboard(
			keyboard,
			standingsButton,
			user,
		)
	}

	for _, fixture := range compView.Fixtures {
		AppendButtonToKeyboard(
			keyboard,
			models.InlineKeyboardButton{
				Text:         generateFixtureButtonText(user, fixture),
				CallbackData: fmt.Sprintf("%s%d", CbdFixtureToggle, fixture.ID),
			},
		)
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
	AppendTranslatedButtonToKeyboard(keyboard, ButtonBack, *user)

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
	AppendTranslatedButtonToKeyboard(keyboard, ButtonBack, *user)

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
			userKeyboard.InlineKeyboard[0] = removeButtonFromBlock(userKeyboard.InlineKeyboard[0], i)
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

func removeButtonFromBlock(slice []models.InlineKeyboardButton, i int) []models.InlineKeyboardButton {
	slice[i] = slice[len(slice)-1]
	slice[len(slice)-1] = models.InlineKeyboardButton{}
	slice = slice[:len(slice)-1]

	return slice
}

func removeButtonBlock(slice [][]models.InlineKeyboardButton, i int) [][]models.InlineKeyboardButton {
	copy(slice[i:], slice[i+1:])
	slice[len(slice)-1] = nil
	slice = slice[:len(slice)-1]

	return slice
}

func generateFixtureButtonText(user model.User, fixture manager.FixtureView) string {
	var homeIcon, awayIcon string
	favIcon := "⭐"
	mainIcon, score := generateIconAndScore(fixture, user)
	if fixture.IsHomeUserFav {
		homeIcon = favIcon + " "
	}
	if fixture.IsAwayUserFav {
		awayIcon = favIcon + " "
	}

	buttonText := fmt.Sprintf(
		"%s %s %s%s %s %s%s",
		mainIcon,
		fixture.Date.Format(TimeFormat),
		homeIcon,
		fixture.HomeTeamName,
		score,
		awayIcon,
		fixture.AwayTeamName,
	)

	if len(buttonText) > keyboardButtonTextLength && fixture.AwayTeamCode != "" {
		buttonText = fmt.Sprintf(
			"%s %s %s%s %s %s%s",
			mainIcon,
			fixture.Date.Format(TimeFormat),
			homeIcon,
			fixture.HomeTeamName,
			score,
			awayIcon,
			fixture.AwayTeamCode,
		)
	}
	if len(buttonText) > keyboardButtonTextLength && fixture.HomeTeamCode != "" {
		buttonText = fmt.Sprintf(
			"%s %s %s%s %s %s%s",
			mainIcon,
			fixture.Date.Format(TimeFormat),
			homeIcon,
			fixture.HomeTeamCode,
			score,
			awayIcon,
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

func findSmallestMissingInt(nums []int) int {
	for i := 0; ; i++ {
		found := false
		for _, num := range nums {
			if num == i {
				found = true
				break
			}
		}
		if !found {
			return i
		}
	}
}
