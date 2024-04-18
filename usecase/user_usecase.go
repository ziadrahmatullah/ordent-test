package usecase

import (
	"context"

	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/hasher"
	"github.com/ziadrahmatullah/ordent-test/repository"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type UserUsecase interface {
	GetAllUser(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	UserProfile(ctx context.Context) (*entity.User, *entity.Profile, error)
	ResetPassword(context.Context, string, string) error
	UpdateProfile(context.Context, *entity.Profile) error
}

type userUsecase struct {
	userRepo   repository.UserRepository
	profilRepo repository.ProfileRepository
	hash       hasher.Hasher
}

func NewUserUsecase(
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	hash hasher.Hasher,
) UserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		profilRepo: profileRepo,
		hash:       hash,
	}
}

func (u *userUsecase) GetAllUser(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.userRepo.FindAllUser(ctx, query)
}

func (u *userUsecase) UserProfile(ctx context.Context) (*entity.User, *entity.Profile, error) {
	userId := ctx.Value("user_id").(uint)
	fetchedUser, err := u.userRepo.FindById(ctx, userId)
	if err != nil {
		return nil, nil, err
	}
	if fetchedUser == nil {
		return nil, nil, apperror.NewClientError(apperror.NewInvalidCredentialsError())
	}
	uidQuery := valueobject.NewQuery().Condition("user_id", valueobject.Equal, userId)
	fetchProfile, err := u.profilRepo.FindOne(ctx, uidQuery)
	if err != nil {
		return nil, nil, err
	}
	return fetchedUser, fetchProfile, nil
}

func (u *userUsecase) ResetPassword(ctx context.Context, oldPassword, newPassword string) error {
	userId := ctx.Value("user_id").(uint)
	fetchedUser, err := u.userRepo.FindById(ctx, userId)
	if err != nil {
		return err
	}
	if fetchedUser == nil {
		return apperror.NewClientError(apperror.NewInvalidCredentialsError())
	}
	if oldPassword == newPassword {
		return apperror.NewClientError(apperror.NewResourceStateError("can't change to the same password"))
	}
	if !(u.hash.Compare(fetchedUser.Password, oldPassword)) {
		return apperror.NewClientError(apperror.NewResourceStateError("incorrect old password"))
	}
	hashedPassword, err := u.hash.Hash(newPassword)
	if err != nil {
		return err
	}
	fetchedUser.Password = string(hashedPassword)
	_, err = u.userRepo.Update(ctx, fetchedUser)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, profile *entity.Profile) error {
	userId := ctx.Value("user_id").(uint)
	fetchedUser, err := u.userRepo.FindById(ctx, userId)
	if err != nil {
		return err
	}
	if fetchedUser == nil {
		return apperror.NewClientError(apperror.NewInvalidCredentialsError())
	}
	profileQuery := valueobject.NewQuery().
		Condition("user_id", valueobject.Equal, fetchedUser.Id).Lock()
	fetchedProfile, err := u.profilRepo.FindOne(ctx, profileQuery)
	if err != nil {
		return err
	}
	fetchedProfile.Name = profile.Name
	_, err = u.profilRepo.Update(ctx, fetchedProfile)
	if err != nil {
		return err
	}
	return err
}
