package entity

type DelegateWithdrawal struct {
	Id            int64  `gorm:"type:bigserial;primaryKey;autoIncrement"`
	UserId        int64  `gorm:""`
	OrderId       int64  `gorm:""`
	TransactionId string `gorm:"size:255;uniqueIndex"`

	Typo   int32  `gorm:""`
	Status string `gorm:""`

	FromBase58 string `gorm:"size:255"`
	ToBase58   string `gorm:"size:255"`
	FromHex    string `gorm:"size:255"`
	ToHex      string `gorm:"size:255"`

	UnDelegatedAmount float64 `gorm:""`
	UnDelegatedSun    int64   `gorm:""`

	FailedTimes int32  `gorm:""`
	Time        uint64 `gorm:""`
	Description string `gorm:"size:1023"`

	TimeAts `gorm:"embedded"`
}
