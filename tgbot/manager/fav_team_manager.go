package manager

import (
	model "schedule/model/bot"
)

func GetFavTeamsForUser(userId int) []model.FavTeam {
	var favTeams []model.FavTeam
	dbGorm.Where("user_id = ?", userId).
		Preload("Team").
		Find(&favTeams)
	return favTeams
}
