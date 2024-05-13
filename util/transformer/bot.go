package transformer

import (
	"github.com/go-telegram/bot/models"
	model "schedule/model/bot"
)

func CreateUserFromChatMember(chatMember *models.ChatMember, update *models.Update) model.User {
	chatUser := chatMember.Member.User
	return model.User{
		ID:        int(chatUser.ID),
		ChatId:    int(update.CallbackQuery.Message.Message.Chat.ID),
		Username:  chatUser.Username,
		FirstName: chatUser.FirstName,
		LastName:  chatUser.LastName,
		Locale:    chatUser.LanguageCode,
	}
}
