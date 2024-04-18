package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/ordent-test/dto"
	"github.com/ziadrahmatullah/ordent-test/usecase"
)

type ShippingMethodHandler struct {
	usecase usecase.ShippingMethodUsecase
}

func NewShippingMethodHandler(shippingMethodUsecase usecase.ShippingMethodUsecase) *ShippingMethodHandler {
	return &ShippingMethodHandler{usecase: shippingMethodUsecase}
}

func (h *ShippingMethodHandler) GetShippingMethod(c *gin.Context) {
	var addressUri dto.AddressUri
	if err := c.ShouldBindUri(&addressUri); err != nil {
		_ = c.Error(err)
		return
	}

	shippingMethods, err := h.usecase.GetShippingMethod(c.Request.Context(), addressUri.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var response []dto.ShippingMethodResponse
	for _, method := range shippingMethods {
		response = append(response, dto.ShippingMethodResponse{
			Name:     method.Name,
			Duration: method.EstimatedDuration,
			Cost:     method.Cost,
		})
	}
	c.JSON(http.StatusOK, dto.Response{Data: response})
}
