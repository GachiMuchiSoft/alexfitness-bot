package models

import (
	"gachiSoft/alexfitness-bot/internal/app/domain/models/user"
	"time"
)

type TrainingType interface {
	GetID() uint
	GetName() string
}

type Training interface {
	GetTime() time.Time
	GetStudent() user.Student
	GetType() TrainingType
	GetID() uint
}
