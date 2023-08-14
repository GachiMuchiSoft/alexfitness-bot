package repository

import (
	"fmt"
	"gachiSoft/alexfitness-bot/internal/app/implementation/models"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{
		db: db,
	}
}

func (u User) GetByID(id int) (user *models.User, err error) {
	err = u.db.First(&user, id).Error
	if err != nil {
		return nil, fmt.Errorf("error during getting user by id: %w", err)
	}

	return user, nil
}
