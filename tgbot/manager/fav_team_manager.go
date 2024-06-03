package manager

import (
	"schedule/model"
	"schedule/model/bot"
)

func GetFavTeamsForUser(userId int) []bot.FavTeam {
	var favTeams []bot.FavTeam
	dbGorm.Where("user_id = ?", userId).
		Preload("Team").
		Find(&favTeams)
	return favTeams
}

func RemoveFavTeamForUser(userId int, teamId int) {
	var favTeam bot.FavTeam
	dbGorm.Where("user_id = ? AND team_id = ?", userId, teamId).First(&favTeam)
	dbGorm.Delete(&favTeam)
}

func FindTeamByUserInput(userInput string) *model.Team {
	var favTeam model.Team
	dbGorm.
		Where("name LIKE ?", "%"+userInput+"%").
		First(&favTeam)

	if favTeam.ID == 0 {
		return nil
	}
	return &favTeam
}

func AddFavTeamForUser(userId int, teamId int) {
	dbGorm.Create(&bot.FavTeam{
		UserID: userId,
		TeamID: teamId,
	})
}
