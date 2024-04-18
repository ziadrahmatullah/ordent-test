package dto

import (
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type AdminShopQueryReq struct {
	Email  *string `form:"email"`
	SortBy *string `form:"sort_by" binding:"omitempty,oneof=email"`
	Order  *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit  *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page   *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

func (p *AdminShopQueryReq) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()

	if p.Page != nil {
		query.WithPage(*p.Page)
	}
	if p.Limit != nil {
		query.WithLimit(*p.Limit)
	}

	if p.Order != nil {
		query.WithOrder(valueobject.Order(*p.Order))
	}

	if p.SortBy != nil {
		query.WithSortBy(*p.SortBy)
	} else {
		query.WithSortBy("id")
	}

	if p.Email != nil {
		query.Condition("email", valueobject.ILike, *p.Email)
	}

	return query, nil
}

type AdminShopRes struct {
	Id     uint   `json:"id"`
	Email  string `json:"email"`
	RoleId uint   `json:"role_id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
}

func NewAdminShopRes(u *entity.User) *AdminShopRes {
	return &AdminShopRes{Id: u.Id, Email: u.Email, RoleId: uint(u.RoleId), Name: u.AdminContact.Name, Phone: u.AdminContact.Phone}
}

type AdminShopReq struct {
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required" json:"password"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required,phonenumberprefix,phonenumberlength"`
}

type AdminShopUpdateReq struct {
	Email string `binding:"required,email" json:"email"`
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required,phonenumberprefix,phonenumberlength"`
}

func (a *AdminShopReq) ToModel() (*entity.User, *entity.AdminContact) {
	return &entity.User{Email: a.Email, Password: a.Password}, &entity.AdminContact{Name: a.Name, Phone: a.Phone}
}

func (a *AdminShopUpdateReq) ToModel() (*entity.User, *entity.AdminContact) {
	return &entity.User{Email: a.Email}, &entity.AdminContact{Name: a.Name, Phone: a.Phone}
}
