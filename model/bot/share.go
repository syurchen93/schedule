package bot

type UserShare struct {
	ID           int  `gorm:"primaryKey"`
	SourceUserId uint `gorm:"not null"`
	SourceUser   User `gorm:"foreignKey:SourceUserId"`
	TargetUserId int  `gorm:"not null"`
	TargetUser   User `gorm:"foreignKey:TargetUserId"` // Subscriber
}
