package manager

import (
	"fmt"
	"schedule/util"

	"github.com/go-telegram/bot/models"
)

const BotMsgPrefix = "bot_message_"

func CacheBotMessage(msg *models.Message) {
	util.SetCacheItem(fmt.Sprintf("%s%d", BotMsgPrefix, msg.ID), msg)
}

func GetCachedBotMessage(msgId int) *models.Message {
	var msg models.Message
	err := util.GetCacheItem(fmt.Sprintf("%s%d", BotMsgPrefix, msgId), &msg)
	if nil != err {
		return nil
	}

	return &msg
}
