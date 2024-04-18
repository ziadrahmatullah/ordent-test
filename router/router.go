package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/handler"
	"github.com/ziadrahmatullah/ordent-test/middleware"
)

type Handlers struct {
	User *handler.UserHandler
	Auth *handler.AuthHandler
}

func New(handlers Handlers) http.Handler {
	router := gin.New()

	router.NoRoute(routeNotFoundHandler)
	router.Use(gin.Recovery())
	router.Use(middleware.Timeout())
	router.Use(middleware.Logger())
	router.Use(middleware.Error())

	auth := router.Group("/auth")
	auth.POST("/register", handlers.Auth.Register)
	auth.POST("/verify", handlers.Auth.Verify)
	auth.POST("/login", handlers.Auth.Login)
	auth.POST("/forgot-password", handlers.Auth.RequestForgotPassword)
	auth.PUT("/forgot-password", handlers.Auth.ApplyPassword)

	user := router.Group("/users")
	user.GET("", handlers.User.GetAllUser)

	return router
}

func routeNotFoundHandler(c *gin.Context) {
	var errRouteNotFound = errors.New("route not found")
	_ = c.Error(apperror.NewClientError(errRouteNotFound).NotFound())
}
