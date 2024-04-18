package usecase

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/imagehelper"
	"github.com/ziadrahmatullah/ordent-test/repository"
	"github.com/ziadrahmatullah/ordent-test/transactor"
	"github.com/ziadrahmatullah/ordent-test/util"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type OrderUsecase interface {
	CreateOrder(context.Context, *entity.ProductOrder) (uint, error)
	ListAllOrders(context.Context, *valueobject.Query) (*valueobject.PagedResult, error)
	GetAvailableProduct(ctx context.Context) (decimal.Decimal, []*entity.ShopProduct, []*entity.OrderItem, []*entity.CartItem, error)
	UploadPaymentProof(context.Context, uint) error
	OrderDetail(context.Context, uint) (*entity.ProductOrder, []*entity.OrderItem, *entity.Address, error)
	UserUpdateOrderStatus(context.Context, *entity.ProductOrder) error
	AdminUpdateOrderStatus(context.Context, *entity.ProductOrder) error
}

type orderUsecase struct {
	manager         transactor.Manager
	imageHelper     imagehelper.ImageHelper
	cartRepo        repository.CartRepository
	orderItemRepo   repository.OrderItemRepository
	orderRepo       repository.OrderRepository
	cartItemRepo    repository.CartItemRepository
	addressRepo     repository.AddressRepository
	shopRepo        repository.ShopRepository
	shopProductRepo repository.ShopProductRepository
}

func NewOrderUsecase(
	manager transactor.Manager,
	imageHelper imagehelper.ImageHelper,
	cartRepo repository.CartRepository,
	orderItemRepo repository.OrderItemRepository,
	orderRepo repository.OrderRepository,
	cartItemRepo repository.CartItemRepository,
	addressRepo repository.AddressRepository,
	shopRepo repository.ShopRepository,
	shopProductRepo repository.ShopProductRepository,
) OrderUsecase {
	return &orderUsecase{
		manager:         manager,
		imageHelper:     imageHelper,
		cartRepo:        cartRepo,
		orderItemRepo:   orderItemRepo,
		orderRepo:       orderRepo,
		cartItemRepo:    cartItemRepo,
		addressRepo:     addressRepo,
		shopRepo:        shopRepo,
		shopProductRepo: shopProductRepo,
	}
}
func (u *orderUsecase) CreateOrder(ctx context.Context, userOrder *entity.ProductOrder) (uint, error) {
	userId := ctx.Value("user_id").(uint)
	total, shopProducts, orderItems, fetchedCartItem, err := u.GetAvailableProduct(ctx)
	if err != nil {
		return 0, err
	}
	var orderId uint
	err = u.manager.Run(ctx, func(c context.Context) error {
		var order entity.ProductOrder
		order.OrderedAt = time.Now()
		order.OrderStatusId = uint(entity.WaitingForPayment)
		order.ProfileId = userId
		order.ExpiredAt = time.Now().Add(time.Hour * 24).Truncate(time.Hour).Add(-time.Minute)
		order.ShippingName = userOrder.ShippingName
		order.ShippingPrice = userOrder.ShippingPrice
		order.ShippingEta = userOrder.ShippingEta
		order.TotalPayment = total.Add(userOrder.ShippingPrice)
		order.PaymentMethod = userOrder.PaymentMethod
		order.AddressId = userOrder.AddressId
		order.ItemOrderQty = len(shopProducts)
		createdOrder, err := u.orderRepo.Create(c, &order)
		if err != nil {
			return err
		}
		for _, item := range orderItems {
			item.OrderId = createdOrder.Id
		}
		err = u.orderItemRepo.BulkCreate(c, orderItems)
		if err != nil {
			return err
		}
		err = u.cartItemRepo.BulkDelete(c, fetchedCartItem)
		if err != nil {
			return err
		}
		orderId = createdOrder.Id
		return nil
	})
	return orderId, err
}

