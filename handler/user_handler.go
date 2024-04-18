package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsercase usecase.UserUsecase
}

func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsercase: u,
	}
}

func (h *UserHandler) GetAllUser(c *gin.Context) {
	var requestParam dto.UserQueryParamReq
	if err := c.ShouldBindQuery(&requestParam); err != nil {
		_ = c.Error(err)
		return
	}
	query := requestParam.ToQuery()
	pageResult, err := h.userUsercase.GetAllUser(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	users := pageResult.Data.([]*entity.User)
	c.JSON(http.StatusOK, dto.Response{
		Data:        users,
		TotalPage:   &pageResult.TotalPage,
		TotalItem:   &pageResult.TotalItem,
		CurrentPage: &pageResult.CurrentPage,
		CurrentItem: &pageResult.CurrentItems,
	})
}

