package model

import "time"

type Team struct {
	ID                 int       `gorm:"primaryKey"`
	Name               string    `gorm:"type:varchar(100); not null"`
	Logo               string    `gorm:"type:varchar(255); not null"`
	Country            *string   `gorm:"type:varchar(100)"`
	Code               *string   `gorm:"type:varchar(100)"`
	TriedToFetchCodeAt time.Time `gorm:"default:null"`
	IsUserFav          bool
}
