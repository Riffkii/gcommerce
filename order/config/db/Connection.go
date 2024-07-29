package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenConnection() (*gorm.DB, func()) {
	dialect := mysql.Open("root:123@tcp(127.0.0.1:3306)/gcommerce?charset=utf8mb4&parseTime=True")
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	return db, func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}
