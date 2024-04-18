package entity

import "github.com/shopspring/decimal"

type CartItem struct {
	Id            uint            `gorm:"primaryKey"`
	ShopProductId uint            `gorm:"not null"`
	ShopProduct   ShopProduct     `gorm:"foreignKey:ShopProductId;references:Id"`
	CartId        uint            `gorm:"not null"`
	Cart          Cart            `gorm:"foreignKey:CartId;references:UserId"`
	Quantity      int             `gorm:"not null"`
	SubAmount     decimal.Decimal `gorm:"not null;type:numeric"`
	IsChecked     bool            `gorm:"not null"`
}
