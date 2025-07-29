package lifecycle

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go-task-service/internal/logger"

	"github.com/gofiber/fiber/v2"
)

func SetupGracefulShutdown(app *fiber.App, cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	cancel()
	_ = app.Shutdown()
	logger.Log.Println("Сервер успешно завершил работу")
}
