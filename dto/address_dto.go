package dto

import (
	"github.com/ziadrahmatullah/ordent-test/entity"
)

type GetAddressResponse struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	StreetName string `json:"street"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
	Detail     string `json:"detail"`
	Province   string `json:"province"`
	City       string `json:"city"`
	IsDefault  bool   `json:"is_default"`
}

type CreateAddressRequest struct {
	Name       string `binding:"required" json:"name"`
	StreetName string `binding:"required" json:"street"`
	PostalCode string `binding:"required" json:"postal_code"`
	Phone      string `binding:"required,phonenumberprefix,phonenumberlength" json:"phone"`
	Detail     string `json:"detail"`
	Province   string `binding:"required" json:"province"`
	City       string `binding:"required" json:"city"`
}

type UpdateAddressRequest struct {
	Name       string `binding:"required" json:"name"`
	StreetName string `binding:"required" json:"street"`
	PostalCode string `binding:"required" json:"postal_code"`
	Phone      string `binding:"required,phonenumberprefix,phonenumberlength" json:"phone"`
	Detail     string `json:"detail"`
	Province   string `binding:"required" json:"province"`
	City       string `binding:"required" json:"city"`
}

type CreateAddressResponse struct {
	AddressId uint `json:"address_id"`
}

type AddressUri struct {
	Id uint `uri:"id" binding:"required,numeric"`
}

func (r *CreateAddressRequest) ToProfile() *entity.Profile {
	return &entity.Profile{
		Name: r.Name,
	}
}

func (r *CreateAddressRequest) ToAddress() *entity.Address {
	address := &entity.Address{
		Name:       r.Name,
		StreetName: r.StreetName,
		PostalCode: r.PostalCode,
		Phone:      r.Phone,
		Detail:     r.Detail,
		Province:   r.Province,
		City:       r.City,
	}
	return address
}

func (r *UpdateAddressRequest) ToAddress() *entity.Address {
	address := &entity.Address{
		Name:       r.Name,
		StreetName: r.StreetName,
		PostalCode: r.PostalCode,
		Phone:      r.Phone,
		Detail:     r.Detail,
		Province:   r.Province,
		City:       r.City,
	}
	return address
}
