package entity

type Contacts struct {
	Telegram string `gorm:"size:255"`
	WhatsApp string `gorm:"column:whatsapp;size:255"`
	WeChat   string `gorm:"column:wechat;size:255"`
	Other    string `gorm:"size:255"`
}

type TimeAts struct {
	UpdatedAt uint64 `gorm:"autoUpdateTime:milli"` // 使用时间戳纳秒数填充更新时间
	DeletedAt uint64 `gorm:"autoDeleteTime:milli"` // 使用时间戳毫秒数填充更新时间
	CreatedAt uint64 `gorm:"autoCreateTime:milli"` // 使用时间戳秒数填充创建时间
}
