package entity

type ExchangeConfig struct {
	Id int64 `gorm:"type:bigserial;primaryKey;autoIncrement"`

	Typo string `gorm:"size:32"`

	RangeFrom float64 `gorm:""`
	RangeTo   float64 `gorm:""`

	Value float64 `gorm:""`

	TimeAts `gorm:"embedded"`
}
