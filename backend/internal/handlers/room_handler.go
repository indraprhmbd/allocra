package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/indraprhmbd/allocra/internal/services"
)

type RoomHandler struct {
    roomService *services.RoomService
}

func NewRoomHandler(roomService *services.RoomService) *RoomHandler {
    return &RoomHandler{roomService: roomService}
}

type CreateRoomRequest struct {
    Name     string `json:"name"`
    Capacity int    `json:"capacity"`
    Type     string `json:"type"`
    Status   string `json:"status"`
}

func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
    var req CreateRoomRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
    
    room, err := h.roomService.CreateRoom(c.Context(), req.Name, req.Capacity, req.Type, req.Status)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(room)
}

func (h *RoomHandler) GetRooms(c *fiber.Ctx) error {
    rooms, err := h.roomService.GetAllRooms(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.JSON(rooms)
}

func (h *RoomHandler) UpdateRoom(c *fiber.Ctx) error {
    id, _ := strconv.Atoi(c.Params("id"))
    var req CreateRoomRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
    }
    
    err := h.roomService.UpdateRoom(c.Context(), id, req.Name, req.Capacity, req.Type, req.Status)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.SendStatus(200)
}

func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
    id, _ := strconv.Atoi(c.Params("id"))
    err := h.roomService.DeleteRoom(c.Context(), id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.SendStatus(200)
}
