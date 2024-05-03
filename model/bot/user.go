package bot

type User struct {
	ID                   int    `gorm:"primaryKey"`
	FirstName            string `gorm:"type:varchar(100)"`
	LastName             string `gorm:"type:varchar(100)"`
	Username             string `gorm:"type:varchar(100)"`
	Locale               string `gorm:"type:varchar(3)"`
	DisabledCountries    []int  `gorm:"type:json"`
	DisabledCompetitions []int  `gorm:"type:json"`
}
