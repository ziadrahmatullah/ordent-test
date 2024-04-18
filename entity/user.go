package entity

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	Id         uint       `gorm:"primaryKey;autoIncrement"`
	Email      string     `gorm:"unique;not null"`
	Username   string     `gorm:"unique;not null"`
	Password   string     `gorm:"not null"`
	Role       Role       `gorm:"not null"`
	IsVerified bool       `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
