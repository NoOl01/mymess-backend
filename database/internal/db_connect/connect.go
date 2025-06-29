package db_connect

import (
	"database/internal/db_models"
	"database/internal/storage"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"results/errs"
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		storage.Env.DbHost, storage.Env.DbUser, storage.Env.DbPass, storage.Env.DbName, storage.Env.DbSsl)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("%s: %s", errs.FailedDatabaseConnect, err.Error())
	}

	autoMigrateErr := db.AutoMigrate(&db_models.User{}, &db_models.Chat{}, &db_models.Message{}, &db_models.UserTheme{})
	if autoMigrateErr != nil {
		log.Fatal(autoMigrateErr)
	}

	return db
}