func (u *orderUsecase) GetAvailableProduct(ctx context.Context) (decimal.Decimal, []*entity.ShopProduct, []*entity.OrderItem, []*entity.CartItem, error) {
	userId := ctx.Value("user_id").(uint)
	cartItemQuery := valueobject.NewQuery().
		Condition("cart_id", valueobject.Equal, userId).
		Condition("is_checked", valueobject.Equal, true)
	fetchedCartItem, err := u.cartItemRepo.Find(ctx, cartItemQuery)
	if err != nil {
		return decimal.Zero, nil, nil, nil, err
	}
	if len(fetchedCartItem) == 0 {
		return decimal.Zero, nil, nil, nil, apperror.NewClientError(apperror.NewResourceStateError("no item in cart"))
	}
	var totalPrice decimal.Decimal
	var orderItems []*entity.OrderItem
	var listOfProductShopId []uint
	for _, value := range fetchedCartItem {
		orderItem := &entity.OrderItem{
			ShopProductId: value.Id,
			Quantity:      value.Quantity,
			SubTotal:      value.SubAmount,
		}
		orderItems = append(orderItems, orderItem)
		totalPrice = totalPrice.Add(orderItem.SubTotal)
		listOfProductShopId = append(listOfProductShopId, value.Id)
	}
	productShopQuery := valueobject.NewQuery().Condition("id", valueobject.In, listOfProductShopId).WithPreload("Product")
	fetchedPPResult, err := u.shopProductRepo.Find(ctx, productShopQuery)
	if err != nil {
		return decimal.Zero, nil, nil, nil, err
	}
	return totalPrice, fetchedPPResult, orderItems, fetchedCartItem, nil
}

func (u *orderUsecase) ListAllOrders(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	userId := ctx.Value("user_id").(uint)
	roleId := ctx.Value("role_id").(entity.RoleId)
	return u.orderRepo.FindAllOrders(ctx, query, userId, roleId)
}

func (u *orderUsecase) UploadPaymentProof(ctx context.Context, orderId uint) error {
	userId := ctx.Value("user_id").(uint)
	orderQuery := valueobject.NewQuery().Condition("id", valueobject.Equal, orderId)
	fetchedOrder, err := u.orderRepo.FindOne(ctx, orderQuery)
	if err != nil {
		return err
	}
	if fetchedOrder == nil {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", orderId))
	}
	if fetchedOrder.ProfileId != userId {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", orderId))
	}
	if fetchedOrder.ExpiredAt.Before(time.Now()) {
		return apperror.NewClientError(apperror.NewResourceStateError("order already expired"))
	}
	if fetchedOrder.OrderStatusId > uint(entity.WaitingForPaymentConfirmation) {
		return apperror.NewClientError(apperror.NewResourceStateError("payment already confirmed"))
	}
	image := ctx.Value("image")
	if image == nil {
		return apperror.NewClientError(apperror.NewResourceStateError("no image inputted"))
	}
	proofKey := fetchedOrder.ProofKey
	if fetchedOrder.ProofKey == "" {
		proofKey = entity.PaymentProofPrefix + util.GenerateRandomString(10)
		fetchedOrder.ProofKey = proofKey
	}
	proofUrl, err := u.imageHelper.Upload(ctx, image.(multipart.File), entity.PaymentProofFolder, proofKey)
	if err != nil {
		return err
	}
	fetchedOrder.PaymentProof = proofUrl
	fetchedOrder.OrderStatusId = uint(entity.WaitingForPaymentConfirmation)
	_, err = u.orderRepo.Update(ctx, fetchedOrder)
	if err != nil {
		return err
	}
	return nil
}

