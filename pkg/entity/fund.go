package entity

type AddressFund struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	Chain   string `gorm:"size:31;uniqueIndex:idx_chain_address_index"`
	Address string `gorm:"size:255;uniqueIndex:idx_chain_address_index"`

	TronUsdt float64 `gorm:""`
	TronUsdc float64 `gorm:""`
	BscUsdt  float64 `gorm:""`
	BscUsdc  float64 `gorm:""`
	EthUsdt  float64 `gorm:""`
	EthUsdc  float64 `gorm:""`

	TimeAts
}

type AddressFundCollectLog struct {
	Id             int64  `gorm:"primaryKey;autoIncrement"`
	UserId         string `gorm:"size:63"`
	AddressGroupId int64  `gorm:""`

	Chain     string  `gorm:"size:31"`
	Currency  string  `gorm:"size:31"`
	Status    string  `gorm:"size:31"`
	ToAddress string  `gorm:"size:255"`
	AmountMin float64 `gorm:""`
	FeeMax    float64 `gorm:""`

	TronUsdt float64 `gorm:""`
	TronUsdc float64 `gorm:""`
	BscUsdt  float64 `gorm:""`
	BscUsdc  float64 `gorm:""`
	EthUsdt  float64 `gorm:""`
	EthUsdc  float64 `gorm:""`

	Description string `gorm:"size:1023"`

	TimeAts
}

type AddressFundCollectDetail struct {
	Id    int64 `gorm:"primaryKey;autoIncrement"`
	LogId int64 `gorm:""`

	Chain       string `gorm:"size:31"`
	Currency    string `gorm:"size:31"`
	FromAddress string `gorm:"size:255"`
	ToAddress   string `gorm:"size:255"`

	Amount float64 `gorm:""`
	Status string  `gorm:"size:31"`

	TimeAts
}
