package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	timer "atomicgo.dev/schedule"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"schedule/tgbot/manager"
	"schedule/tgbot/template"
	"schedule/util"
)

const AlertIntervalSeconds = 10

var defaultLocale = "en"
var supportedLocales = []string{"en", "ru", "de"}

const (
	UserTextInputModeCity        = "city"
	UserTextInputModeAlertOffset = "alert_offset"
	UserTextInputModeFavTeam     = "fav_team"
	UserTextInputModeAlertShare  = "alert_share"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	util.InitTranslator("tgbot/translation", supportedLocales)
	manager.Init(defaultLocale, supportedLocales)
	util.InitCache(time.Hour, 10_000)
	manager.InitImageGenerator("tgbot/images/")
	util.InitGoogleMapsClient(ctx, util.GetEnv("GOOGLE_MAPS_API_KEY"))

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMessageTextHandler("", bot.MatchTypeContains, textHandler),

		bot.WithCallbackQueryDataHandler(template.CbdSetLang, bot.MatchTypePrefix, setLocaleHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettings, bot.MatchTypeExact, settingsGeneralHandler),

		bot.WithCallbackQueryDataHandler(template.CbdSettingsCountry, bot.MatchTypeExact, settingsCountryHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsCountryToggle, bot.MatchTypePrefix, settingsCountryToggleHandler),

		bot.WithCallbackQueryDataHandler(template.CbdSettingsCompetition, bot.MatchTypeExact, settingsCompetitionHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsCompetitionToggle, bot.MatchTypePrefix, settingsCompetitionToggleHandler),

		bot.WithCallbackQueryDataHandler(template.CbdSettingsAlert, bot.MatchTypeExact, settingsAlertHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsUserAlertOffset, bot.MatchTypeExact, settingsUserAlertOffsetHandler),

		bot.WithCallbackQueryDataHandler(template.CbdSettingsUser, bot.MatchTypeExact, settingsUserHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsTimezone, bot.MatchTypeExact, settingsUserTimezoneHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsFavTeam, bot.MatchTypeExact, settingsUserFavTeamHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsFavTeamAddStart, bot.MatchTypeExact, settingsUserFavTeamAddHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsFavTeamRemove, bot.MatchTypePrefix, settingsUserFavTeamRemoveHandler),
		bot.WithCallbackQueryDataHandler(template.CbdSettingsShareManage, bot.MatchTypeExact, settingsUserShareManageHandler),

		bot.WithCallbackQueryDataHandler(template.CbdSchedule, bot.MatchTypeExact, scheduleHandler),
		bot.WithCallbackQueryDataHandler(template.CbdFixtureToggle, bot.MatchTypePrefix, fixtureToggleHandler),
		bot.WithCallbackQueryDataHandler(template.CbdShowStandings, bot.MatchTypePrefix, standingsHandler),
	}

	b, err := bot.New(util.GetEnv("TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {
		panic(err)
	}

	timer.Every(AlertIntervalSeconds*time.Second, func() bool {
		return manager.GetAndFireAlerts(ctx, b)
	})

	b.Start(ctx)
}

func settingsUserShareManageHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)
	userShares, err := manager.GetShareSubscriptionsForUser(user.ID)
	if nil != err {
		panic(err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text:                util.Translate(user.Locale, "SettingsShareSubscriptionsInfo", user.Username),
		ReplyMarkup:         template.GetShareAlertsKeyboardForUser(userShares, *user),
	})
	if nil != err {
		panic(err)
	}
	manager.SetUserTextInputMode(user.ID, UserTextInputModeAlertShare)
}

func settingsUserFavTeamRemoveHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)
	favTeamId, err := strconv.Atoi(update.CallbackQuery.Data[len(template.CbdSettingsFavTeamRemove):])
	if nil != err {
		panic(err)
	}

	manager.RemoveFavTeamForUser(user.ID, favTeamId)

	msg := manager.GetCachedBotMessage(update.CallbackQuery.Message.Message.ID)

	if msg != nil {
		_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			ReplyMarkup: template.RemoveFavTeamFromCachedKeyboard(favTeamId, msg.ReplyMarkup),
		})
		if nil != err {
			panic(err)
		}
		manager.CacheBotMessage(msg)

		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        util.Translate(user.Locale, "SettingsFavTeamRemoveSuccess"),
		ReplyMarkup: template.TranslateKeyboardForUser(*user, template.KeyboardToSchedule),
	})
	if nil != err {
		panic(err)
	}
}

func settingsUserFavTeamAddHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)
	manager.SetUserTextInputMode(user.ID, UserTextInputModeFavTeam)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text:                util.Translate(user.Locale, "SettingsFavTeamAddInput"),
		ReplyMarkup:         template.TranslateKeyboardForUser(*user, template.KeyboardBack),
	})
	if nil != err {
		panic(err)
	}
}

func settingsUserFavTeamHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)
	favTeams := manager.GetFavTeamsForUser(user.ID)

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text:                transateForUpdateUser("SettingsFavTeam", update),
		ReplyMarkup:         template.GetFavTeamKeyboardForUser(favTeams, *user),
	})
	if nil != err {
		panic(err)
	}
	manager.CacheBotMessage(msg)
}

func settingsUserAlertOffsetHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	minOffset := user.AlertOffset / 60

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text:                util.Translate(user.Locale, "SettingsUserAlertOffsetInfo", minOffset),
		ReplyMarkup:         template.TranslateKeyboardForUser(*user, template.KeyboardBack),
	})
	if nil != err {
		panic(err)
	}

	manager.SetUserTextInputMode(user.ID, UserTextInputModeAlertOffset)
}

func settingsUserTimezoneHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)
	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		panic(err)
	}
	currentTimeInUserTimezone := time.Now().In(loc)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text: util.Translate(
			user.Locale,
			"SettingsTimezoneInfo",
			user.Timezone,
			currentTimeInUserTimezone.Format(template.TimeFormat),
		),
		ReplyMarkup: template.TranslateKeyboardForUser(*user, template.KeyboardBack),
	})
	if nil != err {
		panic(err)
	}
	manager.SetUserTextInputMode(user.ID, UserTextInputModeCity)
}

func settingsUserHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text:                transateForUpdateUser("SettingsUser", update),
		ReplyMarkup:         template.GetUserSettingsKeyboardForUser(*user),
	})
	if nil != err {
		panic(err)
	}
}

func settingsAlertHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	alertCompViews := manager.GetAlertCompetitionViewsForUser(user.ID)
	for i, compView := range alertCompViews {
		keyboard := template.GetCompetitionFixturesKeyboardForUser(*user, compView)
		if i == len(alertCompViews)-1 {
			fmt.Println("Adding refresh and settings buttons")
			template.AppendTranslatedButtonToKeyboard(keyboard, template.ButtonSettings, *user, 0)
			template.AppendTranslatedButtonToKeyboard(keyboard, template.ButtonSchedule, *user, 0)
		}
		msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
			Text:                fmt.Sprintf("%s %s", manager.GetCountryEmoji(compView.CountryName), compView.CompName),
			DisableNotification: true,
			ReplyMarkup:         keyboard,
		})
		if nil != err {
			panic(err)
		}
		manager.CacheBotMessage(msg)
	}

	if len(alertCompViews) == 0 {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
			Text:                util.Translate(user.Locale, "NoAlerts"),
			DisableNotification: true,
		})
		if nil != err {
			panic(err)
		}
	}
}

func standingsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	compId, err := strconv.Atoi(update.CallbackQuery.Data[len(template.CbdShowStandings):])
	if nil != err {
		panic(err)
	}
	standings := manager.GetCachedCompetitionStandings(uint(compId))

	user := manager.GetOrCreateUser(ctx, b, update)
	standingsFilePath, err := manager.GetStandingsImage(compId, standings, user.Locale)
	if nil != err {
		panic(err)
	}
	standingsFile, err := os.Open(standingsFilePath)
	if nil != err {
		panic(err)
	}

	_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Photo: &models.InputFileUpload{
			Filename: standingsFilePath,
			Data:     standingsFile,
		},
		Caption: util.Translate(user.Locale, "Standings"),
	})

	if nil != err {
		panic(err)
	}
}

func fixtureToggleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var editedKeyboard *models.InlineKeyboardMarkup

	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)
	fixtureId, err := strconv.Atoi(update.CallbackQuery.Data[len(template.CbdFixtureToggle):])
	if nil != err {
		panic(err)
	}

	originalMsg := manager.GetCachedBotMessage(update.CallbackQuery.Message.Message.ID)
	if originalMsg == nil || originalMsg.ID == 0 {
		competitionFixtures := manager.GetCompetitionFixturesAndToggleByFixtureId(user, fixtureId)
		editedKeyboard = template.GetCompetitionFixturesKeyboardForUser(
			*user,
			competitionFixtures,
		)
	} else {
		fixtureView := manager.GetToggleFixtureViewByFixtureId(user, fixtureId)
		editedKeyboard = template.ToggleFixtureOnCachedKeyboard(*user, fixtureView, originalMsg.ReplyMarkup)
	}

	msg, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: editedKeyboard,
	})
	if nil != err {
		panic(err)
	}
	checkIfSuccessfulMessageEdit(ctx, b, update, msg)

	manager.CacheBotMessage(msg)
}

func scheduleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answerCallbackQuery(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	competitions := manager.GetCompetitionViewsForUser(user)

	for i, compView := range competitions {
		replyMarkup := template.GetCompetitionFixturesKeyboardForUser(*user, compView)
		if i == len(competitions)-1 {
			template.AppendTranslatedButtonToKeyboard(replyMarkup, template.ButtonSettings, *user)
			template.AppendTranslatedButtonToKeyboard(replyMarkup, template.ButtonRefreshSchedule, *user)
		}
		msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
			Text:                fmt.Sprintf("%s %s", manager.GetCountryEmoji(compView.CountryName), compView.CompName),
			DisableNotification: true,
			ReplyMarkup:         replyMarkup,
		})
		if nil != err {
			panic(err)
		}
		manager.CacheBotMessage(msg)
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
	manager.ClearUserTextInputMode(user.ID)

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

	locale := update.CallbackQuery.Data[len(template.CbdSetLang):]
	_ = manager.GetOrCreateUser(ctx, b, update)

	user := manager.GetOrCreateUser(ctx, b, update)

	err := manager.UpdateUserLocale(user, locale)
	if nil != err {
		panic(err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              update.CallbackQuery.Message.Message.Chat.ID,
		Text:                transateForUpdateUser("Greetings", update),
		ReplyMarkup:         template.GetLanguageSelectKeyboardForUser(*user),
	})
	if nil != err {
		panic(err)
	}
}

func textHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := manager.GetOrCreateUser(ctx, b, update)

	mode := manager.GetUserTextInputMode(user.ID)

	if mode == nil {
		defaultHandler(ctx, b, update)
		return
	}
	manager.ClearUserTextInputMode(user.ID)

	switch *mode {
	case UserTextInputModeCity:
		userTextInputModeCityHandler(ctx, b, update)
	case UserTextInputModeAlertOffset:
		userTextInputModeAlertOffsetHandler(ctx, b, update)
	case UserTextInputModeFavTeam:
		userTextInputModeFavTeamHandler(ctx, b, update)
	case UserTextInputModeAlertShare:
		userTextInputModeAlertShareHandler(ctx, b, update)
	}
}

func userTextInputModeAlertShareHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := manager.GetOrCreateUser(ctx, b, update)
	subUsername := update.Message.Text

	err := manager.SubUserByTargetUsername(user, subUsername)

	if nil != err {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        util.Translate(user.Locale, "SettingsAlertShareNotFound"),
			ReplyMarkup: template.TranslateKeyboardForUser(*user, template.KeyboardBack),
		})
		if nil != err {
			panic(err)
		}
		manager.SetUserTextInputMode(user.ID, UserTextInputModeAlertShare)
		return
	}

	// todo edit keyboard, add new sub
}

func userTextInputModeFavTeamHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := manager.GetOrCreateUser(ctx, b, update)
	favTeamInput := update.Message.Text

	team := manager.FindTeamByUserInput(favTeamInput)
	if team == nil {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        util.Translate(user.Locale, "SettingsFavTeamAddNotFound"),
			ReplyMarkup: template.TranslateKeyboardForUser(*user, template.KeyboardBack),
		})
		if nil != err {
			panic(err)
		}
		return
	}
	manager.AddFavTeamForUser(user, team.ID)
	manager.ClearUserTextInputMode(user.ID)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        util.Translate(user.Locale, "SettingsFavTeamAddSuccess", team.Name),
		ReplyMarkup: template.TranslateKeyboardForUser(*user, template.KeyboardToSchedule),
	})
	if nil != err {
		panic(err)
	}
}

func userTextInputModeAlertOffsetHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := manager.GetOrCreateUser(ctx, b, update)
	alertOffset, err := strconv.Atoi(update.Message.Text)
	if nil != err {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        util.Translate(user.Locale, "SettingsUserAlertOffsetInvalid"),
			ReplyMarkup: template.TranslateKeyboardForUser(*user, template.KeyboardBack),
		})
		if nil != err {
			panic(err)
		}
		manager.SetUserTextInputMode(user.ID, UserTextInputModeAlertOffset)
		return
	}

	err = manager.UpdateUserAlertOffset(user, alertOffset)
	if nil != err {
		panic(err)
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        util.Translate(user.Locale, "SettingsUserAlertOffsetSuccess", alertOffset),
		ReplyMarkup: template.TranslateKeyboardForUser(*user, template.KeyboardToSchedule),
	})
	if nil != err {
		panic(err)
	}
}

func userTextInputModeCityHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var timezone string
	var err error

	location := update.Message.Location
	if location != nil {
		timezone, err = util.GetTimezoneByLatLong(location.Latitude, location.Longitude)
	} else {
		timezone, err = util.GetTimezoneByCityName(update.Message.Text)
	}

	handleTimezoneResult(b, ctx, update, timezone, err)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var chatId int64
	answerCallbackQuery(ctx, b, update)

	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Message.Chat.ID
	} else {
		chatId = update.Message.Chat.ID
	}

	user := manager.GetOrCreateUser(ctx, b, update)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		DisableNotification: true,
		ChatID:              chatId,
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
	if (update.CallbackQuery == nil) || (update.CallbackQuery.ID == "") {
		return
	}
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if nil != err {
		defaultHandler(ctx, b, update)
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

func handleTimezoneResult(b *bot.Bot, ctx context.Context, update *models.Update, timezone string, err error) {
	user := manager.GetOrCreateUser(ctx, b, update)

	if err != nil {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        util.Translate(user.Locale, "SettingsTimezoneCityNotFound"),
			ReplyMarkup: template.TranslateKeyboardForUser(*user, template.KeyboardBack),
		})
		if nil != err {
			panic(err)
		}
		manager.SetUserTextInputMode(user.ID, UserTextInputModeCity)
		return
	}

	err = manager.UpdateUserTimezone(user, timezone)
	if err != nil {
		panic(err)
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        util.Translate(user.Locale, "SettingsTimezoneCitySuccess", timezone),
		ReplyMarkup: template.TranslateKeyboardForUser(*user, template.KeyboardToSchedule),
	})
	if nil != err {
		panic(err)
	}
}
