package usecase

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/imagehelper"
	"github.com/ziadrahmatullah/ordent-test/repository"
	"github.com/ziadrahmatullah/ordent-test/transactor"
	"github.com/ziadrahmatullah/ordent-test/util"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type ProductUsecase interface {
	ListAllProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	AddProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	GetProductDetail(ctx context.Context, productId uint) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
}

type productUsecase struct {
	manager         transactor.Manager
	imageHelper     imagehelper.ImageHelper
	productRepo     repository.ProductRepository
	categoryRepo    repository.ProductCategoryRepository
	shopProductRepo repository.ShopProductRepository
}

func NewProductUsecase(
	manager transactor.Manager,
	imageHelper imagehelper.ImageHelper,
	productRepo repository.ProductRepository,
	categoryRepo repository.ProductCategoryRepository,
	shopProductRepo repository.ShopProductRepository,
) ProductUsecase {
	return &productUsecase{
		manager:         manager,
		imageHelper:     imageHelper,
		productRepo:     productRepo,
		categoryRepo:    categoryRepo,
		shopProductRepo: shopProductRepo,
	}
}

func (u *productUsecase) ListAllProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	query.WithJoin("ProductCategory")
	pagedResult, err := u.productRepo.FindAllProducts(ctx, query)
	if err != nil {
		return nil, err
	}

	return pagedResult, nil
}

func (u *productUsecase) AddProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	var createdProduct *entity.Product

	fetchedProductCategory, err := u.categoryRepo.FindById(ctx, product.ProductCategoryId)
	if err != nil {
		return nil, err
	}
	if fetchedProductCategory == nil {
		return nil, apperror.NewResourceNotFoundError("product category", "id", product.ProductCategoryId)
	}

	image := ctx.Value("image")
	imageKey := entity.ProductKeyPrefix + util.GenerateRandomString(10)
	imgUrl, err := u.imageHelper.Upload(ctx, image.(multipart.File), entity.ProductFolder, imageKey)

	if err != nil {
		return nil, err
	}

	product.Image = imgUrl
	product.ImageKey = imageKey
	createdProduct, err = u.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (u *productUsecase) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	var updatedProduct *entity.Product
	fetchedProduct, err := u.productRepo.FindById(ctx, product.Id)
	if err != nil {
		return nil, err
	}
	if fetchedProduct == nil {
		return nil, apperror.NewResourceNotFoundError("product", "id", product.Id)
	}
	fetchedProductCategory, err := u.categoryRepo.FindById(ctx, product.ProductCategoryId)
	if err != nil {
		return nil, err
	}
	if fetchedProductCategory == nil {
		return nil, apperror.NewClientError(fmt.Errorf("product category id :%v not found", product.ProductCategoryId))
	}

	product.Image = fetchedProduct.Image
	product.ImageKey = fetchedProduct.ImageKey
	err = u.manager.Run(ctx, func(c context.Context) error {
		image := ctx.Value("image")
		if image != nil {
			imgUrl, err := u.imageHelper.Upload(ctx, image.(multipart.File), entity.ProductFolder, entity.ProductCategoryKeyPrefix+util.GenerateRandomString(10))
			if err != nil {
				return err
			}
			product.Image = imgUrl
		}
		updatedProduct, err = u.productRepo.Update(c, product)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (u *productUsecase) GetProductDetail(ctx context.Context, productId uint) (*entity.Product, error) {
	query := valueobject.NewQuery().
		Condition("\"products\".id", valueobject.Equal, productId)

	fetchedProduct, err := u.productRepo.FindOne(ctx, query)
	if err != nil {
		return nil, err
	}
	if fetchedProduct == nil {
		return nil, apperror.NewResourceNotFoundError("product", "id", productId)
	}

	return fetchedProduct, nil
}
