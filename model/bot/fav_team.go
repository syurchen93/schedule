package bot

import (
	"schedule/model"
)

type FavTeam struct {
	ID     int        `gorm:"primaryKey"`
	UserID int        `gorm:"type:int"`
	User   User       `gorm:"foreignKey:UserID"`
	TeamID int        `gorm:"type:int"`
	Team   model.Team `gorm:"foreignKey:TeamID"`
}
