package repository

import (
	"context"
	"strings"

	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
	"gorm.io/gorm"
)

type OrderRepository interface {
	BaseRepository[entity.ProductOrder]
	FindAllOrders(context.Context, *valueobject.Query, uint, entity.RoleId) (*valueobject.PagedResult, error)
	FindOrderDetail(context.Context, uint, uint, entity.RoleId) (*entity.ProductOrder, error)
}

type orderRepository struct {
	*baseRepository[entity.ProductOrder]
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db:             db,
		baseRepository: &baseRepository[entity.ProductOrder]{db: db},
	}
}

func (r *orderRepository) FindAllOrders(ctx context.Context, query *valueobject.Query, userId uint, roleId entity.RoleId) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		orderStatus := query.GetConditionValue("order_status")
		name := query.GetConditionValue("name")
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "price":
			query.WithSortBy("\"total_payment\"")
		case "order_date":
			query.WithSortBy("\"ordered_at\"")
		}
		db.Joins("LEFT JOIN order_items ON product_orders.id = order_items.order_id").
			Joins("LEFT JOIN shop_products ON shop_products.id = order_items.shop_product_id").
			Joins("LEFT JOIN products ON shop_products.product_id = products.id").
			Joins("LEFT JOIN order_statuses ON order_statuses.id = product_orders.order_status_id")

		if orderStatus != nil {
			db.Where("order_status_id = ?", orderStatus)
		}
		switch roleId {
		case entity.RoleUser:
			db.Where("profile_id = ?", userId)
			if name != nil {
				db.Where("products.name ILIKE ? ", name)
			}
		case entity.RoleAdmin:
			db.Joins("LEFT JOIN shops ON shops.id = shop_products.shop_id").
				Where("shops.admin_id = ?", userId)
		}
		if roleId != entity.RoleUser {
			db.Joins("LEFT JOIN profiles ON profiles.user_id = product_orders.profile_id")
			if name != nil {
				db.Where("profiles.name ILIKE ? ", name)
			}
			db.Preload("Profile")
		}
		db.Group("product_orders.id")
		db.Preload("OrderItems.ShopProduct.Product")
		db.Preload("OrderStatus")
		return db
	})
}

func (r *orderRepository) FindOrderDetail(ctx context.Context, orderId, userId uint, roleId entity.RoleId) (*entity.ProductOrder, error) {
	var order entity.ProductOrder
	query := r.conn(ctx).
		Model(&entity.ProductOrder{}).
		Joins("LEFT JOIN order_items ON order_items.order_id = product_orders.id").
		Joins("LEFT JOIN shop_products ON shop_products.id = order_items.shop_product_id").
		Joins("LEFT JOIN shops ON shops.id = shop_products.shop_id").
		Joins("LEFT JOIN admin_contacts ON admin_contacts.user_id = shops.admin_id").
		Where("order_id = ? ", orderId).
		Preload("OrderItems.ShopProduct.Shop").
		Preload("OrderStatus")
	switch roleId {
	case entity.RoleUser:
		query = query.Where("profile_id = ?", userId)
	case entity.RoleAdmin:
		query = query.Where("shops.admin_id = ?", userId)
	case entity.RoleSuperAdmin:
		query = query.Preload("OrderItems.ShopProduct.Shop.Admin.AdminContact.User")
	}
	err := query.Find(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil

}
