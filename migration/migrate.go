package migration

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	u := &entity.User{}

	_ = db.Migrator().DropTable(u)
	_ = db.AutoMigrate(u)
}
