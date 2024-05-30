package manager

import (
	"context"
	"schedule/util/transformer"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"

	model "schedule/model/bot"
)

var dbGorm *gorm.DB
var defaultLocale string
var supportedLocales []string

var currentUser *model.User

func Init(db *gorm.DB, defaultLocaleArg string, supportedLocalesArg []string) {
	dbGorm = db
	defaultLocale = defaultLocaleArg
	supportedLocales = supportedLocalesArg
}

func GetCurrentUser() *model.User {
	return currentUser
}

func UpdateCurrentUserLocale(locale string) error {
	result := dbGorm.Model(currentUser).Update("locale", locale)
	if result.Error != nil {
		return result.Error
	}
	currentUser.Locale = locale

	return nil
}

func UpdateUserTimezone(user *model.User, timezone string) error {
	result := dbGorm.Model(user).Update("timezone", timezone)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateUserAlertOffset(user *model.User, offset int) error {
	result := dbGorm.Model(user).Update("alert_offset", offset*60)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetUserLocale(update *models.Update) (string, error) {
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

func GetOrCreateUser(ctx context.Context, b *bot.Bot, update *models.Update) *model.User {
	if currentUser != nil {
		return currentUser
	}

	user, err := getUserByUpdate(update)
	if err != nil {
		user, err = createUser(ctx, b, update)
		if err != nil {
			panic(err)
		}
	}

	return user
}

func getUserByUpdate(update *models.Update) (*model.User, error) {
	var userId int
	user := model.User{}

	if update.CallbackQuery == nil {
		userId = int(update.Message.From.ID)
	} else {
		userId = int(update.CallbackQuery.From.ID)
	}

	result := dbGorm.First(&user, userId)

	if result.Error != nil {
		return &user, result.Error
	}
	currentUser = &user

	return &user, nil
}

func createUser(ctx context.Context, b *bot.Bot, update *models.Update) (*model.User, error) {
	chatMember, err := b.GetChatMember(ctx, &bot.GetChatMemberParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		UserID: update.CallbackQuery.From.ID,
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
