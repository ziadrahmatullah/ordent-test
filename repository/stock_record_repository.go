package repository

import (
	"context"
	"strings"

	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
	"gorm.io/gorm"
)

type StockRecordRepository interface {
	BaseRepository[entity.StockRecord]
	FindAllStockRecord(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	BulkCreate(ctx context.Context, records []*entity.StockRecord) error
}

type stockRecordRepository struct {
	*baseRepository[entity.StockRecord]
	db *gorm.DB
}

func NewStockRecordRepository(db *gorm.DB) StockRecordRepository {
	return &stockRecordRepository{
		db:             db,
		baseRepository: &baseRepository[entity.StockRecord]{db: db},
	}
}

func (r *stockRecordRepository) FindAllStockRecord(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "id":
			query.WithSortBy("\"stock_records\".id ")
		case "name":
			query.WithSortBy("\"ShopProduct__Product\".name")
		}
		db.Joins("ShopProduct.Product.ProductCategory").Joins("ShopProduct.Shop")

		isReduction := query.GetConditionValue("is_reduction")
		name := query.GetConditionValue("name")
		ProductId := query.GetConditionValue("ShopProductId")
		db.Where("\"ShopProduct__Shop\".admin_id =?", ctx.Value("user_id").(uint))
		if isReduction != nil {
			db.Where("\"stock_records\".is_reduction = ?", isReduction)
		}
		if ProductId != nil {
			db.Where("\"stock_records\".shop_product_id", ProductId)
		}
		if name != nil {
			db.Where("\"ShopProduct__Product\".name ILIKE ?", name)
		}
		return db
	})
}

func (r *stockRecordRepository) BulkCreate(ctx context.Context, records []*entity.StockRecord) error {
	return r.conn(ctx).Model(&entity.StockRecord{}).Create(records).Error
}
