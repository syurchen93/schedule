package manager

import (
	"fmt"
	"schedule/util"

	"github.com/go-telegram/bot/models"
)

const BotMsgPrefix = "bot_message_"
const CompStandings = "comp_standings_"

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
