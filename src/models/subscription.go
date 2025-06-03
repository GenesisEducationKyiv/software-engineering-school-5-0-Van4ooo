package models

import "time"

type Subscription struct {
	ID        uint   `gorm:"PrimaryKey"`
	Email     string `gorm:"unique;not null"`
	City      string `gorm:"not null"`
	Frequency string `gorm:"not null"`
	Token     string `gorm:"unique; not null"`
	Confirmed bool   `gorm:"default:false"`
	CreatedAt time.Time
}
