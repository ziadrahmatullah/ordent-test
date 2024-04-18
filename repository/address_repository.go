package repository

import (
	"github.com/ziadrahmatullah/ordent-test/entity"
	"gorm.io/gorm"
)

type AddressRepository interface {
	BaseRepository[entity.Address]
}

type addressRepository struct {
	*baseRepository[entity.Address]
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Address]{db: db},
	}
}
