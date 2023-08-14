package repository

import (
	"gachiSoft/alexfitness-bot/internal/app/domain/models"
	"gachiSoft/alexfitness-bot/internal/app/domain/models/user"
	"time"
)

type Training interface {
	Create(student user.Student, trainingType models.TrainingType, time time.Time) (models.Training, error)
}
