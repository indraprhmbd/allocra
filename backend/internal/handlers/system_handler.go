package handlers

import (
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) GetStats(c *fiber.Ctx) error {
	// Simulate dynamic system metrics
	return c.JSON(fiber.Map{
		"cpu_usage":    10 + rand.Intn(15),
		"memory_usage": 40 + rand.Intn(20),
		"io_wait":      1 + rand.Intn(5),
		"status":       "HEALTHY",
	})
}
