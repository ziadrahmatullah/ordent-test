package usecase

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/repository"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type ShippingMethodUsecase interface {
	GetShippingMethod(ctx context.Context, addressId uint) ([]*entity.CalculatedShippingMethod, error)
}

type shippingMethodUsecase struct {
	addressRepo  repository.AddressRepository
	shippingRepo repository.ShippingMethodRepository
	shopRepo     repository.ShopRepository
	orderUsecase OrderUsecase
}

func NewShippingMethodUsecase(
	addressRepo repository.AddressRepository,
	shippingRepo repository.ShippingMethodRepository,
	shopRepo repository.ShopRepository,
	orderUsecase OrderUsecase,
) ShippingMethodUsecase {
	return &shippingMethodUsecase{
		addressRepo:  addressRepo,
		shippingRepo: shippingRepo,
		shopRepo:     shopRepo,
		orderUsecase: orderUsecase,
	}
}

func (u *shippingMethodUsecase) GetShippingMethod(ctx context.Context, addressId uint) ([]*entity.CalculatedShippingMethod, error) {
	userId := ctx.Value("user_id").(uint)

	_, products, _, _, err := u.orderUsecase.GetAvailableProduct(ctx)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, apperror.NewClientError(fmt.Errorf("no shop available"))
	}

	addressQuery := valueobject.NewQuery().
		WithJoin("City").
		Condition("\"addresses\".id", valueobject.Equal, addressId)
	fetchedAddress, err := u.addressRepo.FindOne(ctx, addressQuery)
	if err != nil {
		return nil, err
	}
	if fetchedAddress == nil {
		return nil, apperror.NewResourceNotFoundError("address", "id", addressId)
	}

	if fetchedAddress.ProfileId != userId {
		return nil, apperror.NewForbiddenActionError("not the address of current logged in user")
	}

	fetchedShop, err := u.shopRepo.FindNearestShopFromAddress(ctx, addressId)
	if err != nil {
		return nil, err
	}
	if fetchedShop == nil || len(fetchedShop) < 1 {
		return nil, apperror.NewClientError(fmt.Errorf("there's no shop available near this address"))
	}

	distanceInKM, err := u.shippingRepo.FindDistanceBetween(ctx, fetchedShop[0].Location, fetchedAddress.Location)
	if err != nil {
		return nil, err
	}
	distanceThresholdInKM := decimal.NewFromInt(25)
	calculatedShippingMethods := make([]*entity.CalculatedShippingMethod, 0)

	if distanceInKM.LessThan(distanceThresholdInKM) {
		fetchedShippingMethod, err := u.shippingRepo.Find(ctx, valueobject.NewQuery())
		if err != nil {
			return nil, err
		}
		distanceInKM = distanceInKM.DivRound(decimal.NewFromInt(1), 0)
		for _, fsm := range fetchedShippingMethod {
			calculatedShippingMethods = append(calculatedShippingMethods, &entity.CalculatedShippingMethod{
				Name:              fsm.Name,
				EstimatedDuration: fsm.Duration,
				Cost:              fsm.PricePerKM.Mul(distanceInKM).String(),
			})
		}
	}

	calculatedThirdPartyShippingMethods, _ := u.shippingRepo.GetThirdPartyShipping(ctx, fetchedShop[0].City.Code, fetchedAddress.City.Code, "10")

	if len(calculatedThirdPartyShippingMethods) > 0 {
		calculatedShippingMethods = append(calculatedShippingMethods, calculatedThirdPartyShippingMethods...)
	}

	return calculatedShippingMethods, nil
}
