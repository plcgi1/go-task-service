package main

import (
	"context"

	"go-task-service/config"
	"go-task-service/internal/handler"
	"go-task-service/internal/lifecycle"
	"go-task-service/internal/logger"
	"go-task-service/internal/metrics"
	"go-task-service/internal/repository"
	"go-task-service/internal/service"

	_ "go-task-service/docs" // Swagger docs

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Обработчик задач API
// @version 1.0
// @description Конкурентный обработчик задач
// @host localhost:8080
// @BasePath /
func main() {
	metrics.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Init()
	cfg := config.Load()
	db := repository.InitDB(cfg)

	logger.Log.Info("Starting Fiber server on :" + cfg.AppPort)

	sqlDB, _ := db.DB()

	logger.Log.Info("DB initialized")

	defer sqlDB.Close()

	repo := repository.NewTaskRepo(db)
	svc := service.NewTaskService(ctx, repo, cfg)
	app := fiber.New()

	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Post(
		"/process",
		handler.New(svc, repo).ProcessHandler,
	)
	app.Get("/tasks", handler.New(svc, repo).GetTasksHandler)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	app.Get("/debug/stack", handler.GoroutineDumpHandler)

	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			logger.Log.Fatal(err)
		}
	}()

	lifecycle.SetupGracefulShutdown(app, cancel)
}
