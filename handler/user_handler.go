package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/ordent-test/dto"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/usecase"
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

func (h *UserHandler) GetProfile(c *gin.Context) {
	user, profile, err := h.userUsercase.UserProfile(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToUserProfileDTO(user, profile))
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var request dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	err := h.userUsercase.ResetPassword(c.Request.Context(), request.OldPassword, request.NewPassword)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ResetPasswordResponse{Message: "password changed"})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var request dto.UpdateProfileRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	profile := request.ToProfile()
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = h.userUsercase.UpdateProfile(c.Request.Context(), profile)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "update success"})
}
