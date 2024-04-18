package migration

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/hasher"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	roles := []*entity.Role{
		{Id: entity.RoleUser, Name: "user"},
		{Id: entity.RoleAdmin, Name: "shop-admin"},
		{Id: entity.RoleSuperAdmin, Name: "super-admin"},
	}

	users := []*entity.User{
		{Email: "alice@example.com", Password: hashPassword("Alice12345"), RoleId: entity.RoleUser, IsVerified: true},
		{Email: "bob@example.com", Password: hashPassword("Bob12345"), RoleId: entity.RoleUser, IsVerified: true},
		{Email: "charlie@example.com", Password: hashPassword("Charlie12345"), RoleId: entity.RoleAdmin, IsVerified: true},
		{Email: "doni@example.com", Password: hashPassword("Doni12345"), RoleId: entity.RoleAdmin, IsVerified: true},
		{Email: "david@example.com", Password: hashPassword("David12345"), RoleId: entity.RoleSuperAdmin, IsVerified: true},
		{Email: "daniel@example.com", Password: hashPassword("Daniel12345"), RoleId: entity.RoleUser, IsVerified: true},
	}

	profiles := []*entity.Profile{
		generateProfile(1, "Alice"),
		generateProfile(2, "Bob"),
		generateProfile(3, "Charlie"),
		generateProfile(4, "David"),
	}
	carts := []entity.Cart{
		{UserId: 1},
		{UserId: 2},
		{UserId: 3},
		{UserId: 4},
	}

	productCategories := ImportProductCategories()

	products := ImportProduct()

	orderStatuses := []*entity.OrderStatus{
		{Name: "waiting for payment"},
		{Name: "waiting for payment confirmation"},
		{Name: "processed"},
		{Name: "sent"},
		{Name: "order confirmed"},
		{Name: "canceled"},
	}

	shops := ImportShop()

	provinces := getProvinces()

	cities := getCities()

	addresses := []*entity.Address{
		{
			Name:       "Alice",
			StreetName: "Jalan Mega Kuningan Barat",
			PostalCode: "12950",
			Phone:      "08772348585",
			Detail:     "",
			Location: &valueobject.Coordinate{
				Latitude:  generateDecimalFromString("-6.230835326032342"),
				Longitude: generateDecimalFromString("106.82413596846786"),
			},
			IsDefault:  false,
			ProfileId:  1,
			ProvinceId: 6,
			CityId:     42,
		},
		{
			Name:       "Alice",
			StreetName: "Jalan Cihampelas No 160",
			PostalCode: "12950",
			Phone:      "08772348585",
			Detail:     "",
			Location: &valueobject.Coordinate{
				Latitude:  generateDecimalFromString("-6.894393363416537"),
				Longitude: generateDecimalFromString("107.60773740044303"),
			},
			IsDefault:  false,
			ProfileId:  1,
			ProvinceId: 9,
			CityId:     64,
		},
		{
			Name:       "Alice",
			StreetName: "Jalan Dharmahusada 144",
			PostalCode: "60285",
			Phone:      "08772348585",
			Detail:     "",
			Location: &valueobject.Coordinate{
				Latitude:  generateDecimalFromString("-7.266614744067581"),
				Longitude: generateDecimalFromString("112.77033130376155"),
			},
			IsDefault:  true,
			ProfileId:  1,
			ProvinceId: 11,
			CityId:     161,
		},
	}

	orders := []entity.ProductOrder{
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 1,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  2,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 2,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  2,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 3,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(309000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  2,
			TotalPayment:  decimal.NewFromInt(209000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 4,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  1,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 5,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  2,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 6,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  1,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.December, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 1,
			ProfileId:     2,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  1,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "Mandiri",
			CreatedAt:     time.Date(2024, time.December, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
	}

	orderItems := []entity.OrderItem{
		{
			OrderId:       1,
			ShopProductId: 1,
			Quantity:      12,
			SubTotal:      decimal.NewFromInt(1196400),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderId:       1,
			ShopProductId: 2,
			Quantity:      4,
			SubTotal:      decimal.NewFromInt(398800),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderId:       2,
			ShopProductId: 4,
			Quantity:      4,
			SubTotal:      decimal.NewFromInt(132400),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderId:       2,
			ShopProductId: 5,
			Quantity:      4,
			SubTotal:      decimal.NewFromInt(372000),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderId:       3,
			ShopProductId: 3,
			Quantity:      8,
			SubTotal:      decimal.NewFromInt(797600),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderId:       3,
			ShopProductId: 3,
			Quantity:      4,
			SubTotal:      decimal.NewFromInt(102400),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderId:       4,
			ShopProductId: 290,
			Quantity:      4,
			SubTotal:      decimal.NewFromInt(102400),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderId:       5,
			ShopProductId: 6,
			Quantity:      8,
			SubTotal:      decimal.NewFromInt(204800),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderId:       5,
			ShopProductId: 4,
			Quantity:      8,
			SubTotal:      decimal.NewFromInt(548000),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderId:       6,
			ShopProductId: 6,
			Quantity:      4,
			SubTotal:      decimal.NewFromInt(398800),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderId:       7,
			ShopProductId: 10,
			Quantity:      4,
			SubTotal:      decimal.NewFromInt(398800),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	pharmacyProduct := generateShopProduct(301, 0)
	pharmacyProduct2 := generateShopProduct(601, 300)
	pharmacyProduct3 := generateShopProduct(908, 600)

	stockRecords := []*entity.StockRecord{
		{ShopProductId: 1, Quantity: 11, IsReduction: false, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{ShopProductId: 1, Quantity: 5, IsReduction: true, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{ShopProductId: 41, Quantity: 15, IsReduction: false, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{ShopProductId: 41, Quantity: 11, IsReduction: true, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{ShopProductId: 31, Quantity: 10, IsReduction: false, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{ShopProductId: 31, Quantity: 2, IsReduction: true, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{ShopProductId: 21, Quantity: 15, IsReduction: false, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{ShopProductId: 21, Quantity: 10, IsReduction: true, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
	}

	users2, adminContact2 := generateAdminShop(len(users), len(users)+17, "admin", "Admin12345")

	db.Create(roles)
	db.Create(users)
	db.Create(users2)
	db.Create(adminContact2)
	db.Create(profiles)
	db.Create(carts)
	db.Create(provinces)
	db.Create(cities)
	db.Create(addresses)
	db.Create(productCategories)
	db.Create(shops)
	db.Create(products)
	db.Create(orderStatuses)
	db.Create(orders)
	db.Create(pharmacyProduct)
	db.Create(pharmacyProduct2)
	db.Create(pharmacyProduct3)
	db.Create(orderItems)
	db.Create(stockRecords)

}

func hashPassword(text string) string {
	h := hasher.NewHasher()
	hashedText, err := h.Hash(text)
	if err != nil {
		return ""
	}
	return string(hashedText)
}