func (u *orderUsecase) OrderDetail(ctx context.Context, orderId uint) (*entity.ProductOrder, []*entity.OrderItem, *entity.Address, error) {
	userId := ctx.Value("user_id").(uint)
	roleId := ctx.Value("role_id").(entity.RoleId)
	fetchedOrder, err := u.orderRepo.FindOrderDetail(ctx, orderId, userId, roleId)
	if err != nil {
		return nil, nil, nil, err
	}
	if fetchedOrder.Id == uint(0) {
		return nil, nil, nil, apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", orderId))
	}
	orderItemQuery := valueobject.NewQuery().Condition("order_id", valueobject.Equal, orderId).WithJoin("ShopProduct").WithJoin("ShopProduct.Product")
	fetchedItemOrder, err := u.orderItemRepo.Find(ctx, orderItemQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	addressQuery := valueobject.NewQuery().Condition("id", valueobject.Equal, fetchedOrder.AddressId).WithPreload("Province").WithPreload("City")
	fetchedAddress, err := u.addressRepo.FindOne(ctx, addressQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	if fetchedAddress == nil {
		return nil, nil, nil, apperror.NewClientError(apperror.NewResourceNotFoundError("address", "id", fetchedOrder.AddressId))
	}
	return fetchedOrder, fetchedItemOrder, fetchedAddress, nil
}

func (u *orderUsecase) UserUpdateOrderStatus(ctx context.Context, order *entity.ProductOrder) error {
	userId := ctx.Value("user_id").(uint)
	orderQuery := valueobject.NewQuery().Condition("id", valueobject.Equal, order.Id)
	fetchedOrder, err := u.orderRepo.FindOne(ctx, orderQuery)
	if err != nil {
		return err
	}
	if fetchedOrder == nil {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", order.Id))
	}
	if fetchedOrder.ProfileId != userId {
		return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", order.Id))
	}
	if order.OrderStatusId == uint(entity.Canceled) {
		if fetchedOrder.OrderStatusId >= uint(entity.WaitingForPaymentConfirmation) {
			return apperror.NewClientError(apperror.NewResourceStateError("cant cancel order"))
		}
		fetchedOrder.OrderStatusId = uint(entity.Canceled)
	}
	if order.OrderStatusId == uint(entity.OrderConfirmed) {
		if fetchedOrder.OrderStatusId != uint(entity.Sent) {
			return apperror.NewClientError(apperror.NewResourceStateError("cant confirm order"))
		}
		fetchedOrder.OrderStatusId = uint(entity.OrderConfirmed)
	}
	_, err = u.orderRepo.Update(ctx, fetchedOrder)
	if err != nil {
		return err
	}
	return nil
}

func (u *orderUsecase) AdminUpdateOrderStatus(ctx context.Context, order *entity.ProductOrder) error {
	roleId := ctx.Value("role_id").(entity.RoleId)
	userId := ctx.Value("user_id").(uint)
	if roleId != entity.RoleAdmin {
		return apperror.NewForbiddenActionError("you're not admin")
	}
	err := u.manager.Run(ctx, func(c context.Context) error {
		orderQuery := valueobject.NewQuery().Condition("id", valueobject.Equal, order.Id).Lock()
		fetchedOrder, err := u.orderRepo.FindOne(c, orderQuery)
		if err != nil {
			return err
		}
		if fetchedOrder == nil {
			return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", order.Id))
		}
		switch order.OrderStatusId {
		case uint(entity.Processed):
			if fetchedOrder.OrderStatusId != uint(entity.WaitingForPaymentConfirmation) {
				return apperror.NewClientError(apperror.NewResourceStateError("cant update order status to processed"))
			}
			fetchedOrderItem, err := u.orderItemRepo.ListOfOrderItem(c, order.Id, userId)
			if err != nil {
				return err
			}
			if len(fetchedOrderItem) == 0 {
				return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", order.Id))
			}
			for _, item := range fetchedOrderItem {
				if item.ShopProduct.Stock >= item.Quantity {
					item.ShopProduct.Stock -= item.Quantity
					_, err := u.shopProductRepo.Update(c, &item.ShopProduct)
					if err != nil {
						return err
					}
				} else {
					return apperror.NewResourceStateError("Not enough stock, order cancelled")
				}
			}
			fetchedOrder.OrderStatusId = uint(entity.Processed)
		case uint(entity.Sent):
			if fetchedOrder.OrderStatusId != uint(entity.Processed) {
				return apperror.NewClientError(apperror.NewResourceStateError("cant update order status to sent"))
			}
			fetchedOrder.OrderStatusId = uint(entity.Sent)
		case uint(entity.Canceled):
			if fetchedOrder.OrderStatusId >= uint(entity.Sent) {
				return apperror.NewClientError(apperror.NewResourceStateError("cant cancel order"))
			}
			if fetchedOrder.OrderStatusId == uint(entity.Processed) {
				fetchedOrderItem, err := u.orderItemRepo.ListOfOrderItem(c, order.Id, userId)
				if err != nil {
					return err
				}
				if len(fetchedOrderItem) == 0 {
					return apperror.NewClientError(apperror.NewResourceNotFoundError("order", "id", order.Id))
				}
				for _, item := range fetchedOrderItem {
					item.ShopProduct.Stock += item.Quantity
					_, err := u.shopProductRepo.Update(c, &item.ShopProduct)
					if err != nil {
						return err
					}

				}
			}
			fetchedOrder.OrderStatusId = uint(entity.Canceled)
		}
		_, err = u.orderRepo.Update(c, fetchedOrder)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
