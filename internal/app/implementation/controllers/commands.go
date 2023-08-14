package controllers

import (
	"context"
	"fmt"
	"gachiSoft/alexfitness-bot/internal/app/implementation/bot"
	"gachiSoft/alexfitness-bot/internal/app/implementation/bot/messages"
	iuser "gachiSoft/alexfitness-bot/internal/app/implementation/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/dig"
)

func NewCommands(container *dig.Container) *Commands {
	return &Commands{
		container: container,
	}
}

type Commands struct {
	container *dig.Container
}

const CommandStart = "start"

func (controller Commands) Start(ctx context.Context, _ tgbotapi.Update, user *iuser.User) error {
	if err := controller.container.Invoke(func(bot *bot.Bot) error {
		var message *tgbotapi.MessageConfig
		switch {
		case user.IsCoach():
			message = messages.StartCoach(int64(user.ChatID))
		case user.IsStudent():
			message = messages.StartStudent(int64(user.ChatID))
		}

		err := bot.Send(ctx, message)

		return err
	}); err != nil {
		return fmt.Errorf("error during invoking di in start controller: %w", err)
	}

	return nil
}
