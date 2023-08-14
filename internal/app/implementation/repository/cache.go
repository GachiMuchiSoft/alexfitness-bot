package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		redis: client,
	}
}

type Cache struct {
	redis *redis.Client
}

func (c Cache) GetStudentTrainings(ctx context.Context, id uint) (int, error) {
	trainingsCount, err := c.redis.Get(ctx, fmt.Sprintf("%d-trainings-count", id)).Int()
	switch {
	case err == redis.Nil:
		return 0, nil
	case err != nil:
		return 0, fmt.Errorf("error during getting user[%d] training balance :%w", id, err)
	default:
		return trainingsCount, nil
	}
}

func (c Cache) DecrStudentTrainings(ctx context.Context, id uint, count int64) error {
	if err := c.redis.DecrBy(ctx, fmt.Sprintf("%d-trainings-count", id), count).Err(); err != nil {
		return fmt.Errorf("error during removing trainings to student: %w", err)
	}

	return nil
}

func (c Cache) IncStudentTrainings(ctx context.Context, id uint, count int64) error {
	if err := c.redis.IncrBy(ctx, fmt.Sprintf("%d-trainings-count", id), count).Err(); err != nil {
		return fmt.Errorf("error during adding trainings to student: %w", err)
	}

	return nil
}

const buttonDataExpiration = time.Minute * 30

func (c Cache) SetButtonData(ctx context.Context, id string, data string) error {
	if err := c.redis.Set(ctx, fmt.Sprintf("button-data-%s", id), data, buttonDataExpiration).Err(); err != nil {
		return fmt.Errorf("error during saving button data: %w", err)
	}

	return nil
}

func (c Cache) GetButtonData(ctx context.Context, id string) (string, error) {
	result, err := c.redis.Get(ctx, fmt.Sprintf("button-data-%s", id)).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return result, nil
}
