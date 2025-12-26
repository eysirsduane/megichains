package entity

type EvmLog struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	Chain    string `gorm:""`
	Currency string `gorm:"size:15"`
	ChainId  uint64 `gorm:"uniqueIndex:idx_evm_trans_index"`
	TxHash   string `gorm:"size:255;uniqueIndex:idx_evm_trans_index"`
	Index    uint   `gorm:"uniqueIndex:idx_evm_trans_index"`
	TxIndex  uint   `gorm:""`

	Amount float64 `gorm:""`
	Sun    int64   `gorm:""`

	FromHex string `gorm:"size:255"`
	ToHex   string `gorm:"size:255"`

	Contract       string `gorm:"size:255"`
	BlockHash      string `gorm:"size:255"`
	BlockNumber    uint64 `gorm:""`
	BlockTimestamp uint64 `gorm:""`

	Removed bool `gorm:""`

	TimeAts `gorm:"embedded"`
}

type TronTransaction struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	Chain         string `gorm:""`
	Currency      string `gorm:"size:15"`
	TransactionId string `gorm:"size:255;uniqueIndex"`

	Amount float64 `gorm:""`
	Sun    int64   `gorm:""`

	FromBase58 string `gorm:"size:255"`
	ToBase58   string `gorm:"size:255"`

	Contract       string `gorm:"size:255"`
	BlockTimestamp uint64 `gorm:""`

	TimeAts `gorm:"embedded"`
}
