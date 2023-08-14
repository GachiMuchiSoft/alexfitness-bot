package main

import (
	"context"
	"fmt"
	"gachiSoft/alexfitness-bot/internal/app/domain"
	drepository "gachiSoft/alexfitness-bot/internal/app/domain/repository"
	"gachiSoft/alexfitness-bot/internal/app/domain/service"
	"gachiSoft/alexfitness-bot/internal/app/implementation/bot"
	"gachiSoft/alexfitness-bot/internal/app/implementation/controllers"
	"gachiSoft/alexfitness-bot/internal/app/implementation/environment"
	"gachiSoft/alexfitness-bot/internal/app/implementation/models"
	irepository "gachiSoft/alexfitness-bot/internal/app/implementation/repository"
	"gachiSoft/alexfitness-bot/internal/app/implementation/routing"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/dig"
)

func main() {
	di, err := bootstrap()
	ctx := context.Background()
	if err != nil {
		panic(fmt.Errorf("error during creating container %w", err))
	}

	if err = di.Invoke(func(bot *tgbotapi.BotAPI, router *routing.Router) error {
		if err != nil {
			return fmt.Errorf("error during authorizing bot: %w", err)
		}
		bot.Debug = true

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		for update := range bot.GetUpdatesChan(u) {
			if err := router.Route(ctx, update); err != nil {
				return fmt.Errorf("error during routing new update: %w", err)
			}
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("error during invoking router: %w", err))
	}
}

func bootstrap() (*dig.Container, error) {
	di := dig.New()

	if err := di.Provide(func() (*tgbotapi.BotAPI, error) {
		return tgbotapi.NewBotAPI(environment.BotToken())
	}); err != nil {
		return nil, fmt.Errorf("error during providing `controllers.Commands` di component: %w", err)
	}

	if err := di.Provide(func() (*controllers.Commands, error) {
		return controllers.NewCommands(di), nil
	}); err != nil {
		return nil, fmt.Errorf("error during providing `controllers.Commands` di component: %w", err)
	}

	if err := di.Provide(func() (*controllers.Coach, error) {
		return controllers.NewCoach(di), nil
	}); err != nil {
		return nil, fmt.Errorf("error during providing `controllers.Coach` di component: %w", err)
	}

	if err := di.Provide(func() (*controllers.Student, error) {
		return controllers.NewStudent(di), nil
	}); err != nil {
		return nil, fmt.Errorf("error during providing `controllers.Student` di component: %w", err)
	}

	if err := di.Provide(routing.NewRouter); err != nil {
		return nil, fmt.Errorf("error during providing `routing.Router` di component: %w", err)
	}

	if err := di.Provide(bot.New); err != nil {
		return nil, fmt.Errorf("error during providing `bot.Bot` di component: %w", err)
	}

	if err := di.Provide(func(db *gorm.DB) (routing.Authenticator, error) {
		if environment.IsProduction() {
			return routing.NewProductionAuthenticator(db), nil
		}

		return routing.NewTestAuthenticator(db), nil
	}); err != nil {
		return nil, fmt.Errorf("error during providing `routing.Authenticator` di component: %w", err)
	}

	if err := di.Provide(func() (*gorm.DB, error) {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			environment.PostgresHost(),
			environment.PostgresUser(),
			environment.PostgresPassword(),
			environment.PostgresDB(),
			environment.PostgresPort(),
		)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("error during connecting to database: %w", err)
		}

		err = db.AutoMigrate(&models.User{}, &models.CoachStudent{})
		if err != nil {
			return nil, fmt.Errorf("error during database migrations: %w", err)
		}

		return db, err
	}); err != nil {
		return nil, fmt.Errorf("error during providing `gorm.DB` di component: %w", err)
	}

	if err := di.Provide(func() *redis.Client {
		return redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", environment.RedisHost(), environment.RedisPort()),
			Password: environment.RedisPassword(),
			DB:       0,
		})
	}); err != nil {
		return nil, fmt.Errorf("error during providing `redis.Client` di component: %w", err)
	}

	if err := di.Provide(irepository.NewCache); err != nil {
		return nil, fmt.Errorf("error during providing `repository.Cache` di component: %w", err)
	}

	if err := di.Provide(irepository.NewUser); err != nil {
		return nil, fmt.Errorf("error during providing `repository.User` di component: %w", err)
	}

	if err := di.Provide(func(
		studentRepository drepository.Student,
		trainingRepository drepository.Training,
		notifier domain.Notifier,
	) service.Coach {
		return service.NewCoach(studentRepository, trainingRepository, notifier)
	}); err != nil {
		return nil, fmt.Errorf("error during providing `drepository.Student` di component: %w", err)
	}

	if err := di.Provide(func(db *gorm.DB) drepository.Training {
		return irepository.NewTraining(db)
	}); err != nil {
		return nil, fmt.Errorf("error during providing `drepository.Training` di component: %w", err)
	}

	if err := di.Provide(func(redis *redis.Client) drepository.Student {
		return irepository.NewStudent(redis)
	}); err != nil {
		return nil, fmt.Errorf("error during providing `drepository.Student` di component: %w", err)
	}

	if err := di.Provide(irepository.NewStudent); err != nil {
		return nil, fmt.Errorf("error during providing `irepository.Student` di component: %w", err)
	}

	if err := di.Provide(func(bot *bot.Bot) domain.Notifier {
		return bot
	}); err != nil {
		return nil, fmt.Errorf("error during providing `domain.Notifier` di component: %w", err)
	}

	return di, nil
}
