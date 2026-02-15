package entity

type MerchOrder struct {
	Id            int64  `gorm:"primaryKey;autoIncrement"`
	LogId         int64  `gorm:""`
	MerchOrderId  string `gorm:"size:63;uniqueIndex"`
	TransactionId string `gorm:"size:255"`
	Chain         string `gorm:"size:15;"`

	Typo         string `gorm:"size:15"`
	Status       string `gorm:"size:15"`
	NotifyStatus string `gorm:"size:15"`
	Currency     string `gorm:"size:15"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    int64   `gorm:""`

	FromAddress string `gorm:"size:255"`
	ToAddress   string `gorm:"size:255"`

	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}
