package migration

import (
	"github.com/ziadrahmatullah/ordent-test/entity"
	"github.com/ziadrahmatullah/ordent-test/hasher"
	"github.com/ziadrahmatullah/ordent-test/valueobject"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
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

	db.Create(users)
	db.Create(profiles)
	db.Create(carts)
	db.Create(addresses)

}

func hashPassword(text string) string {
	h := hasher.NewHasher()
	hashedText, err := h.Hash(text)
	if err != nil {
		return ""
	}
	return string(hashedText)
}
