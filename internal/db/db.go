package db

import (
	"log/slog"
	"os"
	"wb/internal/order"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", slog.String("dsn", dsn), slog.String("error", err.Error()))
		os.Exit(1) // Завершаем процесс при критической ошибке
	}

	logger.Info("Successfully connected to database", slog.String("dsn", dsn))

	// Автоматическая миграция
	err = db.AutoMigrate(&order.Order{})
	if err != nil {
		logger.Error("Failed to migrate database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	db = db.Debug()
	logger.Info("Database migration completed successfully")

	return db
}
