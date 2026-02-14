package services

import (
	"context"
	"fmt"

	"github.com/indraprhmbd/allocra/internal/models"
	"github.com/indraprhmbd/allocra/internal/repository"
)

type RoomService struct {
    roomRepo *repository.RoomRepository
}

func NewRoomService(roomRepo *repository.RoomRepository) *RoomService {
    return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) CreateRoom(ctx context.Context, name string, capacity int, roomType string, status string) (*models.Room, error) {
    if capacity <= 0 {
        return nil, fmt.Errorf("capacity must be positive")
    }
    
    if roomType == "" { roomType = "shared" }
    if status == "" { status = "online" }
    
    return s.roomRepo.Create(ctx, name, capacity, roomType, status)
}

func (s *RoomService) GetAllRooms(ctx context.Context) ([]models.Room, error) {
    return s.roomRepo.GetAll(ctx)
}

func (s *RoomService) UpdateRoom(ctx context.Context, id int, name string, capacity int, roomType string, status string) error {
    if capacity <= 0 {
        return fmt.Errorf("capacity must be positive")
    }
    return s.roomRepo.Update(ctx, id, name, capacity, roomType, status)
}

func (s *RoomService) DeleteRoom(ctx context.Context, id int) error {
    return s.roomRepo.Delete(ctx, id)
}
