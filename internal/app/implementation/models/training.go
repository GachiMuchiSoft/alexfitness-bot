package models

import (
	"gachiSoft/alexfitness-bot/internal/app/domain/models"
	"gachiSoft/alexfitness-bot/internal/app/domain/models/user"
	"gorm.io/gorm"
	"time"
)

func NewTraining(studentID uint, typeID uint, time time.Time) *Training {
	return &Training{
		Time:      time,
		StudentID: studentID,
		TypeID:    typeID,
	}
}

type Training struct {
	Time      time.Time
	Student   User
	StudentID uint
	Type      TrainingType
	TypeID    uint
	gorm.Model
}

func (t Training) GetID() uint {
	return t.Model.ID
}

func (t Training) GetTime() time.Time {
	return t.Time
}

func (t Training) GetStudent() user.Student {
	return &t.Student
}

func (t Training) GetType() models.TrainingType {
	return t.Type
}
