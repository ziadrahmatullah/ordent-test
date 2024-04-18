package usecase

import (
	"context"
	"mime/multipart"

	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/hasher"
	"github.com/ziadrahmatullah/ordent-test/imagehelper"
	"github.com/ziadrahmatullah/ordent-test/repository"
	"github.com/ziadrahmatullah/ordent-test/transactor"
	"github.com/ziadrahmatullah/ordent-test/util"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type UserUsecase interface {
	GetAllUser(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	UserProfile(ctx context.Context) (*entity.User, *entity.Profile, error)
	ResetPassword(context.Context, string, string) error
	UpdateProfile(context.Context, *entity.Profile) error
}

type userUsecase struct {
	userRepo    repository.UserRepository
	profilRepo  repository.ProfileRepository
	hash        hasher.Hasher
	imageHelper imagehelper.ImageHelper
	manager     transactor.Manager
}

func NewUserUsecase(
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	hash hasher.Hasher,
	imageHelper imagehelper.ImageHelper,
	manager transactor.Manager,
) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		profilRepo:  profileRepo,
		hash:        hash,
		imageHelper: imageHelper,
		manager:     manager,
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
	updatedProfileQuery := valueobject.NewQuery().
		Condition("user_id", valueobject.Equal, userId).Lock()
	updatedProfile, err := u.profilRepo.FindOne(ctx, updatedProfileQuery)
	if err != nil {
		return err
	}
	var imageKey string
	err = u.manager.Run(ctx, func(c context.Context) error {
		fetchedUser, err := u.userRepo.FindById(c, userId)
		if err != nil {
			return err
		}
		if fetchedUser == nil {
			return apperror.NewClientError(apperror.NewInvalidCredentialsError())
		}
		profileQuery := valueobject.NewQuery().
			Condition("user_id", valueobject.Equal, fetchedUser.Id).Lock()
		fetchedProfile, err := u.profilRepo.FindOne(c, profileQuery)
		if err != nil {
			return err
		}
		var imgUrl string
		image := c.Value("image")
		if image != nil {
			imageKey = fetchedProfile.ImageKey
			if fetchedProfile.ImageKey == "" {
				imageKey = entity.ProfilePhotoKeyPrefix + util.GenerateRandomString(10)
				fetchedProfile.ImageKey = imageKey
			}
			imgUrl, err = u.imageHelper.Upload(c, image.(multipart.File), entity.ProfilePhotoFolder, entity.ProfilePhotoKeyPrefix+util.GenerateRandomString(10))
			if err != nil {
				return err
			}
			fetchedProfile.Image = imgUrl
		}
		fetchedProfile.Name = profile.Name
		updatedProfile, err = u.profilRepo.Update(c, fetchedProfile)
		if err != nil {
			return err
		}
		if fetchedUser.RoleId == entity.RoleUser {
			return nil
		}
		return nil
	})
	if updatedProfile.Image == "" {
		u.imageHelper.Destroy(ctx, entity.ProfilePhotoFolder, imageKey)
	}
	return err
}
