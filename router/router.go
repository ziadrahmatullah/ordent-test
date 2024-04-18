package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User *handler.UserHandler
}

func New(handlers Handlers) http.Handler {
	router := gin.New()

	router.NoRoute(routeNotFoundHandler)
	router.Use(gin.Recovery())
	router.Use(middleware.Timeout())
	router.Use(middleware.Logger())
	router.Use(middleware.Error())

	user := router.Group("/users")
	user.GET("", handlers.User.GetAllUser)

	return router
}

func routeNotFoundHandler(c *gin.Context) {
	var errRouteNotFound = errors.New("route not found")
	_ = c.Error(apperror.NewClientError(errRouteNotFound).NotFound())
}
