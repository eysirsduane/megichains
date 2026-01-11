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
