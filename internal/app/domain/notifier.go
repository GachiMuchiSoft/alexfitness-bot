package domain

import "gachiSoft/alexfitness-bot/internal/app/domain/models/user"

type Notifier interface {
	Notify(user user.User, message interface{}) error
}
