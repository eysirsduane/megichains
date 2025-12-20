package entity

type Address struct {
	Id      int64 `gorm:"type:bigserial;primaryKey;autoIncrement"`
	GroupId int64 `gorm:""`

	Chain      string `gorm:"size:31"`
	Typo       string `gorm:"size:31"`
	Status     string `gorm:"size:15"`
	AddressHex string `gorm:"size:255"`
	Secrect    string `gorm:"size:2047"`

	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}

type AddressGroup struct {
	Id int64 `gorm:"type:bigserial;primaryKey;autoIncrement"`

	Name        string `gorm:"size:31"`
	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}
