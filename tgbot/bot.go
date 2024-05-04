package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"

	"gorm.io/gorm"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"schedule/db"
	"schedule/tgbot/manager"
	"schedule/tgbot/template"
	"schedule/util"
)

var dbGorm *gorm.DB
var defaultLocale = "en"
var supportedLocales = []string{"en", "ru", "de"}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db.Init()
	dbGorm = db.Db()

	util.InitTranslator("tgbot/translation", supportedLocales)
	manager.Init(dbGorm, defaultLocale, supportedLocales)

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithCallbackQueryDataHandler("set_lang_", bot.MatchTypePrefix, setLocaleHandler),
		bot.WithCallbackQueryDataHandler("settings", bot.MatchTypeExact, settingsGeneralHandler),
		bot.WithCallbackQueryDataHandler("settings_country", bot.MatchTypeExact, settingsCountryHandler),
		bot.WithCallbackQueryDataHandler("settings_country_toggle", bot.MatchTypePrefix, settingsCountryToggleHandler),
	}

	b, err := bot.New(util.GetEnv("TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {
		panic(err)
	}

	b.Start(ctx)
}

func settingsCountryToggleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)
	countryId, err := strconv.Atoi(update.CallbackQuery.Data[len("settings_country_toggle_"):])
	if nil != err {
		panic(err)
	}

	manager.ToggleUserCountrySettings(user, countryId)

	success, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: template.GetUserCountrySettingsKyboard(
			user,
			manager.GetUserCountrySettings(user),
		),
	})
	if nil != err {
		panic(err)
	}
	if success != nil {
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   transateForUpdateUser("Done", update),
	})
	if nil != err {
		panic(err)
	}
}

func settingsCountryHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	userCountrySettings := manager.GetUserCountrySettings(user)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        transateForUpdateUser("SettingsCountryHeader", update),
		ReplyMarkup: template.GetUserCountrySettingsKyboard(user, userCountrySettings),
	})
	if nil != err {
		panic(err)
	}
}

func settingsGeneralHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        transateForUpdateUser("SettingsGeneral", update),
		ReplyMarkup: template.TranslateKeyboardForUser(*user, template.KeyboardSettingsGeneral),
	})
	if nil != err {
		panic(err)
	}
}

func setLocaleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	locale := update.CallbackQuery.Data[len("set_lang_"):]
	_ = manager.GetOrCreateUser(ctx, b, update)

	err := manager.UpdateCurrentUserLocale(locale)
	if nil != err {
		panic(err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        transateForUpdateUser("Greetings", update),
		ReplyMarkup: template.GetLanguageSelectKeyboardForUser(*manager.GetCurrentUser()),
	})
	if nil != err {
		panic(err)
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        transateForUpdateUser("Greetings", update),
		ReplyMarkup: template.GetLanguageSelectKeyboardForUser(*user),
	})
	if nil != err {
		panic(err)
	}
}

func transateForUpdateUser(key string, update *models.Update) string {
	locale, err := manager.GetUserLocale(update)
	if nil != err {
		panic(err)
	}

	return util.Translate(locale, key)
}

func answerCallbackQuery(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if nil != err {
		panic(err)
	}
}
