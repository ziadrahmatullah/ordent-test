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
	User           *handler.UserHandler
	Auth           *handler.AuthHandler
	Address        *handler.AddressHandler
	Product        *handler.ProductHandler
	Shop           *handler.ShopHandler
	ShopProduct    *handler.ShopProductHandler
	Cart           *handler.CartHandler
	ShippingMethod *handler.ShippingMethodHandler
	Order          *handler.OrderHandler
	StockRecord    *handler.StockRecordHandler
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
	shop.GET("/:shop_id", middleware.Auth(entity.RoleAdmin, entity.RoleSuperAdmin), handlers.Shop.GetShopDetail)
	shop.POST("", middleware.Auth(entity.RoleAdmin), handlers.Shop.AddShop)
	shop.PUT("/:shop_id", middleware.Auth(entity.RoleAdmin), handlers.Shop.UpdateShop)
	shop.DELETE("/:shop_id", handlers.Shop.DeleteShop)

	shopProduct := shop.Group("/:shop_id/products")
	shopProduct.GET("", middleware.Auth(entity.RoleAdmin), handlers.ShopProduct.GetAllShopProduct)
	shopProduct.POST("", middleware.Auth(entity.RoleAdmin), handlers.ShopProduct.PostShopProduct)
	shopProduct.GET("/:product_id", middleware.Auth(entity.RoleAdmin), handlers.ShopProduct.GetShopProductDetail)
	shopProduct.PUT("/:product_id", middleware.Auth(entity.RoleAdmin), handlers.ShopProduct.PutShopProduct)

	cart := router.Group("/cart")
	cart.GET("", middleware.Auth(entity.RoleUser), handlers.Cart.GetCart)
	cart.POST("", middleware.Auth(entity.RoleUser), handlers.Cart.AddItem)
	cart.PATCH("/:id", middleware.Auth(entity.RoleUser), handlers.Cart.ChangeQty)
	cart.DELETE("/:id", middleware.Auth(entity.RoleUser), handlers.Cart.DeleteItem)
	cart.PATCH("/check/:id", middleware.Auth(entity.RoleUser), handlers.Cart.CheckItem)
	cart.PUT("/check-all", middleware.Auth(entity.RoleUser), handlers.Cart.CheckAllItem)

	router.GET("/shipping-method/:id", middleware.Auth(entity.RoleUser), handlers.ShippingMethod.GetShippingMethod)

	order := router.Group("/order")
	order.POST("", middleware.Auth(entity.RoleUser), handlers.Order.CreateOrder)
	order.GET("", middleware.Auth(entity.RoleUser, entity.RoleAdmin, entity.RoleSuperAdmin), handlers.Order.OrderHistory)
	order.GET("/:id", middleware.Auth(entity.RoleUser, entity.RoleAdmin, entity.RoleSuperAdmin), handlers.Order.OrderDetail)
	order.GET("/items/:id", middleware.Auth(entity.RoleUser), handlers.Order.GetAvailableProduct)
	order.POST("/:id/upload-proof", middleware.Auth(entity.RoleUser), middleware.ImageUploadMiddleware(), handlers.Order.UploadPaymentProof)
	order.PATCH("/:id/status-admin", middleware.Auth(entity.RoleAdmin), handlers.Order.AdminUpdateOrderStatus)
	order.PATCH("/:id/status-user", middleware.Auth(entity.RoleUser), handlers.Order.UserUpdateOrderStatus)

	stockRecord := router.Group("/stock-records")
	stockRecord.GET("", middleware.Auth(entity.RoleAdmin), handlers.StockRecord.GetAllStockRecord)
	stockRecord.POST("", middleware.Auth(entity.RoleAdmin), handlers.StockRecord.PostStockRecord)

	return router
}

func routeNotFoundHandler(c *gin.Context) {
	var errRouteNotFound = errors.New("route not found")
	_ = c.Error(apperror.NewClientError(errRouteNotFound).NotFound())
}
