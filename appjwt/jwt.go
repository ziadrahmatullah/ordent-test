package appjwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ziadrahmatullah/ordent-test/config"
	"github.com/ziadrahmatullah/ordent-test/entity"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Id   uint        `json:"id"`
	Role entity.Role `json:"role_id"`
}

type Jwt interface {
	GenerateToken(user *entity.User) (string, error)
	ValidateToken(tokenString string) (*entity.User, error)
}

type jwtImpl struct {
	secretKey []byte
}

func NewJwt() Jwt {
	jwtConfig := config.NewJwtConfig()
	return &jwtImpl{
		secretKey: []byte(jwtConfig.Secret),
	}
}

func (j *jwtImpl) GenerateToken(user *entity.User) (string, error) {
	jwtConfig := config.NewJwtConfig()
	appConfig := config.NewAppConfig()

	userId := strconv.Itoa(int(user.Id))
	duration := jwtConfig.ExpiryDuration

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    appConfig.Name,
			Subject:   userId,
		},
		Id:   user.Id,
		Role: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(j.secretKey)
	return signedString, err
}

func (j *jwtImpl) ValidateToken(tokenString string) (*entity.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok {
		user := &entity.User{
			Id:   claims.Id,
			Role: claims.Role,
		}
		return user, nil
	}
	return nil, errors.New("invalid claims type")
}
