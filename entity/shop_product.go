package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ShopProduct struct {
	Id        uint `gorm:"primaryKey;autoIncrement"`
	ProductId uint `gorm:"not null"`
	Product   *Product
	ShopId    uint `gorm:"not null"`
	Shop      *Shop
	Stock     int             `gorm:"not null"`
	Price     decimal.Decimal `gorm:"not null;type:numeric"`
	IsActive  bool            `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
