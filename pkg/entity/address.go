package entity

type Address struct {
	Id int64 `gorm:"type:bigserial;primaryKey;autoIncrement"`

	Chain      string `gorm:"size:31"`
	Typo       string `gorm:"size:31"`
	Status     string `gorm:"size:15"`
	AddressHex string `gorm:"size:255"`

	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}
