package dto

import (
	"errors"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type ShopReq struct {
	Name                    string   `json:"name" binding:"required"`
	Address                 string   `json:"address" binding:"required"`
	CityId                  uint     `json:"city_id" binding:"required"`
	ProvinceId              uint     `json:"province_id" binding:"required"`
	Latitude                string   `json:"latitude" binding:"required,latitude"`
	Longitude               string   `json:"longitude" binding:"required,longitude"`
	StartTime               string   `json:"start_time"   binding:"required"`
	EndTime                 string   `json:"end_time"  binding:"required"`
	OperationalDay          []string `json:"operational_day" binding:"required"`
	PharmacistLicenseNumber string   `json:"pharmacist_license_number" binding:"required"`
	PharmacistPhoneNumber   string   `json:"pharmacist_phone_number" binding:"required,phonenumberprefix,phonenumberlength"`
}

type RequestShopUri struct {
	Id uint `uri:"shop_id" binding:"required,numeric"`
}

func (r *ShopReq) ToModel() (*entity.Shop, error) {
	startTime, err := time.Parse("15:04:05", r.StartTime)
	if err != nil {
		return nil, apperror.NewClientError(err)
	}
	startTime = time.Date(1, 1, 1, startTime.Hour(), startTime.Minute(), 0, 0, time.Local)

	endTime, err := time.Parse("15:04:05", r.EndTime)
	if err != nil {
		return nil, apperror.NewClientError(err)
	}
	endTime = time.Date(1, 1, 1, endTime.Hour(), endTime.Minute(), 0, 0, time.Local)
	if !startTime.Before(endTime) {
		return nil, apperror.NewClientError(errors.New("start time must be less than end time"))
	}
	operationalDayString := strings.Join(r.OperationalDay, ",")

	latitude, _ := decimal.NewFromString(r.Latitude)
	longitude, _ := decimal.NewFromString(r.Longitude)

	return &entity.Shop{
		Name:       r.Name,
		Address:    r.Address,
		CityId:     r.CityId,
		ProvinceId: r.ProvinceId,
		Location: &valueobject.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		},
		StartTime:               startTime,
		EndTime:                 endTime,
		OperationalDay:          operationalDayString,
		PharmacistLicenseNumber: r.PharmacistLicenseNumber,
		PharmacistPhoneNumber:   r.PharmacistPhoneNumber,
	}, nil
}

type ShopRes struct {
	Id                      uint              `json:"id"`
	Name                    string            `json:"name"`
	AdminId                 uint              `json:"admin_id"`
	Address                 string            `json:"address"`
	City                    *CityRes          `json:"city"`
	Province                *ProvinceRes      `json:"province"`
	Latitude                decimal.Decimal   `json:"latitude"`
	Longitude               decimal.Decimal   `json:"longitude"`
	StartTime               string            `json:"start_time"`
	EndTime                 string            `json:"end_time" `
	OperationalDay          []string          `json:"operational_day"`
	PharmacistLicenseNumber string            `json:"pharmacist_license_number"`
	PharmacistPhoneNumber   string            `json:"pharmacist_phone_number"`
	Products                []*ProductShopRes `json:"prodcuts,omitempty"`
}

func NewShopRes(p *entity.Shop) *ShopRes {
	productsTemp := []*ProductShopRes{}
	for _, prodcut := range p.Products {
		prodcutTemp := NewProductPhamarcyRes(&prodcut)
		productsTemp = append(productsTemp, prodcutTemp)
	}
	operationalDay := strings.Split(p.OperationalDay, ",")
	return &ShopRes{
		Id:                      p.Id,
		Name:                    p.Name,
		AdminId:                 p.AdminId,
		City:                    NewCityProvinceRes(&p.City),
		Province:                NewProvinceRes(&p.Province),
		Address:                 p.Address,
		Latitude:                p.Location.Latitude,
		Longitude:               p.Location.Longitude,
		StartTime:               p.StartTime.Format("15:04"),
		EndTime:                 p.EndTime.Format("15:04"),
		OperationalDay:          operationalDay,
		PharmacistLicenseNumber: p.PharmacistLicenseNumber,
		PharmacistPhoneNumber:   p.PharmacistPhoneNumber,
		Products:                productsTemp}
}

type ShopSuperAdminRes struct {
	Id             uint          `json:"id"`
	Name           string        `json:"name"`
	Admin          *AdminShopRes `json:"admin"`
	StartTime      string        `json:"start_time"`
	EndTime        string        `json:"end_time"`
	OperationalDay []string      `json:"operational_day"`
	City           string        `json:"city"`
	Province       string        `json:"province"`
}

func NewShopSuperAdminRes(p *entity.Shop) *ShopSuperAdminRes {
	operationalDay := strings.Split(p.OperationalDay, ",")
	return &ShopSuperAdminRes{
		Id:             p.Id,
		Name:           p.Name,
		Admin:          NewAdminShopRes(&p.Admin),
		OperationalDay: operationalDay,
		StartTime:      p.StartTime.Format("15:04"),
		EndTime:        p.EndTime.Format("15:04"),
		City:           p.City.Name,
		Province:       p.Province.Name}
}

type ListShopQueryParam struct {
	Name   *string `form:"name"`
	SortBy *string `form:"sort_by" binding:"omitempty,oneof=name start_time end_time"`
	Order  *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit  *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page   *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

func (qp *ListShopQueryParam) ToQuery() (*valueobject.Query, error) {
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
	} else {
		query.WithSortBy("id")
	}

	if qp.Name != nil {
		query.Condition("name", valueobject.ILike, *qp.Name)
	}

	return query, nil
}

type ListShopSuperAdminQueryParam struct {
	Name     *string `form:"name"`
	SortBy   *string `form:"sort_by" binding:"omitempty,oneof=shop_name admin_name"`
	Province *uint   `form:"province" binding:"omitempty,numeric,min=1"`
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit    *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page     *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

func (qp *ListShopSuperAdminQueryParam) ToQuery() (*valueobject.Query, error) {
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
	} else {
		query.WithSortBy("id")
	}

	if qp.Name != nil {
		query.Condition("name", valueobject.ILike, *qp.Name)
	}
	if qp.Province != nil {
		query.Condition("province", valueobject.Equal, *qp.Province)
	}

	return query, nil
}
