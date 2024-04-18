package entity

import (
	"time"

	"gorm.io/gorm"
)

type StockRecord struct {
	Id            uint `gorm:"primaryKey;autoIncrement"`
	ShopProductId uint `gorm:"not null"`
	ShopProduct   *ShopProduct
	Quantity      int       `gorm:"not null"`
	IsReduction   bool      `gorm:"not null"`
	ChangeAt      time.Time `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
