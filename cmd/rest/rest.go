package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/ziadrahmatullah/ordent-test/appjwt"
	"github.com/ziadrahmatullah/ordent-test/appvalidator"
	"github.com/ziadrahmatullah/ordent-test/handler"
	"github.com/ziadrahmatullah/ordent-test/hasher"
	"github.com/ziadrahmatullah/ordent-test/logger"
	"github.com/ziadrahmatullah/ordent-test/repository"
	"github.com/ziadrahmatullah/ordent-test/router"
	"github.com/ziadrahmatullah/ordent-test/server"
	"github.com/ziadrahmatullah/ordent-test/transactor"
	"github.com/ziadrahmatullah/ordent-test/usecase"
)

func main() {
	logger.SetLogrusLogger()

	db, err := repository.GetConnection()
	if err != nil {
		logger.Log.Error(err)
	}

	client := resty.New()
	client.SetHeader("Content-Type", "application/json")

	manager := transactor.NewManager(db)
	hash := hasher.NewHasher()
	jwt := appjwt.NewJwt()
	appvalidator.RegisterCustomValidator()

	userR := repository.NewUserRepository(db)
	profileR := repository.NewProfileRepository(db)
	cartR := repository.NewCartRepository(db)
	forgotPassR := repository.NewForgotPasswordRepository(db)
	addressR := repository.NewAddressRepository(db)

	userU := usecase.NewUserUsecase(userR, profileR, hash)
	authU := usecase.NewAuthUsecase(manager, userR, profileR, forgotPassR, cartR, hash, jwt)
	addressU := usecase.NewAddressUsecase(addressR, manager)

	userH := handler.NewUserHandler(userU)
	authH := handler.NewAuthHandler(authU)
	addressH := handler.NewAddressHandler(addressU)

	handlers := router.Handlers{
		User:    userH,
		Auth:    authH,
		Address: addressH,
	}

	r := router.New(handlers)

	s := server.New(r)

	server.StartWithGracefulShutdown(s)
}
