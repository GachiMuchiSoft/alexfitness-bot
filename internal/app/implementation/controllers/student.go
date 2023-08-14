package controllers

import (
	"context"
	"fmt"
	"gachiSoft/alexfitness-bot/internal/app/implementation/bot"
	"gachiSoft/alexfitness-bot/internal/app/implementation/bot/messages"
	"gachiSoft/alexfitness-bot/internal/app/implementation/models"
	"gachiSoft/alexfitness-bot/internal/app/implementation/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func NewStudent(container *dig.Container) *Student {
	return &Student{
		container: container,
	}
}

type Student struct {
	container *dig.Container
}

func (controller Coach) ShowTrainingBalance(ctx context.Context, _ *tgbotapi.Update, user *models.User) error {
	if err := controller.container.Invoke(func(bot *bot.Bot, db *gorm.DB, cache *repository.Cache) error {
		trainingsCount, err := cache.GetStudentTrainings(ctx, user.ID)
		if err != nil {
			return err
		}

		err = bot.Send(ctx, messages.StudentShowTrainingBalance(user.ChatID, trainingsCount))

		return err
	}); err != nil {
		return fmt.Errorf("error during invoking di in student controller: %w", err)
	}

	return nil
}

func (controller Coach) RemoveTraining(ctx context.Context, _ *tgbotapi.Update, user *models.User) error {
	if err := controller.container.Invoke(func(bot *bot.Bot, db *gorm.DB, cache *repository.Cache) error {
		trainingsCount, err := cache.GetStudentTrainings(ctx, user.ID)
		if err != nil {
			return err
		}

		if trainingsCount == 0 {
			err := bot.Send(ctx, messages.StudentHasRunOutOfTrainings(user.ChatID))
			return err
		}

		err = cache.DecrStudentTrainings(ctx, user.ID, 1)
		if err != nil {
			return err
		}

		err = bot.Send(ctx, messages.StudentRemoveTrainingSuccess(user.ChatID))

		return err
	}); err != nil {
		return fmt.Errorf("error during invoking di in student controller: %w", err)
	}

	return nil
}
