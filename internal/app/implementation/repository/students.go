package repository

import (
	"context"
	"fmt"
	"gachiSoft/alexfitness-bot/internal/app/domain/models/user"
	"github.com/redis/go-redis/v9"
)

type Student struct {
	redis *redis.Client
}

func NewStudent(redis *redis.Client) *Student {
	return &Student{
		redis: redis,
	}
}

func (s *Student) AddAvailableTrainingsQuantity(student user.Student, quantity int64) (int, error) {
	result, err := s.redis.IncrBy(context.Background(), fmt.Sprintf("%d-trainings-count", student.GetID()), quantity).Uint64()
	if err != nil {
		return 0, fmt.Errorf("error during adding trainings to student: %w", err)
	}

	return int(result), nil
}

func (s *Student) RemoveAvailableTrainingsQuantity(student user.Student, quantity int64) error {
	if err := s.redis.DecrBy(context.Background(), fmt.Sprintf("%d-trainings-count", student.GetID()), quantity).Err(); err != nil {
		return fmt.Errorf("error during removing student trainings: %w", err)
	}

	return nil
}

func (c *Student) GetTrainings(student user.Student) (int, error) {
	trainingsCount, err := c.redis.Get(context.Background(), fmt.Sprintf("%d-trainings-count", student.GetID())).Int()
	switch {
	case err == redis.Nil:
		return 0, nil
	case err != nil:
		return 0, fmt.Errorf("error during getting user[%d] training balance :%w", student.GetID(), err)
	default:
		return trainingsCount, nil
	}
}
