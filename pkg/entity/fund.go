package entity

type AddressFundCollect struct {
	Id             int64  `gorm:"primaryKey;autoIncrement"`
	UserId         string `gorm:"size:63"`
	AddressGroupId int64  `gorm:""`

	Chain           string  `gorm:"size:31"`
	Currency        string  `gorm:"size:31"`
	Status          string  `gorm:"size:31"`
	ReceiverAddress string  `gorm:"size:255"`
	AmountMin       float64 `gorm:""`
	FeeMax          float64 `gorm:""`

	SuccessAmount float64 `gorm:""`
	TotalCount    int64   `gorm:""`
	SuccessCount  int64   `gorm:""`

	Description string `gorm:"size:1023"`

	TimeAts
}

type AddressFundCollectLog struct {
	Id        int64 `gorm:"primaryKey;autoIncrement"`
	CollectId int64 `gorm:""`

	Chain           string `gorm:"size:31"`
	Currency        string `gorm:"size:31"`
	FromAddress     string `gorm:"size:255"`
	ReceiverAddress string `gorm:"size:255"`

	Amount        float64 `gorm:""`
	Status        string  `gorm:"size:31"`
	TransactionId string  `gorm:"size:255"`

	Description string `gorm:"size:1023"`

	TimeAts
}
