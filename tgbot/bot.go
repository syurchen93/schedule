package main

import (
	"context"
	"os"
	"os/signal"

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
		bot.WithCallbackQueryDataHandler("settings", bot.MatchTypeExact, settingsHandler),
	}

	b, err := bot.New(util.GetEnv("TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {
		panic(err)
	}

	b.Start(ctx)
}

func settingsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func setLocaleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if nil != err {
		panic(err)
	}

	locale := update.CallbackQuery.Data[len("set_lang_"):]
	err = manager.UpdateCurrentUserLocale(locale)
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
	user, err := manager.CreateUser(ctx, b, update)
	if nil != err {
		panic(err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
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
