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

func Init(db *gorm.DB, defaultLocaleArg string, supportedLocalesArg []string) {
	dbGorm = db
	defaultLocale = defaultLocaleArg
	supportedLocales = supportedLocalesArg
}

func UpdateUserLocale(user *model.User, locale string) error {
	result := dbGorm.Model(user).Update("locale", locale)
	if result.Error != nil {
		return result.Error
	}
	ClearUserFromCache(user.ID)

	return nil
}

func UpdateUserTimezone(user *model.User, timezone string) error {
	result := dbGorm.Model(user).Update("timezone", timezone)
	if result.Error != nil {
		return result.Error
	}
	ClearUserFromCache(user.ID)

	return nil
}

func UpdateUserAlertOffset(user *model.User, offset int) error {
	result := dbGorm.Model(user).Update("alert_offset", offset*60)
	if result.Error != nil {
		return result.Error
	}
	ClearUserFromCache(user.ID)

	return nil
}

func GetUserLocale(update *models.Update) (string, error) {
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
	user := &model.User{}

	if update.CallbackQuery == nil {
		userId = int(update.Message.From.ID)
	} else {
		userId = int(update.CallbackQuery.From.ID)
	}

	user = GetUserFromCache(userId)
	if user != nil {
		return user, nil
	}

	result := dbGorm.First(&user, userId)

	if result.Error != nil {
		return user, result.Error
	}
	CacheUser(user)

	return user, nil
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

	return &user, nil
}
