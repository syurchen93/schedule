package league

import (
	"github.com/syurchen93/api-football-client/common"
	"schedule/model"
	"time"
)

type Fixture struct {
	ID            int                  `gorm:"primaryKey"`
	CompetitionID uint                 `gorm:"not null"`
	Competition   Competition          `gorm:"foreignKey:CompetitionID"`
	HomeTeamID    uint                 `gorm:"not null"`
	HomeTeam      model.Team           `gorm:"foreignKey:HomeTeamID"`
	AwayTeamID    uint                 `gorm:"not null"`
	AwayTeam      model.Team           `gorm:"foreignKey:AwayTeamID"`
	Status        common.FixtureStatus `gorm:"type:varchar(4)"`
	GoalsHome     int                  `gorm:"type:int"`
	GoalsAway     int                  `gorm:"type:int"`
	PenaltyHome   int                  `gorm:"type:int"`
	PenaltyAway   int                  `gorm:"type:int"`
	Date          time.Time            `gorm:"type:timestamp"`
	UpdatedAt     time.Time            `gorm:"type:timestamp"`
	HasUserAlert  bool                 `gorm:"column:has_user_alert"`
}
