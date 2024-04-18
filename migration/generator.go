package migration

import (
	"time"

	"github.com/ziadrahmatullah/ordent-test/entity"
)

func generateProfile(id uint, name string) *entity.Profile {
	return &entity.Profile{
		UserId:    id,
		Name:      name,
		Birthdate: time.Date(2000, 01, 01, 0, 0, 0, 0, time.Local),
	}
}
