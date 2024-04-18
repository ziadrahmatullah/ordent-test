package repository

import (
	"github.com/ziadrahmatullah/ordent-test/entity"
	"gorm.io/gorm"
)

type ForgotPasswordRepository interface {
	BaseRepository[entity.ForgotPasswordToken]
}

type forgotPasswordRepository struct {
	*baseRepository[entity.ForgotPasswordToken]
	db *gorm.DB
}

func NewForgotPasswordRepository(db *gorm.DB) ForgotPasswordRepository {
	return &forgotPasswordRepository{
		db:             db,
		baseRepository: &baseRepository[entity.ForgotPasswordToken]{db: db},
	}
}
