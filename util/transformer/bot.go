package transformer

import (
	"github.com/go-telegram/bot/models"
	model "schedule/model/bot"
)

func CreateUserFromChatMember(chatMember *models.ChatMember) model.User {
	chatUser := chatMember.Member.User
	return model.User{
		ID:        int(chatUser.ID),
		Username:  chatUser.Username,
		FirstName: chatUser.FirstName,
		LastName:  chatUser.LastName,
		Locale:    chatUser.LanguageCode,
	}
}
