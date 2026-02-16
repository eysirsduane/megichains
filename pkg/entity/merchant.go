package entity

type Merchant struct {
	Id              int64  `gorm:"primaryKey;autoIncrement"`
	MerchantAccount string `gorm:"size:63;uniqueIndex"`

	Name      string `gorm:"size:31"`
	SecretKey string `gorm:"size:63"`

	Description string `gorm:"size:2047"`

	TimeAts
}
