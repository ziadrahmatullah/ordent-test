package migration

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/ordent-test/entity"
)

func generateProfile(id uint, name string) *entity.Profile {
	return &entity.Profile{
		UserId:    id,
		Name:      name,
		Birthdate: time.Date(2000, 01, 01, 0, 0, 0, 0, time.Local),
	}
}

func generateDecimalFromString(decimalString string) decimal.Decimal {
	zero := decimal.NewFromInt(0)
	decimal, err := decimal.NewFromString(decimalString)
	if err != nil {
		return zero
	}
	return decimal
}

func generateUser(email, password string, roleId entity.RoleId, isVerified bool) *entity.User {
	return &entity.User{Email: email, Password: hashPassword(password), RoleId: roleId, IsVerified: isVerified}
}

func generateAdminContact(userId uint, name, phone string) *entity.AdminContact {
	return &entity.AdminContact{UserId: userId, Name: name, Phone: phone}
}

func generateAdminShop(start, max int, baseName, password string) ([]*entity.User, []*entity.AdminContact) {
	var users []*entity.User
	var adminContacts []*entity.AdminContact
	for i := start; i+1 < max+2; i++ {
		iString := strconv.Itoa(i - start + 1)
		temp := generateUser(baseName+iString+"@gmail.com", password, entity.RoleAdmin, true)
		users = append(users, temp)
		tempAdmin := generateAdminContact(uint(i+1), baseName+iString, "089654749370")
		adminContacts = append(adminContacts, tempAdmin)
	}
	return users, adminContacts
}

func generateShopProduct(max, min int) []*entity.ShopProduct {
	var pps []*entity.ShopProduct
	var productSlice []*entity.Product
	data, _ := os.ReadFile("./migration/data/product-details.json")

	_ = json.Unmarshal(data, &productSlice)
	for i := 1; i <= 18; i++ {
		for _, product := range productSlice[min:max] {
			pps = append(pps, &entity.ShopProduct{
				ProductId: product.Id,
				ShopId:    uint(i),
				Price:     decimal.NewFromInt(int64(randomNumber(int(product.Price.IntPart())-2000, int(product.Price.IntPart())+2000))).Abs(),
				Stock:     randomNumber(10, 30),
			})
		}
	}
	return pps
}

func randomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func RandomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 0
}
