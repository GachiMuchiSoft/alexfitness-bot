package models

import (
	"gorm.io/gorm"
)

type CoachStudent struct {
	gorm.Model
	CoachID   uint
	StudentID uint
	Student   User `gorm:"foreignKey:StudentID"`
}

func (u CoachStudent) TableName() string {
	return "coach_students"
}

func NewCoachStudent(coachID, studentID uint) *CoachStudent {
	return &CoachStudent{
		CoachID:   coachID,
		StudentID: studentID,
	}
}
