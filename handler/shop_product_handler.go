package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/ordent-test/dto"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/usecase"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type ShopProductHandler struct {
	shopProductUsecase usecase.ShopProductUsecase
}

func NewShopProductHandler(u usecase.ShopProductUsecase) *ShopProductHandler {
	return &ShopProductHandler{shopProductUsecase: u}
}

func (h *ShopProductHandler) GetAllShopProduct(c *gin.Context) {
	var request dto.ListShopProductQueryParam
	var requestUri dto.RequestShopUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	query.Condition("shop", valueobject.Equal, requestUri.Id)
	pagedResult, err := h.shopProductUsecase.FindAllShopProduct(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	tempData := []*dto.ProductShopRes{}
	for _, product := range pagedResult.Data.([]*entity.ShopProduct) {
		tempProduct := dto.NewProductPhamarcyRes(product)
		tempData = append(tempData, tempProduct)
	}
	c.JSON(http.StatusOK, dto.Response{
		Data:        tempData,
		CurrentPage: &pagedResult.CurrentPage,
		CurrentItem: &pagedResult.CurrentItems,
		TotalPage:   &pagedResult.TotalPage,
		TotalItem:   &pagedResult.TotalItem,
	})
}

func (h *ShopProductHandler) GetShopProductDetail(c *gin.Context) {
	var requestUri dto.ShopProductUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}

	shopProduct, err := h.shopProductUsecase.FindOneShopPeoduct(c.Request.Context(), &entity.ShopProduct{ProductId: requestUri.ProductId, ShopId: requestUri.ShopId})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{
		Data: dto.NewProductPhamarcyRes(shopProduct),
	})
}

func (h *ShopProductHandler) PostShopProduct(c *gin.Context) {
	var shopProductReq dto.ShopProductReq
	var requestUri dto.RequestShopUri
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	err := c.ShouldBindJSON(&shopProductReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	shopProduct, err := shopProductReq.ToModel()
	if err != nil {
		_ = c.Error(err)
		return
	}
	shopProduct.ShopId = requestUri.Id
	_, err = h.shopProductUsecase.CreateShopProduct(c.Request.Context(), shopProduct)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "created success"})
}

func (h *ShopProductHandler) PutShopProduct(c *gin.Context) {
	var shopProductUpdateReq dto.ShopProductUpdateReq
	var requestProductUri dto.ShopProductUri
	if err := c.ShouldBindUri(&requestProductUri); err != nil {
		_ = c.Error(err)
		return
	}
	err := c.ShouldBindJSON(&shopProductUpdateReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	shopProduct, err := shopProductUpdateReq.ToModel()
	if err != nil {
		_ = c.Error(err)
		return
	}
	shopProduct.ShopId = requestProductUri.ShopId
	shopProduct.Id = requestProductUri.ProductId
	_, err = h.shopProductUsecase.UpdateShopProduct(c.Request.Context(), shopProduct)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "updated success"})
}
