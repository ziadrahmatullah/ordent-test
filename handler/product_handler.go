package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ziadrahmatullah/ordent-test/dto"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/usecase"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase}
}

func (h *ProductHandler) ListProduct(c *gin.Context) {
	var request dto.ListProductQueryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}

	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}

	pagedResult, err := h.productUsecase.ListAllProduct(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	products := pagedResult.Data.([]*entity.Product)

	var response []*dto.ProductResponse
	for _, product := range products {
		response = append(response, dto.NewFromProduct(product))
	}
	c.JSON(200, dto.Response{
		Data:        response,
		CurrentPage: &pagedResult.CurrentPage,
		CurrentItem: &pagedResult.CurrentItems,
		TotalPage:   &pagedResult.TotalPage,
		TotalItem:   &pagedResult.TotalItem,
	})
}

func (h *ProductHandler) AddProduct(c *gin.Context) {
	var request dto.AddProductRequest
	if err := c.ShouldBindWith(&request, binding.Form); err != nil {
		_ = c.Error(err)
		return
	}
	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	product := request.ToProduct()

	createdProduct, err := h.productUsecase.AddProduct(c.Request.Context(), product)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Data: dto.NewFromProduct(createdProduct),
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var request dto.AddProductRequest
	var requestUri dto.RequestUri
	err := c.ShouldBindUri(&requestUri)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindWith(&request, binding.Form); err != nil {
		_ = c.Error(err)
		return
	}
	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	product := request.ToProduct()
	product.Id = uint(requestUri.Id)

	updatedProduct, err := h.productUsecase.UpdateProduct(c.Request.Context(), product)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Data: dto.NewFromProduct(updatedProduct),
	})
}

func (h *ProductHandler) GetProductDetail(c *gin.Context) {
	var uri dto.RequestUri

	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(err)
		return
	}

	fetchedProduct, err := h.productUsecase.GetProductDetail(c.Request.Context(), uint(uri.Id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Data: dto.NewFromShopProduct(fetchedProduct),
	})
}
