package dto

import (
	"errors"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type ListProductQueryParam struct {
	IsUser   *bool   `form:"is_user"`
	Name     *string `form:"name"`
	Category *int    `form:"category" binding:"omitempty,numeric,min=1"`
	SortBy   *string `form:"sort_by" binding:"omitempty,oneof=name price"`
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit    *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page     *int    `form:"page" binding:"omitempty,numeric,min=1"`
	IsHidden *bool   `form:"is_hidden" binding:"omitempty"`
}

func (qp *ListProductQueryParam) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()

	if qp.Page != nil {
		query.WithPage(*qp.Page)
	}
	if qp.Limit != nil {
		query.WithLimit(*qp.Limit)
	}

	if qp.Order != nil {
		query.WithOrder(valueobject.Order(*qp.Order))
	}

	if qp.SortBy != nil {
		query.WithSortBy(*qp.SortBy)
	}

	if qp.Name != nil {
		query.Condition("name", valueobject.ILike, *qp.Name)
	}

	if qp.Category != nil {
		query.Condition("category", valueobject.Equal, *qp.Category)
	}
	if qp.IsHidden != nil {
		query.Condition("is_hidden", valueobject.Equal, *qp.IsHidden)
	}

	return query, nil
}

type AddProductRequest struct {
	Name              string  `form:"name" binding:"required"`
	Manufacture       string  `form:"manufacture" binding:"required"`
	Detail            string  `form:"detail" binding:"required"`
	ProductCategoryId uint    `form:"product_category_id" binding:"required"`
	UnitInPack        string  `form:"unit_in_pack" binding:"required"`
	Weight            string  `form:"weight" binding:"required,numeric"`
	Height            string  `form:"height" binding:"required,numeric"`
	Length            string  `form:"length" binding:"required,numeric"`
	Width             string  `form:"width" binding:"required,numeric"`
	SellingUnit       string  `form:"selling_unit" binding:"required"`
	Price             string  `form:"price" binding:"required,numeric"`
}

func (r *AddProductRequest) Validate() error {
	weight, _ := decimal.NewFromString(r.Weight)
	if weight.LessThanOrEqual(decimal.Zero) {
		return apperror.NewInvalidPathQueryParamError(errors.New("weight should be greater than zero"))
	}

	height, _ := decimal.NewFromString(r.Height)
	if height.LessThanOrEqual(decimal.Zero) {
		return apperror.NewInvalidPathQueryParamError(errors.New("height should be greater than zero"))
	}

	length, _ := decimal.NewFromString(r.Length)
	if length.LessThan(decimal.Zero) {
		return apperror.NewInvalidPathQueryParamError(errors.New("length should be greater than zero"))
	}

	width, _ := decimal.NewFromString(r.Width)
	if width.LessThan(decimal.Zero) {
		return apperror.NewInvalidPathQueryParamError(errors.New("width should be greater than zero"))
	}
	price, _ := decimal.NewFromString(r.Price)
	if price.LessThan(decimal.Zero) {
		return apperror.NewInvalidPathQueryParamError(errors.New("price should be greater than zero"))
	}
	return nil
}

func (r *AddProductRequest) ToProduct() *entity.Product {
	weight, _ := decimal.NewFromString(r.Weight)
	height, _ := decimal.NewFromString(r.Height)
	length, _ := decimal.NewFromString(r.Length)
	width, _ := decimal.NewFromString(r.Width)
	Price, _ := decimal.NewFromString(r.Price)
	return &entity.Product{
		Name:              strings.Trim(r.Name, " "),
		Manufacture:       strings.Trim(r.Manufacture, " "),
		Detail:            strings.Trim(r.Detail, " "),
		ProductCategoryId: r.ProductCategoryId,
		UnitInPack:        r.UnitInPack,
		Weight:            weight,
		Height:            height,
		Length:            length,
		Width:             width,
		Price:             Price,
		SellingUnit:       r.SellingUnit,
	}
}

type SimpleProductResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	TopPrice    string `json:"top_price"`
	FloorPrice  string `json:"floor_price"`
	SellingUnit string `json:"selling_unit"`
	Image       string `json:"image"`
}

type ProductPriceRangeResponse struct {
	Id                uint            `json:"id"`
	Name              string          `json:"name"`
	Manufacture       string          `json:"manufacture"`
	ProductCategoryId uint            `json:"product_category_id"`
	Detail            string          `json:"detail"`
	UnitInPack        string          `json:"unit_in_pack"`
	Weight            decimal.Decimal `json:"weight"`
	Height            decimal.Decimal `json:"height"`
	Length            decimal.Decimal `json:"length"`
	Width             decimal.Decimal `json:"width"`
	Image             string          `json:"image"`
	IsHidden          bool            `json:"is_hidden"`
	Price             string          `json:"price"`
	SellingUnit       string          `json:"selling_unit"`
}

type ProductResponse struct {
	Id                uint            `json:"id"`
	Name              string          `json:"name"`
	Manufacture       string          `json:"manufacture"`
	ProductCategoryId uint            `json:"product_category_id"`
	Detail            string          `json:"detail"`
	UnitInPack        string          `json:"unit_in_pack"`
	Weight            decimal.Decimal `json:"weight"`
	Height            decimal.Decimal `json:"height"`
	Length            decimal.Decimal `json:"length"`
	Width             decimal.Decimal `json:"width"`
	Image             string          `json:"image"`
	IsHidden          bool            `json:"is_hidden"`
	Price             decimal.Decimal `json:"price"`
	SellingUnit       string          `json:"selling_unit"`
}

func NewFromShopProduct(product *entity.Product) *ProductPriceRangeResponse {
	return &ProductPriceRangeResponse{
		Id:                product.Id,
		Name:              product.Name,
		ProductCategoryId: product.ProductCategoryId,
		Price:             product.Price.String(),
		SellingUnit:       product.SellingUnit,
		Manufacture:       product.Manufacture,
		Detail:            product.Detail,
		UnitInPack:        product.UnitInPack,
		Weight:            product.Weight,
		Height:            product.Height,
		Length:            product.Length,
		Width:             product.Weight,
		Image:             product.Image,
		IsHidden:          product.IsHidden,
	}
}

func NewFromProduct(product *entity.Product) *ProductResponse {
	return &ProductResponse{
		Id:                product.Id,
		Name:              product.Name,
		ProductCategoryId: product.ProductCategoryId,
		Price:             product.Price,
		SellingUnit:       product.SellingUnit,
		Manufacture:       product.Manufacture,
		Detail:            product.Detail,
		UnitInPack:        product.UnitInPack,
		Weight:            product.Weight,
		Height:            product.Height,
		Length:            product.Length,
		Width:             product.Weight,
		Image:             product.Image,
		IsHidden:          product.IsHidden,
	}
}
