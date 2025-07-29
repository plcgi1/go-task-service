package handler

import (
	"runtime"

	"github.com/gofiber/fiber/v2"
)

func GoroutineDumpHandler(c *fiber.Ctx) error {
	buf := make([]byte, 1<<20) // 1MB буфер
	n := runtime.Stack(buf, true)
	return c.SendString(string(buf[:n]))
}
