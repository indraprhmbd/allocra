package handlers

import (
	"github.com/indraprhmbd/allocra/internal/services"

	"github.com/gofiber/fiber/v2"
)

type SystemHandler struct {
    bookingService *services.BookingService
}

func NewSystemHandler(bookingService *services.BookingService) *SystemHandler {
    return &SystemHandler{bookingService: bookingService}
}

func (h *SystemHandler) GetStats(c *fiber.Ctx) error {
    stats, err := h.bookingService.GetSystemStats(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(fiber.Map{
        "status": "nominal",
        "uptime": "99.9%",
        "version": "1.0.0",
        "total_bookings": stats.TotalBookings,
        "active_bookings": stats.ActiveBookings,
        "conflicts": stats.Conflicts,
        "utilization": stats.Utilization,
        "cpu_usage": stats.CPUUsage,
        "memory_usage": stats.MemoryUsage,
    })
}

func (h *SystemHandler) ResetAllocations(c *fiber.Ctx) error {
    if err := h.bookingService.ResetAllocations(c.Context()); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"message": "All allocations have been reset"})
}
