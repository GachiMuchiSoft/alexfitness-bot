package messages

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const RoutingKeyStudentShowTrainingBalance = "3"
const RoutingKeyStudentRemoveTraining = "4"

func StartStudent(chatID int64) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, "Вы зарегистрированы в системе как студент, функционала у студентов пока нет")
	//showTrainingBalanceData, _ := json.Marshal(RoutingData{Routing: RoutingKeyStudentShowTrainingBalance})
	//showRemoveTrainingData, _ := json.Marshal(RoutingData{Routing: RoutingKeyStudentRemoveTraining})
	//msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
	//	tgbotapi.NewInlineKeyboardRow(
	//		tgbotapi.NewInlineKeyboardButtonData("Посмотреть баланс тренировок", string(showTrainingBalanceData)),
	//	),
	//	tgbotapi.NewInlineKeyboardRow(
	//		tgbotapi.NewInlineKeyboardButtonData("Зафиксировать тренировку", string(showRemoveTrainingData)),
	//	),
	//)

	return &msg
}

func StudentShowTrainingBalance(chatID uint, count int) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(int64(chatID), fmt.Sprintf("Количество ваших тренировок = %d", count))
	return &msg
}

func StudentRemoveTrainingSuccess(chatID uint) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(int64(chatID), "Тренировка зафиксирована")
	return &msg
}

func StudentHasRunOutOfTrainings(chatID uint) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(int64(chatID), "Тренер зафиксировал вашу последнюю тренировку. У вас нет больше тренировок, оплатите блок у вашего тренера")
	return &msg
}
