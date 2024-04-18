package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
	"gorm.io/gorm"
)

type ShopRepository interface {
	BaseRepository[entity.Shop]
	FindAllShop(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	FindNearestShopFromAddress(ctx context.Context, addressId uint) ([]*entity.Shop, error)
	FindAllShopSuperAdmin(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
}

type shopRepository struct {
	*baseRepository[entity.Shop]
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) ShopRepository {
	return &shopRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Shop]{db: db},
	}
}

func (r *shopRepository) FindAllShop(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "name":
			query.WithSortBy("\"shops\".name")
		case "id":
			query.WithSortBy("\"shops\".id ")
		}
		db.Where("\"shops\".admin_id = ?", ctx.Value("user_id").(uint))
		name := query.GetConditionValue("name")
		db.Joins("City").Joins("Province")

		if name != nil {
			db.Where("\"shops\".name ILIKE ?", name)
		}
		return db
	})
}

func (r *shopRepository) FindAllShopSuperAdmin(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "shop_name":
			query.WithSortBy("\"shops\".name")
		case "admin_name":
			query.WithSortBy("\"Admin__AdminContact\".name")

		case "id":
			query.WithSortBy("\"shops\".id ")
		}
		name := query.GetConditionValue("name")
		db.Joins("City").Joins("Province").Joins("Admin.AdminContact")
		province := query.GetConditionValue("province")
		if name != nil {
			db.Where("\"shops\".name ILIKE ?", name)
		}
		if province != nil {
			db.Where("\"shops\".province_id = ?", province)
		}
		return db
	})
}

func (r *shopRepository) FindNearestShopFromAddress(ctx context.Context, addressId uint) ([]*entity.Shop, error) {
	var shop []*entity.Shop
	err := r.db.
		Joins("City").
		Joins("JOIN addresses a ON st_dwithin(\"shops\".location, a.location, 25000)").
		Where("a.id=?", addressId).
		Order("\"shops\".location <-> a.location").
		Find(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return shop, nil
}
