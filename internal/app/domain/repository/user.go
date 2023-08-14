package repository

import "gachiSoft/alexfitness-bot/internal/app/domain/models/user"

type User interface {
	GetByID(id uint) user.User
}
