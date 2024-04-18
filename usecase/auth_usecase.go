package usecase

import (
	"context"
	"time"

	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/appjwt"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/hasher"
	"github.com/ziadrahmatullah/ordent-test/repository"
	"github.com/ziadrahmatullah/ordent-test/transactor"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type AuthUsecase interface {
	Register(context.Context, *entity.User) (string, error)
	Verify(context.Context, *entity.User, *entity.Profile) error
	Login(context.Context, *entity.User) (*entity.User, error)
	ForgotPassword(ctx context.Context, user *entity.User, tokenEntity *entity.ForgotPasswordToken) (token *entity.ForgotPasswordToken, err error)
	ResetPassword(context.Context, *entity.User, *entity.ForgotPasswordToken) error
}

type authUsecase struct {
	manager            transactor.Manager
	userRepo           repository.UserRepository
	profileRepo        repository.ProfileRepository
	forgotPasswordRepo repository.ForgotPasswordRepository
	cartRepo           repository.CartRepository
	hash               hasher.Hasher
	jwt                appjwt.Jwt
}

func NewAuthUsecase(
	manager transactor.Manager,
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	forgotPasswordRepo repository.ForgotPasswordRepository,
	cartRepo repository.CartRepository,
	hash hasher.Hasher,
	jwt appjwt.Jwt,
) AuthUsecase {
	return &authUsecase{
		manager:            manager,
		userRepo:           userRepo,
		profileRepo:        profileRepo,
		forgotPasswordRepo: forgotPasswordRepo,
		cartRepo:           cartRepo,
		hash:               hash,
		jwt:                jwt,
	}
}

func (u *authUsecase) Register(ctx context.Context, user *entity.User) (newToken string, err error) {
	emailQuery := valueobject.NewQuery().Condition("email", valueobject.Equal, user.Email)
	fetchedUser, err := u.userRepo.FindOne(ctx, emailQuery)
	if err != nil {
		return "", err
	}
	var token string
	if fetchedUser != nil {
		if fetchedUser.IsVerified {
			return "", apperror.NewResourceAlreadyExistError("user", "email", user.Email)
		}
	} else {
		hashedToken, err := u.hash.Hash(user.Email)
		if err != nil {
			return "", err
		}
		token = string(hashedToken)
		user.Token = token
		_, err = u.userRepo.Create(ctx, user)
		if err != nil {
			return "", err
		}
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *authUsecase) Verify(ctx context.Context, user *entity.User, profile *entity.Profile) error {
	err := u.manager.Run(ctx, func(c context.Context) error {
		tokenQuery := valueobject.NewQuery().
			Condition("token", valueobject.Equal, user.Token).Lock()
		fetchedUser, err := u.userRepo.FindOne(c, tokenQuery)
		if err != nil {
			return err
		}
		if fetchedUser == nil {
			return apperror.NewClientError(apperror.NewInvalidTokenError()).BadRequest()
		}
		if fetchedUser.IsVerified {
			return apperror.NewClientError(apperror.NewResourceStateError("Already Verified"))
		}
		hashPass, err := u.hash.Hash(user.Password)
		if err != nil {
			return err
		}
		fetchedUser.Password = string(hashPass)
		fetchedUser.IsVerified = true
		updatedUser, err := u.userRepo.Update(c, fetchedUser)
		if err != nil {
			return err
		}

		profile.UserId = updatedUser.Id
		_, err = u.profileRepo.Create(c, profile)
		if err != nil {
			return err
		}
		if fetchedUser.Role == entity.RoleUser {
			var cart entity.Cart
			cart.UserId = updatedUser.Id
			_, err = u.cartRepo.Create(c, &cart)
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	})
	return err
}

func (u *authUsecase) Login(ctx context.Context, user *entity.User) (*entity.User, error) {
	emailQuery := valueobject.NewQuery().Condition("email", valueobject.Equal, user.Email)
	fetchedUser, err := u.userRepo.FindOne(ctx, emailQuery)
	if err != nil {
		return nil, err
	}
	if fetchedUser == nil {
		return nil, apperror.NewResourceNotFoundError("user", "email", user.Email)
	}
	if !(u.hash.Compare(fetchedUser.Password, user.Password)) {
		return nil, apperror.NewInvalidCredentialsError()
	}
	token, err := u.jwt.GenerateToken(fetchedUser)
	if err != nil {
		return nil, err
	}
	fetchedUser.Token = token
	return fetchedUser, nil
}

func (u *authUsecase) ForgotPassword(ctx context.Context, user *entity.User, tokenEntity *entity.ForgotPasswordToken) (token *entity.ForgotPasswordToken, err error) {
	emailQuery := valueobject.NewQuery().Condition("email", valueobject.Equal, user.Email)
	fetchedUser, err := u.userRepo.FindOne(ctx, emailQuery)
	if err != nil {
		return nil, err
	}
	if fetchedUser == nil {
		return nil, apperror.NewResourceNotFoundError("user", "email", user.Email)
	}
	if !fetchedUser.IsVerified {
		return nil, apperror.NewResourceNotFoundError("user", "email", user.Email)
	}
	hashedToken, err := u.hash.Hash(user.Email)
	if err != nil {
		return nil, err
	}
	tokenEntity.Token = string(hashedToken)
	tokenEntity.UserId = fetchedUser.Id
	tokenEntity, err = u.forgotPasswordRepo.Create(ctx, tokenEntity)
	if err != nil {
		return nil, err
	}
	return tokenEntity, nil
}

func (u *authUsecase) ResetPassword(ctx context.Context, user *entity.User, tokenEntity *entity.ForgotPasswordToken) error {
	err := u.manager.Run(ctx, func(c context.Context) error {
		tokenQuery := valueobject.NewQuery().
			Condition("token", valueobject.Equal, tokenEntity.Token).Lock()
		fetchedTokenEntity, err := u.forgotPasswordRepo.FindOne(c, tokenQuery)
		if err != nil {
			return err
		}
		if fetchedTokenEntity == nil {
			return apperror.NewClientError(apperror.NewInvalidTokenError()).BadRequest()
		}
		if !fetchedTokenEntity.IsActive {
			return apperror.NewInvalidTokenError()
		}
		if fetchedTokenEntity.ExpiredAt.Before(time.Now()) {
			return apperror.NewInvalidTokenError()
		}
		userQuery := valueobject.NewQuery().
			Condition("id", valueobject.Equal, fetchedTokenEntity.UserId).Lock()
		fetchedUser, err := u.userRepo.FindOne(c, userQuery)
		if err != nil {
			return err
		}
		if u.hash.Compare(fetchedUser.Password, user.Password) {
			return apperror.NewClientError(apperror.NewResourceStateError("can't use the same password"))
		}
		hashPass, err := u.hash.Hash(user.Password)
		if err != nil {
			return err
		}
		fetchedUser.Password = string(hashPass)
		fetchedTokenEntity.IsActive = false
		_, err = u.userRepo.Update(c, fetchedUser)
		if err != nil {
			return err
		}
		_, err = u.forgotPasswordRepo.Update(c, fetchedTokenEntity)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}