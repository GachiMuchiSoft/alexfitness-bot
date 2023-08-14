package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"gachiSoft/alexfitness-bot/internal/app/implementation/bot"
	"gachiSoft/alexfitness-bot/internal/app/implementation/bot/messages"
	"gachiSoft/alexfitness-bot/internal/app/implementation/controllers"
	"gachiSoft/alexfitness-bot/internal/app/implementation/models"
	"gachiSoft/alexfitness-bot/internal/app/implementation/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func NewRouter(
	commands *controllers.Commands,
	coach *controllers.Coach,
	student *controllers.Student,
	db *gorm.DB,
	bot *tgbotapi.BotAPI,
	authenticator Authenticator,
	cache *repository.Cache,
) *Router {
	return &Router{
		controllers: struct {
			commands *controllers.Commands
			coach    *controllers.Coach
			student  *controllers.Student
		}{
			commands: commands,
			coach:    coach,
			student:  student,
		},
		db:            db,
		bot:           bot,
		authenticator: authenticator,
		cache:         cache,
	}
}

type Router struct {
	controllers struct {
		commands *controllers.Commands
		coach    *controllers.Coach
		student  *controllers.Student
	}
	db            *gorm.DB
	bot           *tgbotapi.BotAPI
	authenticator Authenticator
	cache         *repository.Cache
}

func (r Router) Route(ctx context.Context, update tgbotapi.Update) error {
	user, err := r.authenticator.authenticate(update)
	if err != nil {
		return fmt.Errorf("error during user authentication: %w", err)
	}

	if user.ID == 0 {
		newStudentUser := models.NewUserStudent(
			fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName),
			uint(update.FromChat().ID),
		)

		if result := r.db.Create(newStudentUser); result.Error != nil {
			return fmt.Errorf("error during registering new user: %w", err)
		}
	}

	switch {
	case update.CallbackQuery != nil:
		if err := r.byCallback(ctx, &update, user); err != nil {
			return fmt.Errorf("error during routing callback update: %w", err)
		}
	case update.Message.IsCommand():
		return r.byCommand(ctx, update, user)
	}

	return nil
}

func (r Router) byCallback(ctx context.Context, update *tgbotapi.Update, user *models.User) error {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := r.bot.Request(callback); err != nil {
		panic(err)
	}

	var callbackDataPointer bot.CallbackData
	_ = json.Unmarshal([]byte(update.CallbackQuery.Data), &callbackDataPointer)

	callbackData, err := r.cache.GetButtonData(ctx, callbackDataPointer.ID)
	if err != nil {
		return fmt.Errorf("error during loading button data from cache: %w", err)
	} else if callbackData == "" {
		_, err = r.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Пожалуйста не пользуйтесь старыми кнопками"))
		return err
	}

	update.CallbackQuery.Data = callbackData

	var data messages.RoutingData
	if err := json.Unmarshal([]byte(callbackData), &data); err != nil {
		return fmt.Errorf("invalid keyboard button data: %w", err)
	}

	switch data.Routing {
	case messages.RoutingKeyCoachCreateTrainingSelectingStudent:
		return r.controllers.coach.CreateTrainingSelectStudent(ctx, update, user)
	case messages.RoutingKeyCoachCreateTrainingSelectingTime:
		return r.controllers.coach.CreateTrainingSelectTime(ctx, update, user)
	case messages.RoutingKeyCoachCreateTrainingSelectingTrainingType:
		return r.controllers.coach.CreateTrainingSelectTrainingType(ctx, update, user)
	case messages.RoutingKeyCoachCreateTraining:
		return r.controllers.coach.CreateTraining(ctx, update, user)

	case messages.RoutingKeyCoachAddTrainingsStudentSelecting:
		return r.controllers.coach.SelectStudentToAddTraining(ctx, update, user)
	case messages.RoutingKeyCoachAddTrainingsStudent:
		return r.controllers.coach.AddTrainingsToStudent(ctx, update, user)

	case messages.RoutingKeyStudentShowTrainingBalance:
		return r.controllers.coach.ShowTrainingBalance(ctx, update, user)

	case messages.RoutingKeyStudentRemoveTraining:
		return r.controllers.coach.RemoveTraining(ctx, update, user)
	}

	return nil
}

func (r Router) byCommand(ctx context.Context, update tgbotapi.Update, user *models.User) error {
	switch update.Message.Command() {
	case controllers.CommandStart:
		if err := r.controllers.commands.Start(ctx, update, user); err != nil {
			return fmt.Errorf("error during commands controller start:%w", err)
		}
	default:
		return fmt.Errorf("unknown chat bot command: %s", update.Message.Text)
	}

	return nil
}
