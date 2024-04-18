package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/ordent-test/apperror"
	"github.com/ziadrahmatullah/ordent-test/appjwt"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/util"
)

func Auth(roles ...entity.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		token, err := extractBearerToken(bearerToken)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			return
		}

		newJwt := appjwt.NewJwt()
		claims, err := newJwt.ValidateToken(token)
		if err != nil {
			c.Abort()
			_ = c.Error(apperror.NewInvalidTokenError())
			return
		}

		ctx := context.WithValue(c.Request.Context(), "user_id", claims.Id)
		ctx = context.WithValue(ctx, "role", claims.Role)
		c.Request = c.Request.WithContext(ctx)

		if !util.IsMemberOf(roles, claims.Role) {
			c.Abort()
			_ = c.Error(apperror.NewForbiddenActionError("permission denied"))
			return
		}

		c.Next()
	}
}

func extractBearerToken(bearerToken string) (string, error) {
	if bearerToken == "" {
		return "", apperror.NewMissingTokenError()
	}
	token := strings.Split(bearerToken, " ")
	if len(token) != 2 {
		return "", apperror.NewInvalidTokenError()
	}
	return token[1], nil
}
