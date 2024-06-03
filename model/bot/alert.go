package bot

import "schedule/model/league"

type Alert struct {
	ID               uint           `gorm:"primaryKey"`
	UserID           uint           `gorm:"uniqueIndex:user_fixture_time_before_idx;not null"`
	User             User           `gorm:"foreignKey:UserID"`
	FixtureID        uint           `gorm:"uniqueIndex:user_fixture_time_before_idx;not null"`
	Fixture          league.Fixture `gorm:"foreignKey:FixtureID"`
	TimeBefore       int            `gorm:"uniqueIndex:user_fixture_time_before_idx;not null"`
	IsFired          bool           `gorm:"default:false"`
	IsFavTeamCreated bool           `gorm:"default:false"`
}
