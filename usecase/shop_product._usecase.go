package usecase

import (
	"context"
	"fmt"

	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/repository"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type ShopProductUsecase interface {
	FindAllShopProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	FindOneShopPeoduct(ctx context.Context, shopProduct *entity.ShopProduct) (*entity.ShopProduct, error)
	CreateShopProduct(ctx context.Context, shopProduct *entity.ShopProduct) (*entity.ShopProduct, error)
	UpdateShopProduct(ctx context.Context, shopProduct *entity.ShopProduct) (*entity.ShopProduct, error)
}

type shopProductUsecase struct {
	shopProductRepository repository.ShopProductRepository
	productRepository     repository.ProductRepository
	shopRepository        repository.ShopRepository
}

func NewShopProductUsecase(rp repository.ShopProductRepository, pr repository.ShopRepository, p repository.ProductRepository) ShopProductUsecase {
	return &shopProductUsecase{shopProductRepository: rp, productRepository: p, shopRepository: pr}
}

func (u *shopProductUsecase) FindAllShopProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	idShop := query.GetConditionValue("shop").(uint)
	checkShop, err := u.shopRepository.FindById(ctx, idShop)
	if err != nil {
		return nil, err
	}
	if checkShop == nil {
		return nil, apperror.NewResourceNotFoundError("shop", "id", idShop)
	}
	if checkShop.AdminId != ctx.Value("user_id").(uint) {
		return nil, apperror.NewForbiddenActionError("dont have access to this shop")
	}
	return u.shopProductRepository.FindAllShopProducts(ctx, query)
}

func (u *shopProductUsecase) FindOneShopPeoduct(ctx context.Context, shopProduct *entity.ShopProduct) (*entity.ShopProduct, error) {
	userId := ctx.Value("user_id").(uint)
	query := valueobject.NewQuery().
		Condition("\"shop_products\".product_id", valueobject.Equal, shopProduct.ProductId).
		Condition("\"shop_products\".shop_id", valueobject.Equal, shopProduct.ShopId).
		WithJoin("Product.ProductCategory").WithPreload("Shop")
	selectShopProduct, err := u.shopProductRepository.FindOne(ctx, query)
	if err != nil {
		return nil, err
	}
	if selectShopProduct == nil {
		return nil, apperror.NewResourceNotFoundError("shop product", "id", shopProduct.ProductId)
	}
	if selectShopProduct.Shop.AdminId != userId {
		return nil, apperror.NewResourceNotFoundError("shop product", "id", shopProduct.ProductId)
	}
	return selectShopProduct, nil
}

func (u *shopProductUsecase) CreateShopProduct(ctx context.Context, shopProduct *entity.ShopProduct) (*entity.ShopProduct, error) {
	product, err := u.productRepository.FindById(ctx, shopProduct.ProductId)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, apperror.NewClientError(fmt.Errorf("product with id %v not found", shopProduct.ProductId))
	}

	shop, err := u.shopRepository.FindById(ctx, shopProduct.ShopId)
	if err != nil {
		return nil, err
	}

	if shop == nil {
		return nil, apperror.NewResourceNotFoundError("shop", "id", shopProduct.ShopId)
	}

	checkPharProduct, err := u.shopProductRepository.FindOne(ctx, valueobject.NewQuery().Condition("shop_id", valueobject.Equal, shopProduct.ShopId).Condition("product_id", valueobject.Equal, shopProduct.ProductId))
	if err != nil {
		return nil, err
	}

	if checkPharProduct != nil {
		return nil, apperror.NewClientError(fmt.Errorf("cannot add duplicate product on this shop id %v", shopProduct.ShopId))
	}

	if shop.AdminId != ctx.Value("user_id").(uint) {
		return nil, apperror.NewForbiddenActionError("cannot have access to add product to this shop")
	}

	newShopProduct, err := u.shopProductRepository.Create(ctx, shopProduct)
	if err != nil {
		return nil, err
	}

	return newShopProduct, nil
}

func (u *shopProductUsecase) UpdateShopProduct(ctx context.Context, shopProduct *entity.ShopProduct) (*entity.ShopProduct, error) {
	shop, err := u.shopRepository.FindById(ctx, shopProduct.ShopId)
	if err != nil {
		return nil, err
	}

	if shop == nil {
		return nil, apperror.NewResourceNotFoundError("shop", "id", shopProduct.ShopId)
	}

	checkPharProduct, err := u.shopProductRepository.FindById(ctx, shopProduct.Id)
	if err != nil {
		return nil, err
	}

	if checkPharProduct == nil {
		return nil, apperror.NewResourceNotFoundError("shop product", "id", shopProduct.Id)
	}

	if shop.AdminId != ctx.Value("user_id").(uint) {
		return nil, apperror.NewForbiddenActionError("cannot have access to add profuct to this shop")
	}
	shopProduct.ProductId = checkPharProduct.ProductId
	shopProduct.Stock = checkPharProduct.Stock
	newShopProduct, err := u.shopProductRepository.Update(ctx, shopProduct)
	if err != nil {
		return nil, err
	}

	return newShopProduct, nil
}
