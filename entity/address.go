package entity

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	Id         uint   `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"not null"`
	StreetName string `gorm:"not null"`
	PostalCode string `gorm:"not null"`
	Phone      string `gorm:"not null"`
	Detail     string
	IsDefault  bool    `gorm:"not null"`
	ProfileId  uint    `gorm:"not null"`
	Profile    Profile `gorm:"foreignKey:ProfileId;references:UserId"`
	Province   string
	City       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
