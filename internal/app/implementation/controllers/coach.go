package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"gachiSoft/alexfitness-bot/internal/app/domain/service"
	"gachiSoft/alexfitness-bot/internal/app/implementation/bot"
	"gachiSoft/alexfitness-bot/internal/app/implementation/bot/messages"
	"gachiSoft/alexfitness-bot/internal/app/implementation/models"
	"gachiSoft/alexfitness-bot/internal/app/implementation/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"time"
)

func NewCoach(container *dig.Container) *Coach {
	return &Coach{
		container: container,
	}
}

type Coach struct {
	container *dig.Container
}

func (controller Coach) SelectStudentToAddTraining(ctx context.Context, _ *tgbotapi.Update, user *models.User) error {
	if err := controller.container.Invoke(func(bot *bot.Bot) error {
		students, err := controller.getCoachStudents(user)
		if err != nil {
			return err
		}

		err = bot.Send(ctx, messages.SelectStudent(
			int64(user.ChatID),
			students,
			messages.RoutingKeyCoachAddTrainingsStudent,
		))

		return err
	}); err != nil {
		return fmt.Errorf("error during invoking di in coach controller: %w", err)
	}

	return nil
}

func (controller Coach) AddTrainingsToStudent(_ context.Context, update *tgbotapi.Update, user *models.User) error {
	if err := controller.container.Invoke(func(studentsRepository *repository.User, coachService service.Coach) error {
		var data struct {
			SelectedStudentID uint `json:"student_id"`
		}

		if err := json.Unmarshal([]byte(update.CallbackQuery.Data), &data); err != nil {
			return fmt.Errorf("wrong button data for coach controller, action add trainings to student: %w", err)
		}

		student, err := studentsRepository.GetByID(int(data.SelectedStudentID))
		if err != nil {
			return err
		}

		if err := coachService.AddAvailableTrainingsQuantity(user, student); err != nil {
			return err
		}

		return err
	}); err != nil {
		return fmt.Errorf("error during invoking di in coach controller: %w", err)
	}

	return nil
}

func (controller Coach) CreateTrainingSelectStudent(ctx context.Context, _ *tgbotapi.Update, user *models.User) error {
	if err := controller.container.Invoke(func(bot *bot.Bot) error {
		students, err := controller.getCoachStudents(user)
		if err != nil {
			return err
		}

		err = bot.Send(ctx, messages.SelectStudent(
			int64(user.ChatID),
			students,
			messages.RoutingKeyCoachCreateTrainingSelectingTrainingType,
		))

		return err
	}); err != nil {
		return fmt.Errorf("error during invoking di in coach controller: %w", err)
	}

	return nil
}

func (controller Coach) CreateTrainingSelectTrainingType(ctx context.Context, update *tgbotapi.Update, user *models.User) error {
	if err := controller.container.Invoke(func(bot *bot.Bot, studentRepository *repository.Student) error {
		var data messages.SelectStudentData
		_ = json.Unmarshal([]byte(update.CallbackQuery.Data), &data)

		studentAvailableTrainingsCount, err := studentRepository.GetTrainings(&models.User{Model: gorm.Model{ID: data.StudentID}})
		if err != nil {
			return err
		}

		if studentAvailableTrainingsCount == 0 {
			return bot.Send(
				ctx,
				messages.CoachUnableToAddTrainingStudentHasNotEnoughAvailableTrainings(user.ChatID),
			)
		}

		trainingTypes, err := controller.getCoachTrainingTypes(user)
		if err != nil {
			return err
		}

		err = bot.Send(ctx, messages.SelectTrainingType(
			int64(user.ChatID),
			messages.RoutingKeyCoachCreateTrainingSelectingTime,
			trainingTypes,
			map[string]interface{}{
				"student_id": data.StudentID,
			},
		))

		return err
	}); err != nil {
		return fmt.Errorf("error during invoking di in coach controller: %w", err)
	}

	return nil
}

func (controller Coach) CreateTrainingSelectTime(ctx context.Context, update *tgbotapi.Update, user *models.User) error {
	if err := controller.container.Invoke(func(bot *bot.Bot) error {
		var data struct {
			StudentID uint `json:"student_id"`
			TypeID    uint `json:"type_id"`
		}

		_ = json.Unmarshal([]byte(update.CallbackQuery.Data), &data)

		err := bot.Send(ctx, messages.SelectTime(
			int64(user.ChatID),
			messages.RoutingKeyCoachCreateTraining,
			map[string]interface{}{
				"student_id": data.StudentID,
				"type_id":    data.TypeID,
			},
		))

		return err
	}); err != nil {
		return fmt.Errorf("error during invoking di in coach controller: %w", err)
	}

	return nil
}

func (controller Coach) CreateTraining(_ context.Context, update *tgbotapi.Update, user *models.User) error {
	if err := controller.container.Invoke(func(bot *bot.Bot, coachService service.Coach, userRepository *repository.User) error {
		var data struct {
			StudentID uint                  `json:"student_id"`
			TypeID    uint                  `json:"type_id"`
			Time      messages.HoursMinutes `json:"time"`
		}
		_ = json.Unmarshal([]byte(update.CallbackQuery.Data), &data)

		student, err := userRepository.GetByID(int(data.StudentID))
		if err != nil {
			return err
		}

		now := time.Now()
		return coachService.CreateTraining(
			user,
			student,
			&models.TrainingType{Model: gorm.Model{ID: data.TypeID}},
			time.Date(now.Year(), now.Month(), now.Day(), time.Time(data.Time).Hour(), time.Time(data.Time).Minute(), 0, 0, time.Local),
		)
	}); err != nil {
		return fmt.Errorf("error during invoking di in coach controller: %w", err)
	}

	return nil
}

func (controller Coach) getCoachStudents(user *models.User) (students []models.CoachStudent, err error) {
	if err = controller.container.Invoke(func(db *gorm.DB) error {
		if err = db.Where("coach_id = ?", user.ID).Joins("Student").Find(&students).Error; err != nil {
			return fmt.Errorf("error during searching coach students: %w", err)
		}

		return err
	}); err != nil {
		return nil, fmt.Errorf("error during invoking di in coach controller: %w", err)
	}

	return students, err
}

func (controller Coach) getCoachTrainingTypes(user *models.User) (types []models.TrainingType, err error) {
	if err = controller.container.Invoke(func(db *gorm.DB) error {
		if err = db.Where("coach_id = ?", user.ID).Find(&types).Error; err != nil {
			return fmt.Errorf("error during searching coach training types: %w", err)
		}

		return err
	}); err != nil {
		return nil, fmt.Errorf("error during invoking di in coach controller: %w", err)
	}

	return types, err
}
