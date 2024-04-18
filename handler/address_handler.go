package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/ordent-test/dto"
	"github.com/ziadrahmatullah/ordent-test/usecase"
)

type AddressHandler struct {
	usecase usecase.AddressUsecase
}

func NewAddressHandler(u usecase.AddressUsecase) *AddressHandler {
	return &AddressHandler{
		usecase: u,
	}
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var request dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	address := request.ToAddress()
	err := h.usecase.CreateAddress(c.Request.Context(), address)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "create success"})
}

func (h *AddressHandler) GetAddress(c *gin.Context) {
	addresses, err := h.usecase.GetAddress(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	var addressesDto []*dto.GetAddressResponse
	for _, address := range addresses {
		var addressDto dto.GetAddressResponse
		addressDto.Id = address.Id
		addressDto.Name = address.Name
		addressDto.StreetName = address.StreetName
		addressDto.PostalCode = address.PostalCode
		addressDto.Phone = address.Phone
		addressDto.Detail = address.Detail
		addressDto.Province = address.Province
		addressDto.City = address.City
		addressDto.IsDefault = address.IsDefault
		addressesDto = append(addressesDto, &addressDto)
	}
	c.JSON(http.StatusOK, dto.Response{Data: addressesDto})
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	var addressUri dto.AddressUri
	var request dto.UpdateAddressRequest
	if err := c.ShouldBindUri(&addressUri); err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	address := request.ToAddress()
	address.Id = addressUri.Id
	err := h.usecase.UpdateAddress(c.Request.Context(), address)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "address updated"})
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	var addressUri dto.AddressUri
	if err := c.ShouldBindUri(&addressUri); err != nil {
		_ = c.Error(err)
		return
	}
	err := h.usecase.DeleteAddress(c.Request.Context(), addressUri.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "address deleted"})
}

func (h *AddressHandler) ChangeDefaultAddress(c *gin.Context) {
	var addressUri dto.AddressUri
	if err := c.ShouldBindUri(&addressUri); err != nil {
		_ = c.Error(err)
		return
	}
	err := h.usecase.ChangeDefaultAddress(c.Request.Context(), addressUri.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "address default changed"})
}
