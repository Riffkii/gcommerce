package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        int64  `gorm:"primary_key;column:id;autoIncrement"`
	Name      string `gorm:"column:name"`
	Stock     int32  `gorm:"column:stock"`
	Price     int64  `gorm:"column:price"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
