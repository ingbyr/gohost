package dal

import (
	"fmt"

	"github.com/ingbyr/gohost/dal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(dsn string) (db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("fail to connect db: %w", err))
	}
	return db
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&model.Host{},
		&model.Group{},
	)
	if err != nil {
		panic(fmt.Errorf("fail to migrate: %w", err))
	}
}
