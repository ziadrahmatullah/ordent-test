package migration

import (
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	users := []*entity.User{
		{Email: "alice@example.com", Username: "alice", Password: hashPassword("Alice12345"), Role: entity.RoleUser, IsVerified: true},
		{Email: "bob@example.com", Username: "bob", Password: hashPassword("Bob12345"), Role: entity.RoleUser, IsVerified: true},
		{Email: "charlie@example.com", Username: "charlie", Password: hashPassword("Charlie12345"), Role: entity.RoleUser, IsVerified: true},
		{Email: "doni@example.com", Username: "doni", Password: hashPassword("Doni12345"), Role: entity.RoleUser, IsVerified: true},
		{Email: "david@example.com", Username: "david", Password: hashPassword("David12345"), Role: entity.RoleUser, IsVerified: true},
		{Email: "daniel@example.com", Username: "daniel", Password: hashPassword("Daniel12345"), Role: entity.RoleAdmin, IsVerified: true},
	}

	db.Create(users)
}

func hashPassword(text string) string {
	h := hasher.NewHasher()
	hashedText, err := h.Hash(text)
	if err != nil {
		return ""
	}
	return string(hashedText)
}
