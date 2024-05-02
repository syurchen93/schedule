package main

import (
	"context"
	"os"
	"os/signal"

	"gorm.io/gorm"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"schedule/db"
	"schedule/util"
	"schedule/util/transformer"

	model "schedule/model/bot"
)

var dbGorm *gorm.DB
var defaultLocale = "en"
var supportedLocales = []string{"en"}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db.Init()
	dbGorm = db.Db()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, callbackHandler),
	}

	b, err := bot.New(util.GetEnv("TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {
		panic(err)
	}

	b.Start(ctx)
}

func callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if nil != err {
		panic(err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You selected the button: " + update.CallbackQuery.Data,
	})
	if nil != err {
		panic(err)
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := createUser(ctx, b, update)
	if nil != err {
		panic(err)
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Button 1", CallbackData: "button_1"},
				{Text: "Button 2", CallbackData: "button_2"},
			}, {
				{Text: "Button 3", CallbackData: "button_3"},
			},
		},
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        transateForUpdateUser("Greetings", update),
		ReplyMarkup: kb,
	})
	if nil != err {
		panic(err)
	}
}

func transateForUpdateUser(key string, update *models.Update) string {
	util.InitLocale("tgbot/translation", supportedLocales)
	locale, err := getUserLocale(update)
	if nil != err {
		panic(err)
	}
	localizer := util.GetLocalizer(locale)

	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: key,
	})
}

func getUserLocale(update *models.Update) (string, error) {
	user := model.User{}
	result := dbGorm.First(&user, update.Message.From.ID)
	if result.Error != nil {
		return "", result.Error
	}

	for _, locale := range supportedLocales {
		if user.Locale == locale {
			return user.Locale, nil
		}
	}

	return defaultLocale, nil
}

func createUser(ctx context.Context, b *bot.Bot, update *models.Update) error {
	chatMember, err := b.GetChatMember(ctx, &bot.GetChatMemberParams{
		ChatID: update.Message.Chat.ID,
		UserID: update.Message.From.ID,
	})
	if nil != err {
		return err
	}

	user := transformer.CreateUserFromChatMember(chatMember)
	result := dbGorm.FirstOrCreate(&user, model.User{ID: user.ID})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
