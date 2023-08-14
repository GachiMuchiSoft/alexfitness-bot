package service

import (
	"gachiSoft/alexfitness-bot/internal/app/domain"
	"gachiSoft/alexfitness-bot/internal/app/domain/messages"
	"gachiSoft/alexfitness-bot/internal/app/domain/models"
	"gachiSoft/alexfitness-bot/internal/app/domain/models/user"
	"gachiSoft/alexfitness-bot/internal/app/domain/repository"
	"time"
)

type Coach interface {
	CreateTraining(coach user.Coach, student user.Student, trainingType models.TrainingType, time time.Time) error
	AddAvailableTrainingsQuantity(coach user.Coach, student user.Student) error
}

func NewCoach(
	studentRepository repository.Student,
	trainingRepository repository.Training,
	notifier domain.Notifier,
) *coach {
	return &coach{
		studentsRepository: studentRepository,
		trainingRepository: trainingRepository,
		notifier:           notifier,
	}
}

type coach struct {
	studentsRepository repository.Student
	trainingRepository repository.Training
	notifier           domain.Notifier
}

const trainingBlockQuantity = 10

func (c *coach) AddAvailableTrainingsQuantity(coach user.Coach, student user.Student) error {
	currentCount, err := c.studentsRepository.AddAvailableTrainingsQuantity(student, trainingBlockQuantity)
	if err != nil {
		return err
	}

	return c.notifier.Notify(coach, messages.CoachAddedAvailableTrainingsToStudent{
		Student: student,
		Count:   currentCount,
	})
}

func (c *coach) CreateTraining(coach user.Coach, student user.Student, trainingType models.TrainingType, time time.Time) error {
	availableTrainingsQuantity, err := c.studentsRepository.GetTrainings(student)
	if availableTrainingsQuantity == 0 {
		return c.notifier.Notify(coach, messages.CoachStudentHasNotEnoughAvailableTrainings{})
	}

	training, err := c.trainingRepository.Create(student, trainingType, time)
	if err != nil {
		return err
	}

	if err = c.studentsRepository.RemoveAvailableTrainingsQuantity(student, 1); err != nil {
		return err
	}

	availableTrainingsQuantity--
	if availableTrainingsQuantity == 0 {
		if err = c.notifier.Notify(coach, messages.CoachStudentHasRunOutOfTrainings{}); err != nil {
			return err
		}

		if err = c.notifier.Notify(student, messages.StudentLastTrainingFinished{}); err != nil {
			return err
		}
	}

	return c.notifier.Notify(coach, messages.CoachTrainingCreated{
		Training:               training,
		TrainingRemainingCount: availableTrainingsQuantity,
	})
}
