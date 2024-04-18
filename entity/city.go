package entity

import (
	"time"

	"github.com/ziadrahmatullah/ordent-test/valueobject"
	"gorm.io/gorm"
)

type City struct {
	Id         uint                    `gorm:"primaryKey;autoIncrement"`
	Name       string                  `gorm:"not null"`
	Code       string                  `gorm:"not null"`
	CityGid    uint                    `json:"city_gid"`
	Location   *valueobject.Coordinate `gorm:"not null;type:geography(POINT,4326)"`
	ProvinceId uint                    `gorm:"not null"`
	Province   Province
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
