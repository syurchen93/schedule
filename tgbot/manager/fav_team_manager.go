package manager

import (
	"schedule/model"
	"schedule/model/bot"
	"schedule/model/league"
)

func GetFavTeamsForUser(userId int) []bot.FavTeam {
	var favTeams []bot.FavTeam
	dbGorm.Where("user_id = ?", userId).
		Preload("Team").
		Find(&favTeams)
	return favTeams
}

func RemoveFavTeamForUser(userId int, teamId int) {
	dbGorm.Unscoped().Where("user_id = ? AND team_id = ?", userId, teamId).Delete(&bot.FavTeam{})
	RemoveAlertsForUserFavTeam(&bot.User{ID: userId}, teamId)
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

func AddFavTeamForUser(user *bot.User, teamId int) {
	dbGorm.Create(&bot.FavTeam{
		UserID: user.ID,
		TeamID: teamId,
	})
	CreateAlertsForUserFavTeamFixtures(user)
}

func CreateAlertsForUserFavTeamFixtures(user *bot.User) {
	var favTeams []bot.FavTeam
	dbGorm.Where("user_id = ?", user.ID).Find(&favTeams)

	for _, favTeam := range favTeams {
		var fixtures []league.Fixture
		dbGorm.
			Where("(home_team_id = ? OR away_team_id = ?) AND date > NOW()", favTeam.TeamID, favTeam.TeamID).
			Find(&fixtures)

		for _, fixture := range fixtures {
			alert := bot.Alert{
				UserID:           uint(user.ID),
				FixtureID:        uint(fixture.ID),
				TimeBefore:       user.AlertOffset,
				IsFavTeamCreated: true,
			}
			dbGorm.Create(&alert)
		}
	}
}

func RemoveAlertsForUserFavTeam(user *bot.User, teamId int) {
	var alerts []bot.Alert
	dbGorm.
		Where("user_id = ? AND is_fav_team_created = 1", user.ID).
		Joins("join fixture on alert.fixture_id = fixture.id").
		Where("home_team_id = ? OR away_team_id = ?", teamId, teamId).
		Find(&alerts)

	for _, alert := range alerts {
		dbGorm.Delete(&alert)
	}
}
