package messages

import (
	"gachiSoft/alexfitness-bot/internal/app/domain/models"
	"gachiSoft/alexfitness-bot/internal/app/domain/models/user"
)

type CoachTrainingCreated struct {
	Training               models.Training
	TrainingRemainingCount int
}

type CoachStudentHasRunOutOfTrainings struct {
}

type CoachStudentHasNotEnoughAvailableTrainings struct {
}

type CoachAddedAvailableTrainingsToStudent struct {
	Student user.Student
	Count   int
}
