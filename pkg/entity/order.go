package entity

type EvmOrder struct {
	Id           int64  `gorm:"primaryKey;autoIncrement"`
	LogId        int64  `gorm:""`
	MerchOrderId string `gorm:"size:63;uniqueIndex"`
	TxHash       string `gorm:"size255"`
	Chain        string `gorm:"size:15;"`

	Typo     string `gorm:"size:15"`
	Status   string `gorm:"size:15"`
	Currency string `gorm:"size:15"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    int64   `gorm:""`

	FromHex string `gorm:"size:255"`
	ToHex   string `gorm:"size:255"`

	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}

type TronOrder struct {
	Id            int64  `gorm:"primaryKey;autoIncrement"`
	LogId         int64  `gorm:""`
	MerchOrderId  string `gorm:"size:63;uniqueIndex"`
	TransactionId string `gorm:"size:255;"`
	Chain         string `gorm:"size:15;"`

	Typo     string `gorm:"size:15"`
	Status   string `gorm:"size:15"`
	Currency string `gorm:"size:15"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    int64   `gorm:""`

	FromBase58 string `gorm:"size:255"`
	ToBase58   string `gorm:"size:255"`

	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}

type SolanaOrder struct {
	Id           int64  `gorm:"primaryKey;autoIncrement"`
	MerchOrderId string `gorm:"size:63"`

	Chain  string `gorm:"size:15"`
	TxHash string `gorm:"size:255;uniqueIndex"`

	Typo     string `gorm:"size:15"`
	Status   string `gorm:"size:15"`
	Currency string `gorm:"size:15"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    int64   `gorm:""`

	FromBase58 string `gorm:"size:255"`
	ToBase58   string `gorm:"size:255"`

	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}
