package repository

import (
	"context"

	"github.com/ziadrahmatullah/ordent-test/entity"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	BaseRepository[entity.OrderItem]
	BulkCreate(context.Context, []*entity.OrderItem) error
	ListOfOrderItem(ctx context.Context, orderId uint, userId uint) ([]*entity.OrderItem, error)
}

type orderItemRepository struct {
	*baseRepository[entity.OrderItem]
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &orderItemRepository{
		db:             db,
		baseRepository: &baseRepository[entity.OrderItem]{db: db},
	}
}

func (r *orderItemRepository) BulkCreate(ctx context.Context, items []*entity.OrderItem) error {
	return r.conn(ctx).Model(&entity.OrderItem{}).Create(items).Error
}

func (r *orderItemRepository) ListOfOrderItem(ctx context.Context, orderId uint, userId uint) ([]*entity.OrderItem, error) {
	var orderItems []*entity.OrderItem
	err := r.conn(ctx).
		Model(&entity.OrderItem{}).
		Joins("JOIN shop_products ON shop_products.id = order_items.shop_product_id ").
		Joins("JOIN shops ON shops.id = shop_products.shop_id").
		Joins("JOIN product_orders ON product_orders.id = order_items.order_id").
		Where("order_id = ?", orderId).
		Where("shops.admin_id = ?", userId).
		Preload("ShopProduct").
		Preload("ShopProduct.Shop").
		Find(&orderItems).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}
