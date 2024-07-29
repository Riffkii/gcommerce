package model

import "time"

type Customer struct {
	ID        int64 `gorm:"primary_key;column:id;autoIncrement"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
