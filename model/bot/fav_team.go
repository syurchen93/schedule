package bot

import (
	"schedule/model"
)

type FavTeam struct {
	ID     int        `gorm:"primaryKey"`
	UserID int        `gorm:"type:int;uniqueIndex:user_team_idx"`
	User   User       `gorm:"foreignKey:UserID"`
	TeamID int        `gorm:"type:int;uniqueIndex:user_team_idx"`
	Team   model.Team `gorm:"foreignKey:TeamID"`
}
