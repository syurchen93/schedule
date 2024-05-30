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
	Timezone             string         `gorm:"type:varchar(100);default:Europe/Berlin"`
	DisabledCountries    datatypes.JSON `gorm:"type:json"`
	DisabledCompetitions datatypes.JSON `gorm:"type:json"`
	EnableSpoilers       bool           `gorm:"type:boolean;default:false"`
	AlertOffset          int            `gorm:"type:int;default:1800"`
}

func (u *User) GetDisabledCountries() []int {
	var result []int
	_ = json.Unmarshal(u.DisabledCountries, &result)
	return result
}

func (u *User) SetDisabledCountries(ids []int) {
	jsonData, _ := json.Marshal(ids)
	u.DisabledCountries = datatypes.JSON(jsonData)
}

func (u *User) GetDisabledCompetitions() []int {
	var result []int
	_ = json.Unmarshal(u.DisabledCompetitions, &result)
	return result
}

func (u *User) SetDisabledCompetitons(ids []int) {
	jsonData, _ := json.Marshal(ids)
	u.DisabledCompetitions = datatypes.JSON(jsonData)
}
