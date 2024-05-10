package bot

import "schedule/model/league"

type Alert struct {
	ID         uint           `gorm:"primaryKey"`
	UserID     uint           `gorm:"not null"`
	User       User           `gorm:"foreignKey:UserID"`
	FixtureID  uint           `gorm:"not null"`
	Fixture    league.Fixture `gorm:"foreignKey:FixtureID"`
	TimeBefore int            `gorm:"not null"`
	IsFired    bool           `gorm:"default:false"`
}
