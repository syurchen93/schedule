package template

import (
	"github.com/go-telegram/bot/models"
	model "schedule/model/bot"
	"schedule/util"
	"strings"
)

var LanguageSelectKeyboard = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "🇬🇧 English", CallbackData: "set_lang_en"},
			{Text: "🇷🇺 Русский", CallbackData: "set_lang_ru"},
			{Text: "🇩🇪 Deutsch", CallbackData: "set_lang_de"},
		},
	},
}

var KeyboardSettingsGeneral = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "SettingsCountry", CallbackData: "settings_country"},
			{Text: "SettingsCompetition", CallbackData: "settings_competition"},
			{Text: "SettingsAlert", CallbackData: "settings_alert"},
		},
	},
}

var ButtonSettings = models.InlineKeyboardButton{
	Text:         "ToSettings",
	CallbackData: "settings",
}

var ButtonSchedule = models.InlineKeyboardButton{
	Text:         "ToSchedule",
	CallbackData: "schedule",
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
