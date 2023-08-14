package messages

import (
	"encoding/json"
	"fmt"
	dmodels "gachiSoft/alexfitness-bot/internal/app/domain/models"
	"gachiSoft/alexfitness-bot/internal/app/implementation/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/maps"
)

const (
	RoutingKeyCoachAddTrainingsStudentSelecting = "1.1"
	RoutingKeyCoachAddTrainingsStudent          = "1.2"
)

const (
	RoutingKeyCoachCreateTrainingSelectingStudent      = "2.1"
	RoutingKeyCoachCreateTrainingSelectingTrainingType = "2.2"
	RoutingKeyCoachCreateTrainingSelectingTime         = "2.3"
	RoutingKeyCoachCreateTraining                      = "2.4"
)

func StartCoach(chatId int64) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, "Вы зарегистрированы в системе как тренер, что будем делать?")
	addTrainingsBlockData, _ := json.Marshal(&RoutingData{Routing: RoutingKeyCoachAddTrainingsStudentSelecting})
	CreateTrainingData, _ := json.Marshal(&RoutingData{Routing: RoutingKeyCoachCreateTrainingSelectingStudent})
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить блок тренировок", string(addTrainingsBlockData)),
			tgbotapi.NewInlineKeyboardButtonData("Назначить тренировку", string(CreateTrainingData)),
		),
	)

	return &msg
}

type SelectStudentData struct {
	Routing   string `json:"routing"`
	StudentID uint   `json:"student_id"`
}

func SelectStudent(chatID int64, students []models.CoachStudent, routing string) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, "Выберите ученика")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(func() (buttons []tgbotapi.InlineKeyboardButton) {
			for _, student := range students {
				keyboardData, _ := json.Marshal(&SelectStudentData{
					Routing:   routing,
					StudentID: student.StudentID,
				})

				buttons = append(
					buttons,
					tgbotapi.NewInlineKeyboardButtonData(student.Student.Name, string(keyboardData)),
				)
			}

			return buttons
		}()...),
	)

	return &msg
}

func SelectTrainingType(chatID int64, routing string, types []models.TrainingType, data map[string]interface{}) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, "Выберите тип тренировки")
	buttonData := func(ttype models.TrainingType) string {
		dataToMerge := data
		maps.Copy(dataToMerge, map[string]interface{}{
			"routing": routing,
			"type_id": ttype.ID,
		})

		marshalledData, _ := json.Marshal(dataToMerge)
		return string(marshalledData)
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(func() (buttons []tgbotapi.InlineKeyboardButton) {
			for _, ttype := range types {
				buttons = append(
					buttons,
					tgbotapi.NewInlineKeyboardButtonData(ttype.Name, buttonData(ttype)),
				)
			}

			return buttons
		}()...),
	)

	return &msg
}

func CoachSuccessfullyAddTrainingsBlock(chatId uint, balance int) *tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(int64(chatId), fmt.Sprintf("Ученику успешно добавлены тренировки (Баланс = %d)", balance))
	return &message
}

func CoachAddedTrainingToStudentAndStudentHasRunOutOfTrainings(chatId uint) *tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(int64(chatId), fmt.Sprintf("Это последняя тренировка в блоке, ученику было отправлено уведомление об оплате"))
	return &message
}

func CoachUnableToAddTrainingStudentHasNotEnoughAvailableTrainings(chatId uint) *tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(int64(chatId), fmt.Sprintf("Нельзя добавить тренировку, ученику нужно оплатить блок"))
	return &message
}

func CoachTrainingCreated(chatId uint, training dmodels.Training, trainingsRemainingCount int) *tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(
		int64(chatId),
		fmt.Sprintf("**Тренировка создана**\nТип: %s\nВремя: %s\nУченик: %s\nТренировок осталось: %d", training.GetType().GetName(), training.GetTime().Format("02.01 15:04"), training.GetStudent().GetName(), trainingsRemainingCount),
	)
	return &message
}
