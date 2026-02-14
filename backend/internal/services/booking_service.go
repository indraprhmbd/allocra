package services

import (
	"context"
	"fmt"
	"time"

	"github.com/indraprhmbd/allocra/internal/models"
	"github.com/indraprhmbd/allocra/internal/repository"
)

type BookingService struct {
    bookingRepo *repository.BookingRepository
}

func NewBookingService(bookingRepo *repository.BookingRepository) *BookingService {
    return &BookingService{bookingRepo: bookingRepo}
}

// CreateBooking validates and creates a new booking
func (s *BookingService) CreateBooking(ctx context.Context, req *models.CreateBookingRequest) (*models.Booking, error) {
    // Business logic validation
    if req.StartTime.After(req.EndTime) || req.StartTime.Equal(req.EndTime) {
        return nil, fmt.Errorf("invalid time range: start must be before end")
    }
    
    // Allow a 2-minute grace period for "immediate" bookings to account for clock drift
    if req.StartTime.Before(time.Now().Add(-2 * time.Minute)) {
        return nil, fmt.Errorf("cannot book in the past (beyond 2min grace period)")
    }
    
    // Delegate to repository (transaction handled there)
    return s.bookingRepo.CreateWithTransaction(ctx, req)
}

func (s *BookingService) ApproveBooking(ctx context.Context, bookingID int) error {
    return s.bookingRepo.ApproveBooking(ctx, bookingID)
}

func (s *BookingService) RejectBooking(ctx context.Context, bookingID int) error {
    return s.bookingRepo.RejectBooking(ctx, bookingID)
}

func (s *BookingService) GetAllBookings(ctx context.Context) ([]models.Booking, error) {
    return s.bookingRepo.GetAll(ctx)
}

func (s *BookingService) GetBookingsByRoom(ctx context.Context, roomID int) ([]models.Booking, error) {
    return s.bookingRepo.GetByRoomID(ctx, roomID)
}

func (s *BookingService) GetMonthlyReport(ctx context.Context) ([]models.MonthlyUsageReport, error) {
    return s.bookingRepo.GetMonthlyUsage(ctx)
}

func (s *BookingService) ForceAllocate(ctx context.Context, bookingID int) error {
    return s.bookingRepo.PreemptBooking(ctx, bookingID)
}

func (s *BookingService) ResetAllocations(ctx context.Context) error {
    return s.bookingRepo.DeleteAll(ctx)
}

func (s *BookingService) GetSystemStats(ctx context.Context) (*repository.SystemStats, error) {
    return s.bookingRepo.GetSystemStats(ctx)
}
