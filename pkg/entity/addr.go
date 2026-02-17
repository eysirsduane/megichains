package entity

type Address struct {
	Id      int64 `gorm:"primaryKey;autoIncrement"`
	GroupId int64 `gorm:""`

	Chain             string `gorm:"size:31"`
	Typo              string `gorm:"size:31"`
	Status            string `gorm:"size:15"`
	Address           string `gorm:"size:255;uniqueIndex"`
	Address2          string `gorm:"size:255"`
	PrivateKey        string `gorm:"size:2047"`
	PublicKey         string `gorm:"size:2047"`
	LastUsedTimestamp uint64 `gorm:""`

	Description string `gorm:"size:2047"`

	AddressBalance AddressBalance `gorm:"foreignKey:Address;references:Address"`

	TimeAts
}

type AddressGroup struct {
	Id int64 `gorm:"type:bigserial;primaryKey;autoIncrement"`

	Name        string `gorm:"size:31"`
	Status      string `gorm:"size:15"`
	Description string `gorm:"size:2047"`

	TimeAts
}

type AddressBalance struct {
	Id      int64  `gorm:"primaryKey;autoIncrement"`
	Address string `gorm:"size:255;uniqueIndex"`

	TronUsdt   float64 `gorm:""`
	TronUsdc   float64 `gorm:""`
	BscUsdt    float64 `gorm:""`
	BscUsdc    float64 `gorm:""`
	EthUsdt    float64 `gorm:""`
	EthUsdc    float64 `gorm:""`
	SolanaUsdt float64 `gorm:""`
	SolanaUsdc float64 `gorm:""`

	TimeAts
}
