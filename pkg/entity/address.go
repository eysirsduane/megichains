package entity

type Address struct {
	Id      int64 `gorm:"primaryKey;autoIncrement"`
	GroupId int64 `gorm:""`

	Chain      string `gorm:"size:31"`
	Typo       string `gorm:"size:31"`
	Status     string `gorm:"size:15"`
	Address    string `gorm:"size:255"`
	Address2   string `gorm:"size:255"`
	PrivateKey string `gorm:"size:2047"`
	PublicKey  string `gorm:"size:2047"`

	Description string `gorm:"size:2047"`

	TimeAts
}

type AddressGroup struct {
	Id int64 `gorm:"type:bigserial;primaryKey;autoIncrement"`

	Name        string `gorm:"size:31"`
	Status      string `gorm:"size:15"`
	Description string `gorm:"size:2047"`

	TimeAts
}
