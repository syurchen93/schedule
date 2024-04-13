package league

type Country struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(256); not null"`
	Code string `gorm:"type:varchar(2); unique_index; not null"`
	Flag string `gorm:"type:varchar(256); not null"`
}
