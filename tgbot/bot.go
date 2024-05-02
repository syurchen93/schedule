package main

import (
	"context"
	"os"
	"os/signal"

	"gorm.io/gorm"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"schedule/db"
	"schedule/tgbot/template"
	"schedule/util"
	"schedule/util/transformer"

	model "schedule/model/bot"
)

var dbGorm *gorm.DB
var defaultLocale = "en"
var supportedLocales = []string{"en", "ru", "de"}

var currentUser *model.User

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db.Init()
	dbGorm = db.Db()

	util.InitTranslator("tgbot/translation", supportedLocales)

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithCallbackQueryDataHandler("set_lang_", bot.MatchTypePrefix, callbackHandler),
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

	locale := update.CallbackQuery.Data[len("set_lang_"):]
	err = updateCurrentUserLocale(locale)
	if nil != err {
		panic(err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        transateForUpdateUser("Greetings", update),
		ReplyMarkup: template.GetLanguageSelectKeyboardForUser(*currentUser),
	})
	if nil != err {
		panic(err)
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user, err := createUser(ctx, b, update)
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
	locale, err := getUserLocale(update)
	if nil != err {
		panic(err)
	}

	return util.Translate(locale, key)
}

func createUser(ctx context.Context, b *bot.Bot, update *models.Update) (*model.User, error) {
	if currentUser != nil {
		return currentUser, nil
	}

	chatMember, err := b.GetChatMember(ctx, &bot.GetChatMemberParams{
		ChatID: update.Message.Chat.ID,
		UserID: update.Message.From.ID,
	})
	if nil != err {
		return nil, err
	}

	user := transformer.CreateUserFromChatMember(chatMember)
	result := dbGorm.FirstOrCreate(&user, model.User{ID: user.ID})
	if result.Error != nil {
		return nil, result.Error
	}
	currentUser = &user

	return &user, nil
}

func updateCurrentUserLocale(locale string) error {
	result := dbGorm.Model(currentUser).Update("locale", locale)
	if result.Error != nil {
		return result.Error
	}
	currentUser.Locale = locale

	return nil
}

func getUserLocale(update *models.Update) (string, error) {
	if currentUser != nil {
		return currentUser.Locale, nil
	}
	user, err := getUserByUpdate(update)
	if err != nil {
		return "", err
	}

	for _, locale := range supportedLocales {
		if user.Locale == locale {
			return user.Locale, nil
		}
	}

	return defaultLocale, nil
}

func getUserByUpdate(update *models.Update) (model.User, error) {
	user := model.User{}

	result := dbGorm.First(&user, update.Message.From.ID)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
