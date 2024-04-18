package migration

import (
	"github.com/ziadrahmatullah/ordent-test/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	u := &entity.User{}
	p := &entity.Profile{}
	fp := &entity.ForgotPasswordToken{}
	c := &entity.Cart{}

	_ = db.Migrator().DropTable(u, p, fp, c)
	_ = db.AutoMigrate(u, p, fp, c)
}
