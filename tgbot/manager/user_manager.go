package manager

import (
	"context"
	"gorm.io/gorm"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/bot"
	"schedule/util/transformer"

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

func CreateUser(ctx context.Context, b *bot.Bot, update *models.Update) (*model.User, error) {
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

func UpdateCurrentUserLocale(locale string) error {
	result := dbGorm.Model(currentUser).Update("locale", locale)
	if result.Error != nil {
		return result.Error
	}
	currentUser.Locale = locale

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

func getUserByUpdate(update *models.Update) (model.User, error) {
	user := model.User{}

	result := dbGorm.First(&user, update.Message.From.ID)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
