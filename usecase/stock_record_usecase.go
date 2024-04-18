package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/logger"
	"github.com/ziadrahmatullah/ordent-test/repository"
	"github.com/ziadrahmatullah/ordent-test/transactor"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type StockRecordUsecase interface {
	FindAllStockRecord(ctx context.Context, querry *valueobject.Query) (*valueobject.PagedResult, error)
	CreateStockRecord(ctx context.Context, StockRecord *entity.StockRecord) (*entity.StockRecord, error)
}

type stockRecordUsecase struct {
	stockRecordRepository repository.StockRecordRepository
	shopProductRepository repository.ShopProductRepository
	manager               transactor.Manager
}

func NewStockRecordUsecase(
	stockRecordRepo repository.StockRecordRepository,
	shopProductRepo repository.ShopProductRepository,
	manager transactor.Manager,
) StockRecordUsecase {
	return &stockRecordUsecase{
		stockRecordRepository: stockRecordRepo,
		shopProductRepository: shopProductRepo,
		manager:               manager}
}

func (u *stockRecordUsecase) FindAllStockRecord(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.stockRecordRepository.FindAllStockRecord(ctx, query)
}

func (u *stockRecordUsecase) CreateStockRecord(ctx context.Context, stockRecord *entity.StockRecord) (*entity.StockRecord, error) {
	var newStockRecord *entity.StockRecord
	stockRecord.ChangeAt = time.Now()
	shopProduct, err := u.shopProductRepository.FindOne(ctx, valueobject.NewQuery().Condition("\"shop_products\".id", valueobject.Equal, stockRecord.ShopProductId).WithJoin("Shop"))
	if err != nil {
		return nil, err
	}
	if shopProduct == nil {
		return nil, apperror.NewClientError(fmt.Errorf("product with id %v not found", stockRecord.ShopProductId))
	}
	if shopProduct.Shop.AdminId != ctx.Value("user_id").(uint) {
		logger.Log.Info(ctx.Value("user_id").(uint))
		logger.Log.Info(shopProduct.Shop.AdminId)
		return nil, apperror.NewForbiddenActionError("cannot have access to change stock")
	}
	err = u.manager.Run(ctx, func(c context.Context) error {
		_, err = u.shopProductRepository.FindOne(c, valueobject.NewQuery().Condition("\"shop_products\".id", valueobject.Equal, stockRecord.ShopProductId).Lock())
		if err != nil {
			return err
		}

		newStockRecord, err = u.stockRecordRepository.Create(c, stockRecord)
		if err != nil {
			return err
		}
		number := int(shopProduct.Stock)
		if stockRecord.IsReduction {
			number -= int(stockRecord.Quantity)
		} else {
			number += int(stockRecord.Quantity)
		}
		if number < 0 {
			return apperror.NewClientError(errors.New("product's stock cannot below zero"))
		}
		shopProduct.Stock = number
		_, err = u.shopProductRepository.Update(c, shopProduct)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return newStockRecord, nil
}
