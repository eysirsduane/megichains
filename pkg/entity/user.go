package entity

type User struct {
	Id string `gorm:"type:varchar(64);primaryKey;"`

	DisplayId int32  `gorm:"uniqueIndex"`
	Nickname  string `gorm:"size:63;uniqueIndex"`
	Avatar    string `gorm:"size:127"`

	Username string `gorm:"size:63;uniqueIndex"`
	Password string `gorm:"size:255"`

	Status string `gorm:"size:15"`

	Contacts `gorm:"embedded"`
	TimeAts  `gorm:"embedded"`
}
