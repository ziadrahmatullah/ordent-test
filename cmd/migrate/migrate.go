package main

import (
	"os"

	"github.com/ziadrahmatullah/ordent-test/logger"
	"github.com/ziadrahmatullah/ordent-test/migration"
	"github.com/ziadrahmatullah/ordent-test/repository"
)

func main() {
	logger.SetLogrusLogger()

	_ = os.Setenv("APP_ENV", "debug")

	db, err := repository.GetConnection()
	if err != nil {
		logger.Log.Error(err)
	}

	migration.Migrate(db)
}
