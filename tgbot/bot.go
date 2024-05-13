package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"gorm.io/gorm"

	timer "atomicgo.dev/schedule"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"schedule/db"
	"schedule/tgbot/manager"
	"schedule/tgbot/template"
	"schedule/util"
)

const AlertInterval = 60

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
		bot.WithCallbackQueryDataHandler(template.CbdSetLang, bot.MatchTypePrefix, setLocaleHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettings, bot.MatchTypeExact, settingsGeneralHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsCountry, bot.MatchTypeExact, settingsCountryHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsCountryToggle, bot.MatchTypePrefix, settingsCountryToggleHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsCompetition, bot.MatchTypeExact, settingsCompetitionHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsCompetitionToggle, bot.MatchTypePrefix, settingsCompetitionToggleHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSchedule, bot.MatchTypeExact, scheduleHandler),
		bot.WithCallbackQueryDataHandler(template.CbdFixtureToggle, bot.MatchTypePrefix, fixtureToggleHandler),
	}

	b, err := bot.New(util.GetEnv("TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {
		panic(err)
	}

	timer.Every(AlertInterval*time.Second, func() bool {
		return manager.GetAndFireAlerts(ctx, b)
	})

	b.Start(ctx)
}

func fixtureToggleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)
	fixtureId, err := strconv.Atoi(update.CallbackQuery.Data[len(template.CbdFixtureToggle):])
	if nil != err {
		panic(err)
	}

	competitionFixtures := manager.GetCompetitionFixturesAndToggleByFixtureId(user, fixtureId)
	success, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: template.GetFixturesKeyboardForUser(
			*user,
			competitionFixtures.Fixtures,
		),
	})
	if nil != err {
		panic(err)
	}
	checkIfSuccessfulMessageEdit(ctx, b, update, success)
}

func scheduleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	competitionFixtures := manager.GetCompetitionFixturesForUser(user)

	for _, comp := range competitionFixtures {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
			Text:                fmt.Sprintf("%s %s", manager.GetCountryEmoji(comp.CountryName), comp.CompName),
			DisableNotification: true,
			ReplyMarkup:         template.GetFixturesKeyboardForUser(*user, comp.Fixtures),
		})
		if nil != err {
			panic(err)
		}
	}
}

func settingsCompetitionToggleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)
	compId, err := strconv.Atoi(update.CallbackQuery.Data[len(template.CbdSettingsCompetitionToggle):])
	if nil != err {
		panic(err)
	}

	manager.ToggleUserCompetitionSettings(user, compId)

	success, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: template.GetUserCompetitonSettingsKyboard(
			user,
			manager.GetUserCountryCompetitionSettings(user, manager.GetCompetitionCountryID(compId)),
		),
	})
	if nil != err {
		panic(err)
	}
	checkIfSuccessfulMessageEdit(ctx, b, update, success)
}

func settingsCompetitionHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	userCountries := manager.GetUserEnabledCountries(user)

	baseMessageKey := "SettingsCompetitionHeader"

	if len(userCountries) == 0 {
		baseMessageKey = "SettingsCompetitionHeaderNoCountries"
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		DisableNotification: true,
		Text:                util.Translate(user.Locale, baseMessageKey),
	})
	if nil != err {
		panic(err)
	}

	for _, country := range userCountries {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			DisableNotification: true,
			ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
			Text:                manager.GetCountryWithEmoji(country.Name),
			ReplyMarkup: template.GetUserCompetitonSettingsKyboard(
				user,
				manager.GetUserCountryCompetitionSettings(user, country.ID),
			),
		})
		if nil != err {
			panic(err)
		}
	}
}

func settingsCountryToggleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)
	countryId, err := strconv.Atoi(update.CallbackQuery.Data[len(template.CbdSettingsCountryToggle):])
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
	checkIfSuccessfulMessageEdit(ctx, b, update, success)
}

func settingsCountryHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	userCountrySettings := manager.GetUserCountrySettings(user)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text:                transateForUpdateUser("SettingsCountryHeader", update),
		ReplyMarkup:         template.GetUserCountrySettingsKyboard(user, userCountrySettings),
	})
	if nil != err {
		panic(err)
	}
}

func settingsGeneralHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text:                transateForUpdateUser("SettingsGeneral", update),
		ReplyMarkup:         template.TranslateKeyboardForUser(*user, template.KeyboardSettingsGeneral),
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
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text:                transateForUpdateUser("Greetings", update),
		ReplyMarkup:         template.GetLanguageSelectKeyboardForUser(*manager.GetCurrentUser()),
	})
	if nil != err {
		panic(err)
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text:                transateForUpdateUser("Greetings", update),
		ReplyMarkup:         template.GetLanguageSelectKeyboardForUser(*user),
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

func checkIfSuccessfulMessageEdit(ctx context.Context, b *bot.Bot, update *models.Update, success *models.Message) {
	if success != nil {
		return
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   transateForUpdateUser("Done", update),
	})
	if nil != err {
		panic(err)
	}
}
