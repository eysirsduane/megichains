package entity

type User struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	DisplayId int32  `gorm:"uniqueIndex"`
	Nickname  string `gorm:"size:63;uniqueIndex"`
	Avatar    string `gorm:"size:127"`

	Username string `gorm:"size:63;uniqueIndex"`
	Password string `gorm:"size:255"`

	Contacts `gorm:"embedded"`
	TimeAts  `gorm:"embedded"`
}
