package routing

import (
	"fmt"
	"gachiSoft/alexfitness-bot/internal/app/implementation/environment"
	"gachiSoft/alexfitness-bot/internal/app/implementation/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type Authenticator interface {
	authenticate(update tgbotapi.Update) (*models.User, error)
}

func NewProductionAuthenticator(db *gorm.DB) *ProductionAuthenticator {
	return &ProductionAuthenticator{db: db}
}

type ProductionAuthenticator struct {
	db *gorm.DB
}

func (a ProductionAuthenticator) authenticate(update tgbotapi.Update) (*models.User, error) {
	userID := update.FromChat().ID

	var user models.User
	if result := a.db.Where("ChatID = ?", userID).Find(&user); result.Error != nil {
		return nil, fmt.Errorf("error during searching user by id %d: %w", update.Message.Chat.ID, result.Error)
	}

	return &user, nil
}

func NewTestAuthenticator(db *gorm.DB) *TestAuthenticator {
	return &TestAuthenticator{db: db}
}

type TestAuthenticator struct {
	db *gorm.DB
}

func (t TestAuthenticator) authenticate(_ tgbotapi.Update) (*models.User, error) {
	var user models.User

	testUserID, err := environment.TestUserID()
	if err != nil {
		return nil, err
	}

	if result := t.db.Find(&user, testUserID); result.Error != nil {
		return nil, fmt.Errorf("invalid test user id : %w", result.Error)
	}

	return &user, nil
}
