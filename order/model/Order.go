package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         int64 `gorm:"primary_key;column:id;autoIncrement"`
	OrderCode  string
	CustomerId int64
	ProductId  int64
	Quantity   int32
	FinalPrice int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
	Customer   Customer `gorm:"foreignKey:customer_id;references:id"`
}
