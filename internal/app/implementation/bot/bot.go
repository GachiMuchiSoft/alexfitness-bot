package bot

import (
	"context"
	"encoding/json"
	"fmt"
	dmessages "gachiSoft/alexfitness-bot/internal/app/domain/messages"
	"gachiSoft/alexfitness-bot/internal/app/domain/models/user"
	imessages "gachiSoft/alexfitness-bot/internal/app/implementation/bot/messages"
	iuser "gachiSoft/alexfitness-bot/internal/app/implementation/models"
	"gachiSoft/alexfitness-bot/internal/app/implementation/repository"
	"gachiSoft/alexfitness-bot/internal/pkg/strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"reflect"
)

type Bot struct {
	cache *repository.Cache
	bot   *tgbotapi.BotAPI
}

func (b *Bot) Notify(u user.User, message interface{}) error {
	user, _ := u.(*iuser.User)
	ctx := context.Background()
	switch message.(type) {
	case dmessages.CoachTrainingCreated:
		msg := message.(dmessages.CoachTrainingCreated)
		err := b.Send(ctx, imessages.CoachTrainingCreated(user.ChatID, msg.Training, msg.TrainingRemainingCount))
		return err
	case dmessages.CoachStudentHasRunOutOfTrainings:
		err := b.Send(ctx, imessages.CoachAddedTrainingToStudentAndStudentHasRunOutOfTrainings(user.ChatID))
		return err
	case dmessages.CoachStudentHasNotEnoughAvailableTrainings:
		err := b.Send(ctx, imessages.CoachUnableToAddTrainingStudentHasNotEnoughAvailableTrainings(user.ChatID))
		return err
	case dmessages.CoachAddedAvailableTrainingsToStudent:
		msg := message.(dmessages.CoachAddedAvailableTrainingsToStudent)
		err := b.Send(ctx, imessages.CoachSuccessfullyAddTrainingsBlock(user.ChatID, msg.Count))
		return err
	case dmessages.StudentLastTrainingFinished:
		err := b.Send(ctx, imessages.StudentHasRunOutOfTrainings(user.ChatID))
		return err
	}

	return fmt.Errorf("unknown message type: %s", reflect.TypeOf(message))
}

func New(cache *repository.Cache, api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		cache: cache,
		bot:   api,
	}
}

type CallbackData struct {
	ID string `json:"id"`
}

func (b *Bot) Send(ctx context.Context, message *tgbotapi.MessageConfig) (err error) {
	message.ParseMode = tgbotapi.ModeMarkdown
	if message.ReplyMarkup != nil {
		keyboard := message.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
		if err = b.saveCompressedKeyboardData(ctx, &keyboard); err != nil {
			return err
		}
	}

	_, err = b.bot.Send(message)

	return err
}

func (b Bot) saveCompressedKeyboardData(ctx context.Context, keyboard *tgbotapi.InlineKeyboardMarkup) error {
	if len(keyboard.InlineKeyboard) != 0 {
		for bri, buttonRow := range keyboard.InlineKeyboard {
			for bi, button := range buttonRow {
				id := strings.RandomString(10)
				marshalled, _ := json.Marshal(CallbackData{id})

				if err := b.cache.SetButtonData(ctx, id, *button.CallbackData); err != nil {
					return err
				}

				marshalledData := string(marshalled)
				keyboard.InlineKeyboard[bri][bi].CallbackData = &marshalledData
			}
		}
	}

	return nil
}
