package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/dto"
	"github.com/ziadrahmatullah/ordent-test/usecase"
)

type AuthHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(u usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: u,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	user := request.ToUser()
	token, err := h.usecase.Register(c.Request.Context(), user)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
}

func (h *AuthHandler) Verify(c *gin.Context) {
	var request dto.VerifyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	token := c.Query("token")
	if token == "" {
		err := apperror.NewInvalidPathQueryParamError(apperror.NewInvalidTokenError())
		_ = c.Error(err)
		return
	}
	user := request.ToUser(token)
	profile, err := request.ToProfile()
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = h.usecase.Verify(c.Request.Context(), user, profile)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "verify success"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	user := request.ToUser()
	tokenUser, err := h.usecase.Login(c.Request.Context(), user)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.LoginResponse{Token: tokenUser.Token}})
}

func (h *AuthHandler) RequestForgotPassword(c *gin.Context) {
	var request dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	user := request.ToUser()
	token := dto.ToForgotPasswordEntity()
	token, err := h.usecase.ForgotPassword(c.Request.Context(), user, token)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.TokenResponse{Token: token.Token})
}

func (h *AuthHandler) ApplyPassword(c *gin.Context) {
	var request dto.ApplyPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	user := request.ToUser()
	token := dto.ToTokenEntity(c.Query("token"))
	err := h.usecase.ResetPassword(c.Request.Context(), user, token)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "password changed"})
}
