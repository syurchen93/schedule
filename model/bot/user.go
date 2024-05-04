package bot

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type User struct {
	ID                   int            `gorm:"primaryKey"`
	FirstName            string         `gorm:"type:varchar(100)"`
	LastName             string         `gorm:"type:varchar(100)"`
	Username             string         `gorm:"type:varchar(100)"`
	Locale               string         `gorm:"type:varchar(3)"`
	DisabledCountries    datatypes.JSON `gorm:"type:json"`
	DisabledCompetitions datatypes.JSON `gorm:"type:json"`
}

func (u *User) GetDisabledCountries() []int {
	var result []int
	err := json.Unmarshal(u.DisabledCountries, &result)
	if err != nil {
		return result
	}
	return result
}

func (u *User) SetDisabledCountries(ids []int) {
	jsonData, _ := json.Marshal(ids)
	u.DisabledCountries = datatypes.JSON(jsonData)
}

func (u *User) GetDisabledCompetitions() []int {
	var result []int
	err := json.Unmarshal(u.DisabledCompetitions, &result)
	if err != nil {
		return result
	}
	return result
}

func (u *User) SetDisabledCompetitons(ids []int) {
	jsonData, _ := json.Marshal(ids)
	u.DisabledCompetitions = datatypes.JSON(jsonData)
}
