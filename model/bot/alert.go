package bot

import "schedule/model/league"

type Alert struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	User       User
	FixtureID  uint
	Fixture    league.Fixture
	TimeBefore int
	IsFired    bool
}
