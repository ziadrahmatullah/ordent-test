package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/handler"
	"github.com/ziadrahmatullah/ordent-test/middleware"
)

type Handlers struct {
	User    *handler.UserHandler
	Auth    *handler.AuthHandler
	Address *handler.AddressHandler
	Product *handler.ProductHandler
	Shop    *handler.ShopHandler
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
	user.GET("", middleware.Auth(entity.RoleAdmin), handlers.User.GetAllUser)
	user.GET("/profile", middleware.Auth(entity.RoleUser), handlers.User.GetProfile)
	user.POST("/reset-password", middleware.Auth(entity.RoleUser), handlers.User.ResetPassword)
	user.PUT("/profile", middleware.Auth(entity.RoleUser), handlers.User.UpdateProfile)

	address := router.Group("/addresses")
	address.GET("", middleware.Auth(entity.RoleUser), handlers.Address.GetAddress)
	address.POST("", middleware.Auth(entity.RoleUser), handlers.Address.CreateAddress)
	address.PUT("/:id", middleware.Auth(entity.RoleUser), handlers.Address.UpdateAddress)
	address.DELETE("/:id", middleware.Auth(entity.RoleUser), handlers.Address.DeleteAddress)
	address.PATCH("/:id/default", middleware.Auth(entity.RoleUser), handlers.Address.ChangeDefaultAddress)

	products := router.Group("/products")
	products.GET("", handlers.Product.ListProduct)
	products.POST("", middleware.Auth(entity.RoleSuperAdmin), middleware.ImageUploadMiddleware(), handlers.Product.AddProduct)
	products.PUT("/:id", middleware.Auth(entity.RoleSuperAdmin), middleware.ImageUploadMiddleware(), handlers.Product.UpdateProduct)
	products.GET("/:id", handlers.Product.GetProductDetail)

	shop := router.Group("/shops")
	shop.GET("", middleware.Auth(entity.RoleAdmin), handlers.Shop.GetAllShop)
	shop.GET("/:pharmacy_id", middleware.Auth(entity.RoleAdmin, entity.RoleSuperAdmin), handlers.Shop.GetShopDetail)
	shop.POST("", middleware.Auth(entity.RoleAdmin), handlers.Shop.AddShop)
	shop.PUT("/:pharmacy_id", middleware.Auth(entity.RoleAdmin), handlers.Shop.UpdateShop)
	shop.DELETE("/:pharmacy_id", handlers.Shop.DeleteShop)

	return router
}

func routeNotFoundHandler(c *gin.Context) {
	var errRouteNotFound = errors.New("route not found")
	_ = c.Error(apperror.NewClientError(errRouteNotFound).NotFound())
}
