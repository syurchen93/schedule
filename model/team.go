package model

type Team struct {
	ID      int    `gorm:"primaryKey"`
	Name    string `gorm:"type:varchar(100); not null"`
	Code    string `gorm:"type:varchar(100); not null"`
	Country string `gorm:"type:varchar(100); not null"`
	Logo    string `gorm:"type:varchar(255); not null"`
}
