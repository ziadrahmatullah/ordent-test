package repository

import (
	"context"
	"strings"

	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
	"gorm.io/gorm"
)

type ShopProductRepository interface {
	BaseRepository[entity.ShopProduct]
	FindAllShopProducts(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	BulkCreate(ctx context.Context, products []*entity.ShopProduct) ([]*entity.ShopProduct, error)
}

type shopProductRepository struct {
	*baseRepository[entity.ShopProduct]
	db *gorm.DB
}

func NewShopProductRepository(db *gorm.DB) ShopProductRepository {
	return &shopProductRepository{
		db:             db,
		baseRepository: &baseRepository[entity.ShopProduct]{db: db},
	}
}

func (r *shopProductRepository) FindAllShopProducts(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "name":
			query.WithSortBy("\"Product\".name")
		case "id":
			query.WithSortBy("\"shop_products\".id ")
		}

		category := query.GetConditionValue("category")
		name := query.GetConditionValue("name")
		shop := query.GetConditionValue("shop")
		isActive := query.GetConditionValue("is_active")
		db.Joins("Product").Joins("Product.ProductCategory").Joins("Shop")
		if category != nil {
			db.Where("\"Product\".product_category_id", category)
		}

		if name != nil {
			db.Where("\"Product\".name ILIKE ?", name)
		}
		if isActive != nil {
			db.Where("is_active", isActive)
		}
		db.Where("Shop.id", shop)
		return db
	})
}

func (r *shopProductRepository) BulkCreate(ctx context.Context, products []*entity.ShopProduct) ([]*entity.ShopProduct, error) {
	err := r.conn(ctx).Model(&entity.ShopProduct{}).Create(products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
