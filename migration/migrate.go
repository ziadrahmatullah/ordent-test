package migration

import (
	"github.com/ziadrahmatullah/ordent-test/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	u := &entity.User{}
	pc := &entity.ProductCategory{}
	p := &entity.Product{}
	r := &entity.Role{}
	profile := &entity.Profile{}
	ftp := &entity.ForgotPasswordToken{}
	pharmacy := &entity.Shop{}
	pharmacyProduct := &entity.ShopProduct{}
	pr := &entity.Province{}
	ct := &entity.City{}
	ors := &entity.OrderStatus{}
	spm := &entity.ShippingMethod{}
	a := &entity.Address{}
	c := &entity.Cart{}
	ci := &entity.CartItem{}
	ts := &entity.StockRecord{}
	po := &entity.ProductOrder{}
	oi := &entity.OrderItem{}
	ac := &entity.AdminContact{}

	_ = db.Migrator().DropTable(u, pc, p, r, profile, ftp, pharmacy, pharmacyProduct, pr, ct, ors, spm, a, c, ci, ts, ac, po, oi)

	_ = db.AutoMigrate(u, pc, p, r, profile, ftp, pharmacy, pharmacyProduct, pr, ct, ors, spm, a, c, ci, ts, ac, po, oi)
}
