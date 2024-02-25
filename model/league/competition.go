package league

type CompetitionType string

const (
    LEAGUE CompetitionType = "league"
    CUP CompetitionType = "cup"
)


type Competition struct {
	ID uint `gorm:"primaryKey"`
	CountryID uint `gorm:"not null"`
	Country Country `gorm:"foreignKey:CountryID"`
	Name string `gorm:"type:varchar(256); not null"`
	Type string `gorm:"type:enum('league', 'cup'); not null"`
	Logo string `gorm:"type:varchar(256); not null"`
}