package entity

type Order struct {
	Id            int64  `gorm:"type:BIGSERIAL;primaryKey;autoIncrement"`
	UserId        int64  `gorm:""`
	TransactionId string `gorm:"size:255;uniqueIndex"`

	Typo     int32  `gorm:""`
	Status   string `gorm:"size:15"`
	Currency string `gorm:"size:15"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    uint64  `gorm:""`

	FromBase58 string `gorm:"size:255"`
	ToBase58   string `gorm:"size:255"`
	FromHex    string `gorm:"size:255"`
	ToHex      string `gorm:"size:255"`

	DelegateAmount float64 `gorm:""`
	DelegateSun    int64   `gorm:""`

	Time         uint64 `gorm:""`
	Expires      uint64 `gorm:""`
	WithdrawTime uint64
	FailedTimes  uint32

	Description string `gorm:"size:2047"`

	Contacts `gorm:"embedded"`
	TimeAts  `gorm:"embedded"`
}

type ExchangeOrder struct {
	Id            int64  `gorm:"type:bigserial;primaryKey;autoIncrement"`
	UserId        int64  `gorm:""`
	TransactionId string `gorm:"size:255;uniqueIndex"`

	Typo     string `gorm:"size:15"`
	Status   string `gorm:"size:15"`
	Currency string `gorm:"size:15"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    int64   `gorm:""`

	FromBase58 string `gorm:"size:255"`
	ToBase58   string `gorm:"size:255"`
	FromHex    string `gorm:"size:255"`
	ToHex      string `gorm:"size:255"`

	ThenRate         float64
	ExchangeRate     float64
	ExchangeDiscount float64
	ExchangeAmount   float64 `gorm:""`
	ExchangeSun      int64   `gorm:""`

	Time    uint64 `gorm:""`
	Expires uint64 `gorm:""`

	Description string `gorm:"size:1023"`

	Contacts `gorm:"embedded"`
	TimeAts  `gorm:"embedded"`
}
