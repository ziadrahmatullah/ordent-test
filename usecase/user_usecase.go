package usecase

import (
	"context"
)

type UserUsecase interface {
	GetAllUser(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(
	userRepo repository.UserRepository,
) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetAllUser(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.userRepo.FindAllUser(ctx, query)
}
