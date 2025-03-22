// @title Song Library API
// @version 1.0
// @description This is a sample server for a song library.
// @host localhost:8080
// @BasePath /
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"music-api/internal/config"
	"music-api/internal/handler"
	"music-api/internal/migration/migration"
	"music-api/internal/repository"
	"music-api/internal/service"
	"os"
)

func main() {
	// Загрузка конфигурации из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Инициализация логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Инициализация конфигурации
	cfg := config.NewConfig()

	// Инициализация базы данных
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Применение миграций
	if err := migration.RunMigrations(cfg.DatabaseURL); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Инициализация репозитория, сервиса и хэндлера
	songRepo := repository.NewSongRepository(db, logger)
	songService := service.NewSongService(songRepo, logger, cfg.APIURL) // Передаем apiURL
	songHandler := handler.NewSongHandler(songService, logger)

	// Инициализация роутера
	r := gin.Default()

	// Регистрация роутов
	r.GET("/songs", songHandler.GetSongs)
	r.GET("/songs/:id/text", songHandler.GetSongText)
	r.DELETE("/songs/:id", songHandler.DeleteSong)
	r.PUT("/songs/:id", songHandler.UpdateSong)
	r.POST("/songs", songHandler.AddSong)

	// Запуск сервера
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
