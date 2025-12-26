package entity

type EvmLog struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	ChainId uint64 `gorm:"uniqueIndex:idx_evm_trans_index"`
	TxHash  string `gorm:"size:255;uniqueIndex:idx_evm_trans_index"`
	Index   uint   `gorm:"uniqueIndex:idx_evm_trans_index"`
	TxIndex uint   `gorm:""`

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
