package migration

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/ordent-test/entity"
)

func generateProfile(id uint, name string) *entity.Profile {
	return &entity.Profile{
		UserId:    id,
		Name:      name,
		Birthdate: time.Date(2000, 01, 01, 0, 0, 0, 0, time.Local),
	}
}

func generateDecimalFromString(decimalString string) decimal.Decimal {
	zero := decimal.NewFromInt(0)
	decimal, err := decimal.NewFromString(decimalString)
	if err != nil {
		return zero
	}
	return decimal
}
