package models

import (
	"gorm.io/gorm"
)

type Role int

const (
	coach   Role = iota
	student Role = iota
)

func NewUserStudent(name string, chatID uint) *User {
	return &User{
		Role:   int(student),
		Name:   name,
		ChatID: chatID,
	}
}

func NewUserCoach(name string, chatID uint) *User {
	return &User{
		Role:   int(coach),
		Name:   name,
		ChatID: chatID,
	}
}

type User struct {
	gorm.Model
	Role   int
	Name   string
	ChatID uint
}

func (u *User) GetName() string {
	return u.Name
}

func (u User) TableName() string {
	return "users"
}

func (u User) GetID() uint {
	return u.ID
}

func (u User) IsStudent() bool {
	return u.Role == int(student)
}

func (u User) IsCoach() bool {
	return u.Role == int(coach)
}
