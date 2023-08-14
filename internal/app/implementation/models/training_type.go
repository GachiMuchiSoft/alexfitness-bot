package models

import (
	"gorm.io/gorm"
)

func NewTrainingType(name string, coachID uint) *TrainingType {
	return &TrainingType{
		Name:    name,
		CoachID: coachID,
	}
}

type TrainingType struct {
	Name    string
	CoachID uint
	Coach   User `gorm:"foreignKey:CoachID"`
	gorm.Model
}

func (t TrainingType) GetID() uint {
	return t.ID
}

func (t TrainingType) GetName() string {
	return t.Name
}
