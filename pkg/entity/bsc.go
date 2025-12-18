package entity

type BscTransaction struct {
	Id      int64  `gorm:"type:BIGSERIAL;primaryKey;autoIncrement"`
	TxHash  string `gorm:"size:255;uniqueIndex"`
	TxIndex uint   `gorm:""`

	Typo     int32  `gorm:""`
	Status   string `gorm:"size:15"`
	Currency string `gorm:"size:15"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    int64   `gorm:""`

	FromHex string `gorm:"size:255"`
	ToHex   string `gorm:"size:255"`
	Index   uint   `gorm:""`

	Address        string `gorm:"size:255"`
	BlockHash      string `gorm:"size:255"`
	BlockNumber    uint64 `gorm:""`
	BlockTimestamp uint64 `gorm:""`

	Removed     bool   `gorm:""`
	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}
