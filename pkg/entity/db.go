package entity

import (
	"fmt"
	"megichains/pkg/global"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg *global.Config) (db *gorm.DB, err error) {
	// envs := getEnvs([]string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_PORT", "DB_NAME", "DB_TIMEZONE"})
	dsn := fmt.Sprintf("host=%v user=%v password=%v port=%v dbname=%v sslmode=disable TimeZone=%v", cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Port, cfg.DB.Name, cfg.DB.Timezone)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	
	err = db.AutoMigrate(&RangeConfig{}, &User{}, Address{}, &BscTransaction{})
	if err != nil {
		panic(err)
	}

	return
}
