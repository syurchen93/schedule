package manager

import (
	"fmt"
	"schedule/model/bot"
	"schedule/util"

	"github.com/go-telegram/bot/models"
)

const (
	BotMsgPrefix  = "bot_message_"
	CompStandings = "comp_standings_"
	UserPrefix    = "user_"

	UserTextInputModePrefix = "user_input_mode_"
)

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

func GetCachedCompetitionStandings(compId uint) []StandingsData {
	var standings []StandingsData
	cacheKey := fmt.Sprintf("%s%d", CompStandings, compId)

	err := util.GetCacheItem(cacheKey, &standings)
	if nil == err && len(standings) > 0 {
		return standings
	}

	standings = buildStandingDatas(fetchUpToDateCompetitionStandings(compId))
	util.SetCacheItem(cacheKey, standings)

	return standings
}

func SetUserTextInputMode(userId int, mode string) {
	util.SetCacheItem(fmt.Sprintf("%s%d", UserTextInputModePrefix, userId), mode)
}

func GetUserTextInputMode(userId int) *string {
	mode, ok := util.GetCacheString(fmt.Sprintf("%s%d", UserTextInputModePrefix, userId))
	if !ok {
		return nil
	}

	return &mode
}

func ClearUserTextInputMode(userId int) {
	util.DeleteCacheItem(fmt.Sprintf("%s%d", UserTextInputModePrefix, userId))
}

func GetUserFromCache(userId int) *bot.User {
	var user bot.User
	err := util.GetCacheItem(fmt.Sprintf("%s%d", UserPrefix, userId), &user)
	if nil != err {
		return nil
	}

	return &user
}

func CacheUser(user *bot.User) {
	util.SetCacheItem(fmt.Sprintf("%s%d", UserPrefix, user.ID), user)
}

func ClearUserFromCache(userId int) {
	util.DeleteCacheItem(fmt.Sprintf("%s%d", UserPrefix, userId))
}
