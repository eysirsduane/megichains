package entity

type EvmOrder struct {
	Id           int64  `gorm:"primaryKey;autoIncrement"`
	MerchOrderId string `gorm:"size:63"`

	Chain   string `gorm:"size:15"`
	ChainId uint64 `gorm:"uniqueIndex:idx_trans_index"`
	TxHash  string `gorm:"size:255;uniqueIndex:idx_trans_index"`
	Index   uint   `gorm:"uniqueIndex:idx_trans_index"`
	TxIndex uint   `gorm:""`

	Typo     string `gorm:"size:15"`
	Status   string `gorm:"size:15"`
	Currency string `gorm:"size:15"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    int64   `gorm:""`

	FromHex string `gorm:"size:255"`
	ToHex   string `gorm:"size:255"`

	Contract       string `gorm:"size:255"`
	BlockHash      string `gorm:"size:255"`
	BlockNumber    uint64 `gorm:""`
	BlockTimestamp uint64 `gorm:""`

	Removed     bool   `gorm:""`
	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}
