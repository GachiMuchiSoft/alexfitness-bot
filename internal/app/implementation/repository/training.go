package repository

import (
	dmodels "gachiSoft/alexfitness-bot/internal/app/domain/models"
	duser "gachiSoft/alexfitness-bot/internal/app/domain/models/user"
	imodels "gachiSoft/alexfitness-bot/internal/app/implementation/models"
	"gorm.io/gorm"
	"time"
)

type Training struct {
	db *gorm.DB
}

func NewTraining(db *gorm.DB) *Training {
	return &Training{
		db: db,
	}
}

func (t *Training) Create(student duser.Student, trainingType dmodels.TrainingType, time time.Time) (dmodels.Training, error) {
	training := imodels.NewTraining(student.GetID(), trainingType.GetID(), time)
	if err := t.db.Create(training).Error; err != nil {
		return nil, err
	}

	return training, t.db.Joins("Type").Joins("Student").Find(&training).Error
}
