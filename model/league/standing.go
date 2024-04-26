package league

import (
	"schedule/model"
	"time"
)

type Standing struct {
	ID            uint        `gorm:"primaryKey"`
	TeamID        uint        `gorm:"uniqueIndex:idx_team_competition"`
	Team          model.Team  `gorm:"foreignKey:TeamID"`
	CompetitionID uint        `gorm:"uniqueIndex:idx_team_competition"`
	Competition   Competition `gorm:"foreignKey:CompetitionID"`
	Rank          int         `gorm:"type:int"`
	Points        int         `gorm:"type:int"`
	GoalsDiff     int         `gorm:"type:int"`
	Group         string      `gorm:"type:varchar(100)"`
	Form          string      `gorm:"type:varchar(100)"`
	Status        string      `gorm:"type:varchar(100)"`
	Description   string      `gorm:"type:varchar(255)"`
	Played        int         `gorm:"type:int"`
	Won           int         `gorm:"type:int"`
	Drawn         int         `gorm:"type:int"`
	Lost          int         `gorm:"type:int"`
	GoalsFor      int         `gorm:"type:int"`
	GoalsAgainst  int         `gorm:"type:int"`
	UpdatedApi    time.Time   `gorm:"type:timestamp"`
	UpdatedAt     time.Time   `gorm:"type:timestamp"`
}
