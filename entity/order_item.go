package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderItem struct {
	Id            uint            `gorm:"primaryKey;autoIncrement"`
	OrderId       uint            `gorm:"not null"`
	Order         ProductOrder    `gorm:"foreignKey:OrderId;references:Id"`
	ShopProductId uint            `gorm:"not null"`
	ShopProduct   ShopProduct     `gorm:"foreignKey:ShopProductId;references:Id"`
	Quantity      int             `gorm:"not null"`
	SubTotal      decimal.Decimal `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
