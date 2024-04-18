package usecase

import (
	"context"
	"fmt"

	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/repository"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type ShopUsecase interface {
	FindAllShop(ctx context.Context, querry *valueobject.Query) (*valueobject.PagedResult, error)
	FindOneShopDetail(ctx context.Context, shop *entity.Shop) (*entity.Shop, error)
	CreateShop(ctx context.Context, shop *entity.Shop) (*entity.Shop, error)
	UpdateShop(ctx context.Context, shop *entity.Shop) (*entity.Shop, error)
	DeleteShop(ctx context.Context, shop entity.Shop) error
}

type shopUsecase struct {
	shopRepository     repository.ShopRepository
	provinceRepository repository.ProvinceRepository
	cityRepository     repository.CityRepository
}

func NewShopUsecase(rp repository.ShopRepository, pr repository.ProvinceRepository, cr repository.CityRepository) ShopUsecase {
	return &shopUsecase{shopRepository: rp, provinceRepository: pr, cityRepository: cr}
}

func (u *shopUsecase) FindAllShop(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.shopRepository.FindAllShop(ctx, query)
}

func (u *shopUsecase) FindOneShopDetail(ctx context.Context, shop *entity.Shop) (*entity.Shop, error) {
	selectShop, err := u.shopRepository.FindOne(ctx, valueobject.NewQuery().Condition("\"shops\".id", valueobject.Equal, shop.Id).WithJoin("Province").WithJoin("City"))
	if err != nil {
		return nil, err
	}
	if selectShop == nil {
		return nil, apperror.NewResourceNotFoundError("shop", "id", shop.Id)
	}
	return selectShop, nil
}

func (u *shopUsecase) CreateShop(ctx context.Context, shop *entity.Shop) (*entity.Shop, error) {
	province, err := u.provinceRepository.FindById(ctx, shop.ProvinceId)
	if err != nil {
		return nil, err
	}
	if province == nil {
		return nil, apperror.NewClientError(fmt.Errorf("province with id %v not found", shop.ProvinceId))
	}
	city, err := u.cityRepository.FindById(ctx, shop.CityId)
	if err != nil {
		return nil, err
	}

	if city == nil {
		return nil, apperror.NewClientError(fmt.Errorf("city with id %v not found", shop.CityId))
	}
	if province.Id != city.ProvinceId {
		return nil, apperror.NewClientError(fmt.Errorf("city with id %v doesn't belong to province id %v", city.Id, shop.ProvinceId))
	}
	shop.AdminId = ctx.Value("user_id").(uint)
	newShop, err := u.shopRepository.Create(ctx, shop)
	if err != nil {
		return nil, err
	}
	return newShop, nil
}

func (u *shopUsecase) UpdateShop(ctx context.Context, shop *entity.Shop) (*entity.Shop, error) {
	checkShop, err := u.shopRepository.FindById(ctx, shop.Id)
	if err != nil {
		return nil, err
	}

	if checkShop == nil {
		return nil, apperror.NewResourceNotFoundError("shop", "id", shop.Id)
	}
	province, err := u.provinceRepository.FindById(ctx, shop.ProvinceId)
	if err != nil {
		return nil, err
	}

	if province == nil {
		return nil, apperror.NewClientError(fmt.Errorf("province with id %v not found", shop.ProvinceId))
	}

	city, err := u.cityRepository.FindById(ctx, shop.CityId)
	if err != nil {
		return nil, err
	}

	if city == nil {
		return nil, apperror.NewClientError(fmt.Errorf("city with id %v not found", shop.CityId))
	}

	if province.Id != city.ProvinceId {
		return nil, apperror.NewClientError(fmt.Errorf("city with id %v doesn't belong to province id %v", city.Id, shop.ProvinceId))
	}

	if checkShop.AdminId != ctx.Value("user_id").(uint) {
		return nil, apperror.NewForbiddenActionError("cannot have access to update this shop")
	}

	shop.AdminId = checkShop.AdminId
	updatedShop, err := u.shopRepository.Update(ctx, shop)
	if err != nil {
		return nil, err
	}
	return updatedShop, nil
}

func (u *shopUsecase) DeleteShop(ctx context.Context, shop entity.Shop) error {
	checkShop, err := u.shopRepository.FindById(ctx, shop.Id)
	if err != nil {
		return err
	}
	if checkShop == nil {
		return apperror.NewResourceNotFoundError("product category", "id", shop.Id)
	}
	if checkShop.AdminId != ctx.Value("user_id").(uint) {
		return apperror.NewForbiddenActionError("cannot have access to delete this shop")
	}
	return u.shopRepository.Delete(ctx, &shop)
}
