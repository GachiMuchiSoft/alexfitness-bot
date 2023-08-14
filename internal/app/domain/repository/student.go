package repository

import (
	"gachiSoft/alexfitness-bot/internal/app/domain/models/user"
)

type Student interface {
	GetTrainings(student user.Student) (int, error)
	AddAvailableTrainingsQuantity(student user.Student, quantity int64) (int, error)
	RemoveAvailableTrainingsQuantity(student user.Student, quantity int64) error
}
