package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"` // "admin" or "user"
    CreatedAt time.Time `json:"created_at"`
}

type Room struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Capacity  int       `json:"capacity"`
    Type      string    `json:"type"`   // "shared" or "exclusive"
    Status    string    `json:"status"` // "online", "maintenance", "offline"
    CreatedAt time.Time `json:"created_at"`
}

type Booking struct {
    ID        int       `json:"id"`
    RoomID    int       `json:"room_id"`
    UserID    int       `json:"user_id"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    Status    string    `json:"status"` // "pending", "approved", "rejected"
    CreatedAt time.Time `json:"created_at"`
}

// CreateBookingRequest represents the booking creation payload
type CreateBookingRequest struct {
    RoomID    int       `json:"room_id"`
    UserID    int       `json:"user_id"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
}

// MonthlyUsageReport represents aggregated room usage
type MonthlyUsageReport struct {
    RoomID        int     `json:"room_id"`
    RoomName      string  `json:"room_name"`
    TotalBookings int     `json:"total_bookings"`
    TotalHours    float64 `json:"total_hours"`
}
