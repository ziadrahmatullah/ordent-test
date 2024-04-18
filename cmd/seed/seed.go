package main

import (
	"os"

)

func main() {
	logger.SetLogrusLogger()

	_ = os.Setenv("APP_ENV", "debug")

	db, err := repository.GetConnection()
	if err != nil {
		logger.Log.Error(err)
	}

	migration.Seed(db)
}
