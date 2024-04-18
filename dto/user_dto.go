package dto

import (
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
)

type UserProfileResponse struct {
	Name      string `json:"name"`
	Image     string `json:"image"`
	Email     string `json:"email"`
	Birthdate string `json:"dob"`
}

type ResetPasswordRequest struct {
	OldPassword string `binding:"required" json:"old_password"`
	NewPassword string `binding:"required" json:"new_password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type UpdateProfileRequest struct {
	Name string `binding:"required" json:"name"`
}

func ToUserProfileDTO(user *entity.User, profile *entity.Profile) UserProfileResponse {
	var userProfileResponse UserProfileResponse
	userProfileResponse.Name = profile.Name
	userProfileResponse.Image = profile.Image
	userProfileResponse.Email = user.Email
	userProfileResponse.Birthdate = profile.Birthdate.Format("2006-01-02")
	return userProfileResponse
}

func (r *UpdateProfileRequest) ToProfile() *entity.Profile {
	return &entity.Profile{
		Name: r.Name,
	}
}

type UserQueryParamReq struct {
	Email      *string `form:"email"`
	IsVerified *bool   `form:"is_verified"`
	SortBy     *string `form:"sort_by" binding:"omitempty,oneof=email"`
	Order      *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit      *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page       *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

func (p *UserQueryParamReq) ToQuery() *valueobject.Query {
	query := valueobject.NewQuery()

	if p.Page != nil {
		query.WithPage(*p.Page)
	}
	if p.Limit != nil {
		query.WithLimit(*p.Limit)
	}

	if p.Order != nil {
		query.WithOrder(valueobject.Order(*p.Order))
	}

	if p.SortBy != nil {
		query.WithSortBy(*p.SortBy)
	} else {
		query.WithSortBy("id")
	}

	if p.Email != nil {
		query.Condition("email", valueobject.ILike, *p.Email)
	}
	if p.IsVerified != nil {
		query.Condition("is_verified", valueobject.Equal, *p.IsVerified)
	}

	return query
}

type UserRes struct {
	Id         uint   `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

type UpdateReq struct {
	Name string `json:"name" binding:"required"`
}
