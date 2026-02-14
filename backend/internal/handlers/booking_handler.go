package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/indraprhmbd/allocra/internal/models"
	"github.com/indraprhmbd/allocra/internal/services"
)

type BookingHandler struct {
    bookingService *services.BookingService
}

func NewBookingHandler(bookingService *services.BookingService) *BookingHandler {
    return &BookingHandler{bookingService: bookingService}
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
    var req models.CreateBookingRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
    
    booking, err := h.bookingService.CreateBooking(c.Context(), &req)
    if err != nil {
        if err.Error() == "booking conflict detected" {
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(booking)
}

func (h *BookingHandler) ApproveBooking(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid booking ID",
        })
    }
    
    err = h.bookingService.ApproveBooking(c.Context(), id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.SendStatus(fiber.StatusOK)
}

func (h *BookingHandler) RejectBooking(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid booking ID",
        })
    }
    
    err = h.bookingService.RejectBooking(c.Context(), id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.SendStatus(fiber.StatusOK)
}

func (h *BookingHandler) GetAllBookings(c *fiber.Ctx) error {
    bookings, err := h.bookingService.GetAllBookings(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.JSON(bookings)
}

func (h *BookingHandler) GetBookingsByRoom(c *fiber.Ctx) error {
    roomID, err := strconv.Atoi(c.Query("room_id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid room ID",
        })
    }
    
    bookings, err := h.bookingService.GetBookingsByRoom(c.Context(), roomID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(bookings)
}

func (h *BookingHandler) GetMonthlyReport(c *fiber.Ctx) error {
    report, err := h.bookingService.GetMonthlyReport(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(report)
}

func (h *BookingHandler) ForceAllocate(c *fiber.Ctx) error {
    id, _ := strconv.Atoi(c.Params("id"))
    err := h.bookingService.ForceAllocate(c.Context(), id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.SendStatus(fiber.StatusOK)
}
