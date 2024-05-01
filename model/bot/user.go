package bot

type User struct {
	ID        int    `gorm:"primaryKey"`
	FirstName string `gorm:"type:varchar(100)"`
	LastName  string `gorm:"type:varchar(100)"`
	Username  string `gorm:"type:varchar(100)"`
	Locale    string `gorm:"type:varchar(3)"`
}
