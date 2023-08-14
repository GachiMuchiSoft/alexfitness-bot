package main

import (
	"fmt"
	"gachiSoft/alexfitness-bot/internal/app/implementation/environment"
	"gachiSoft/alexfitness-bot/internal/app/implementation/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func main() {
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
		panic(fmt.Errorf("error during connecting to database: %w", err))
	}

	tables := []interface{}{&models.User{}, &models.CoachStudent{}, &models.TrainingType{}, &models.Training{}}
	err = db.AutoMigrate(tables...)
	if err != nil {
		panic(fmt.Errorf("error during database migrations: %w", err))
	}

	if environment.IsProduction() {
		seedProductionEnvironment(db)
	} else {
		seedTestEnvironment(db)
	}
}

func seedProductionEnvironment(db *gorm.DB) {
	db.Create(models.NewUserCoach("Gagik", 381671645))
	db.Create(models.NewUserStudent("Андрюша", 381671645))
	db.Create(models.NewUserStudent("Женя", 381671645))
	db.Create(models.NewCoachStudent(1, 2))
	db.Create(models.NewCoachStudent(1, 3))
}

func seedTestEnvironment(db *gorm.DB) {
	userGagik := models.NewUserCoach("Gagik", 381671645)
	db.Create(&userGagik)

	gagikTrainingTypeLegs := models.NewTrainingType("ноги", userGagik.ID)
	db.Create(&gagikTrainingTypeLegs)
	gagikTrainingTypeArms := models.NewTrainingType("руки", userGagik.ID)
	db.Create(&gagikTrainingTypeArms)
	gagikTrainingTypeBody := models.NewTrainingType("вверх", userGagik.ID)
	db.Create(&gagikTrainingTypeBody)

	userAndrey := models.NewUserStudent("Андрей", 381671645)
	db.Create(&userAndrey)

	db.Create(models.NewCoachStudent(userGagik.ID, userAndrey.ID))

	andreyTraining := models.NewTraining(userAndrey.ID, gagikTrainingTypeArms.ID, time.Now())
	db.Create(&andreyTraining)

	userEvgeniy := models.NewUserStudent("Женя", 381671645)
	db.Create(&userEvgeniy)
	db.Create(models.NewCoachStudent(userGagik.ID, userEvgeniy.ID))
	evgeniyTraining := models.NewTraining(userEvgeniy.ID, gagikTrainingTypeLegs.ID, time.Now())
	db.Create(&evgeniyTraining)
}
