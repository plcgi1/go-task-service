package main

import (
	"os"
	"os/signal"

	"go-task-service/config"
	"go-task-service/internal/handler"
	"go-task-service/internal/logger"
	"go-task-service/internal/repository"
	"go-task-service/internal/service"

	_ "go-task-service/docs" // Swagger docs

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Обработчик задач API
// @version 1.0
// @description Конкурентный обработчик задач
// @host localhost:8080
// @BasePath /
func main() {
	logger.Init()
	cfg := config.Load()
	db := repository.InitDB(cfg)

	logger.Log.Info("Starting Fiber server on :" + cfg.AppPort)

	sqlDB, _ := db.DB()

	logger.Log.Info("DB initialized")

	defer sqlDB.Close()

	repo := repository.NewTaskRepo(db)
	svc := service.NewTaskService(repo, cfg)
	app := fiber.New()

	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Post("/process", handler.New(svc, repo).ProcessHandler)
	app.Get("/tasks", handler.New(svc, repo).GetTasksHandler)

	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			logger.Log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Log.Info("Shutting down...")
	err := app.Shutdown()
	if err != nil {
		logger.Log.WithError(err).Info("Shutting down errored")
	}
}
