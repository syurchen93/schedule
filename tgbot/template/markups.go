package template

import (
	"github.com/go-telegram/bot/models"
	model "schedule/model/bot"
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

	return userKeyboard
}

func remove(slice []models.InlineKeyboardButton, i int) []models.InlineKeyboardButton {
	slice[i] = slice[len(slice)-1]
	slice[len(slice)-1] = models.InlineKeyboardButton{}
	slice = slice[:len(slice)-1]

	return slice
}
