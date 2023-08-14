package environment

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const (
	keyBotToken    = "BOT_TOKEN"
	keyEnvironment = "ENVIRONMENT"
	keyTestUserID  = "TEST_USER_ID"
)

const (
	keyPostgresDB       = "POSTGRES_DB"
	keyPostgresUser     = "POSTGRES_USER"
	keyPostgresPort     = "POSTGRES_PORT"
	keyPostgresHost     = "POSTGRES_HOST"
	keyPostgresPassword = "POSTGRES_PASSWORD"
)

const (
	keyRedisHost     = "REDIS_HOST"
	keyRedisPort     = "REDIS_PORT"
	keyRedisPassword = "REDIS_PASSWORD"
)

func init() {
	dir, _ := os.Getwd()
	err := godotenv.Load(fmt.Sprintf("%s/.env", dir))
	if err != nil {
		panic(fmt.Errorf("error during loading environment: %w", err))
	}
}

func IsProduction() bool {
	return os.Getenv(keyEnvironment) == "prod"
}

func BotToken() string {
	return os.Getenv(keyBotToken)
}

func TestUserID() (uint, error) {
	id, err := strconv.Atoi(os.Getenv(keyTestUserID))
	if err != nil {
		panic(err)
	}

	return uint(id), nil
}

func PostgresDB() string {
	return os.Getenv(keyPostgresDB)
}

func PostgresUser() string {
	return os.Getenv(keyPostgresUser)
}

func PostgresPassword() string {
	return os.Getenv(keyPostgresPassword)
}

func PostgresPort() string {
	return os.Getenv(keyPostgresPort)
}

func PostgresHost() string {
	return os.Getenv(keyPostgresHost)
}

func RedisPassword() string {
	return os.Getenv(keyRedisPassword)
}

func RedisHost() string {
	return os.Getenv(keyRedisHost)
}

func RedisPort() string {
	return os.Getenv(keyRedisPort)
}
