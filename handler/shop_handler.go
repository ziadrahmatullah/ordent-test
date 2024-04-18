package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/ordent-test/dto"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/usecase"
)

type ShopHandler struct {
	shopUsecase usecase.ShopUsecase
}

func NewShopHandler(u usecase.ShopUsecase) *ShopHandler {
	return &ShopHandler{shopUsecase: u}
}

func (h *ShopHandler) GetAllShop(c *gin.Context) {
	var request dto.ListShopQueryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pageResult, err := h.shopUsecase.FindAllShop(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	shops := pageResult.Data.([]*entity.Shop)
	shopsRes := []*dto.ShopRes{}
	for _, shop := range shops {
		shopres := dto.NewShopRes(shop)
		shopsRes = append(shopsRes, shopres)
	}
	c.JSON(http.StatusOK, dto.Response{
		Data:        shopsRes,
		TotalPage:   &pageResult.TotalPage,
		TotalItem:   &pageResult.TotalItem,
		CurrentPage: &pageResult.CurrentPage,
		CurrentItem: &pageResult.CurrentItems,
	})
}

func (h *ShopHandler) GetShopDetail(c *gin.Context) {
	var requestUri dto.RequestShopUri
	err := c.ShouldBindUri(&requestUri)
	if err != nil {
		_ = c.Error(err)
		return
	}
	shop, err := h.shopUsecase.FindOneShopDetail(c.Request.Context(), &entity.Shop{Id: uint(requestUri.Id)})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.NewShopRes(shop)})
}

func (h *ShopHandler) AddShop(c *gin.Context) {
	var shopReq dto.ShopReq
	err := c.ShouldBindJSON(&shopReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	shop, err := shopReq.ToModel()
	if err != nil {
		_ = c.Error(err)
		return
	}
	_, err = h.shopUsecase.CreateShop(c.Request.Context(), shop)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "created success"})
}

func (h *ShopHandler) UpdateShop(c *gin.Context) {
	var shopReq dto.ShopReq
	var requestUri dto.RequestShopUri

	err := c.ShouldBindUri(&requestUri)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = c.ShouldBindJSON(&shopReq)
	if err != nil {
		_ = c.Error(err)
		return
	}

	shop, err := shopReq.ToModel()
	if err != nil {
		_ = c.Error(err)
		return
	}

	shop.Id = uint(requestUri.Id)
	_, err = h.shopUsecase.UpdateShop(c.Request.Context(), shop)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "update success"})
}

func (h *ShopHandler) DeleteShop(c *gin.Context) {
	var requestUri dto.RequestShopUri
	var shop entity.Shop
	err := c.ShouldBindUri(&requestUri)
	if err != nil {
		_ = c.Error(err)
		return
	}

	shop.Id = uint(requestUri.Id)
	err = h.shopUsecase.DeleteShop(c.Request.Context(), shop)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "delete success"})
}
