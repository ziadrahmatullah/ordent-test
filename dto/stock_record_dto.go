package dto

import (
	"time"

	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type StockRecordParams struct {
	Name        *string `form:"name"`
	IsReduction *bool   `form:"is_reduction"`
	SortBy      *string `form:"sort_by" binding:"omitempty,oneof=quantity name"`
	Order       *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit       *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page        *int    `form:"page" binding:"omitempty,numeric,min=1"`
	ProductId   *uint   `form:"product_id" binding:"omitempty,numeric,min=1"`
}

func (qp *StockRecordParams) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()
	if qp.Name != nil {
		query.Condition("name", valueobject.ILike, *qp.Name)
	}
	if qp.IsReduction != nil {
		query.Condition("is_reduction", valueobject.Equal, qp.IsReduction)
	}
	if qp.Page != nil {
		query.WithPage(*qp.Page)
	}
	if qp.ProductId != nil {
		query.Condition("ShopProductId", valueobject.Equal, qp.ProductId)
	}
	if qp.Limit != nil {
		query.WithLimit(*qp.Limit)
	}
	if qp.Order != nil {
		query.WithOrder(valueobject.Order(*qp.Order))
	}
	if qp.SortBy != nil {
		query.WithSortBy(*qp.SortBy)
	} else {
		query.WithSortBy("id")
	}

	return query, nil
}

type StockMonthlyReportParams struct {
	Name   *string `form:"name"`
	SortBy *string `form:"sort_by" binding:"omitempty,oneof= additions deductions final_stock shop_name product_name"`
	Order  *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit  *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page   *int    `form:"page" binding:"omitempty,numeric,min=1"`
	Month  *uint   `form:"month" binding:"omitempty,numeric,min=1,max=12"`
}

func (qp *StockMonthlyReportParams) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()
	if qp.Name != nil {
		query.Condition("product_name", valueobject.ILike, *qp.Name)
	}
	if qp.Page != nil {
		query.WithPage(*qp.Page)
	}
	if qp.Month != nil {
		query.Condition("month", valueobject.Equal, *qp.Month)
	}
	if qp.Limit != nil {
		query.WithLimit(*qp.Limit)
	}
	if qp.Order != nil {
		query.WithOrder(valueobject.Order(*qp.Order))
	}
	if qp.SortBy != nil {
		query.WithSortBy(*qp.SortBy)
	} else {
		query.WithSortBy("id")
	}

	return query, nil
}

type ReportRes struct {
	Additions   uint   `json:"additions"`
	Deductions  uint   `json:"deductions"`
	FinalStock  int    `json:"final_stock"`
	Month       string `json:"month"`
	ShopName    string `json:"shop_name"`
	ProductName string `json:"product_name"`
}

type StockRecordRes struct {
	Id          uint             `json:"id"`
	Quantity    int              `json:"quantity"`
	IsReduction bool             `json:"is_reduction"`
	ChangeAt    time.Time        `json:"change_at"`
	Product     *ProductStockRes `json:"product"`
	ShopName    string           `json:"shop_name"`
}

type ProductStockRes struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	CategoryName string `json:"category_name"`
	Image        string `json:"image"`
}

type StockRecordReq struct {
	ShopProductId uint  `json:"shop_product_id" binding:"required,min=1"`
	Quantity      int   `json:"quantity" binding:"required,min=1"`
	IsReduction   *bool `json:"is_reduction" binding:"required"`
}

func NewStockRecordRes(p *entity.StockRecord) *StockRecordRes {
	var product *ProductStockRes
	if p.ShopProduct != nil {
		product = NewStockProductRes(p.ShopProduct)
	}
	return &StockRecordRes{Id: p.Id, Quantity: p.Quantity, IsReduction: p.IsReduction, ChangeAt: p.ChangeAt, Product: product, ShopName: p.ShopProduct.Shop.Name}
}

func NewStockProductRes(p *entity.ShopProduct) *ProductStockRes {
	return &ProductStockRes{Id: p.Id, Name: p.Product.Name, CategoryName: p.Product.ProductCategory.Name, Image: p.Product.Image}
}

func (r *StockRecordReq) ToModel() *entity.StockRecord {
	return &entity.StockRecord{ShopProductId: r.ShopProductId, Quantity: r.Quantity, IsReduction: *r.IsReduction}
}
