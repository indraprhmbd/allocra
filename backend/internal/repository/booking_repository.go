package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/indraprhmbd/allocra/internal/models"
)

type BookingRepository struct {
    db *Database
}

func NewBookingRepository(db *Database) *BookingRepository {
    return &BookingRepository{db: db}
}

// CheckConflict detects time range overlap for approved bookings
// Conflict exists when: existing.start_time < new_end AND existing.end_time > new_start
func (r *BookingRepository) CheckConflict(ctx context.Context, tx *sql.Tx, roomID int, start, end time.Time) (bool, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()
    
    // This query uses the composite index idx_bookings_room_time
    query := `
        SELECT 1
        FROM bookings
        WHERE room_id = $1
          AND status = 'approved'
          AND start_time < $3
          AND end_time > $2
        LIMIT 1
    `
    
    var exists int
    err := tx.QueryRowContext(ctx, query, roomID, start, end).Scan(&exists)
    
    if err == sql.ErrNoRows {
        return false, nil // No conflict
    }
    if err != nil {
        return false, fmt.Errorf("conflict check failed: %w", err)
    }
    
    return true, nil // Conflict exists
}

// CreateWithTransaction creates a booking within a transaction
// Transaction boundary: conflict check + insert must be atomic
func (r *BookingRepository) CreateWithTransaction(ctx context.Context, req *models.CreateBookingRequest) (*models.Booking, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    // BEGIN TRANSACTION
    tx, err := r.db.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback() // Rollback if not committed
    
    // Step 1: Check for conflicts
    hasConflict, err := r.CheckConflict(ctx, tx, req.RoomID, req.StartTime, req.EndTime)
    if err != nil {
        return nil, err
    }
    if hasConflict {
        return nil, fmt.Errorf("booking conflict detected")
    }
    
    // Step 2: Insert booking with status = 'approved' (Auto-Approval)
    query := `
        INSERT INTO bookings (room_id, user_id, start_time, end_time, status)
        VALUES ($1, $2, $3, $4, 'approved')
        RETURNING id, room_id, user_id, start_time, end_time, status, created_at
    `
    
    var booking models.Booking
    err = tx.QueryRowContext(ctx, query,
        req.RoomID,
        req.UserID,
        req.StartTime,
        req.EndTime,
    ).Scan(
        &booking.ID,
        &booking.RoomID,
        &booking.UserID,
        &booking.StartTime,
        &booking.EndTime,
        &booking.Status,
        &booking.CreatedAt,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to insert booking: %w", err)
    }
    
    // COMMIT TRANSACTION
    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return &booking, nil
}

// ApproveBooking updates booking status to 'approved' with conflict re-check
func (r *BookingRepository) ApproveBooking(ctx context.Context, bookingID int) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    tx, err := r.db.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()
    
    // Fetch booking details
    var booking models.Booking
    query := `SELECT room_id, start_time, end_time, status FROM bookings WHERE id = $1 FOR UPDATE`
    err = tx.QueryRowContext(ctx, query, bookingID).Scan(
        &booking.RoomID,
        &booking.StartTime,
        &booking.EndTime,
        &booking.Status,
    )
    if err != nil {
        return fmt.Errorf("failed to fetch booking: %w", err)
    }
    
    if booking.Status != "pending" {
        return fmt.Errorf("booking is not pending")
    }
    
    // Re-check conflict before approval
    hasConflict, err := r.CheckConflict(ctx, tx, booking.RoomID, booking.StartTime, booking.EndTime)
    if err != nil {
        return err
    }
    if hasConflict {
        return fmt.Errorf("conflict detected, cannot approve")
    }
    
    // Update status
    updateQuery := `UPDATE bookings SET status = 'approved' WHERE id = $1`
    _, err = tx.ExecContext(ctx, updateQuery, bookingID)
    if err != nil {
        return fmt.Errorf("failed to approve booking: %w", err)
    }
    
    return tx.Commit()
}

// RejectBooking updates booking status to 'rejected'
func (r *BookingRepository) RejectBooking(ctx context.Context, bookingID int) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    query := `UPDATE bookings SET status = 'rejected' WHERE id = $1 AND status = 'pending'`
    result, err := r.db.DB.ExecContext(ctx, query, bookingID)
    if err != nil {
        return fmt.Errorf("failed to reject booking: %w", err)
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("booking not found or not pending")
    }
    
    return nil
}

// GetAll fetches all bookings across all rooms
func (r *BookingRepository) GetAll(ctx context.Context) ([]models.Booking, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    query := `
        SELECT id, room_id, user_id, start_time, end_time, status, created_at
        FROM bookings
        ORDER BY start_time DESC
    `
    
    rows, err := r.db.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch all bookings: %w", err)
    }
    defer rows.Close()
    
    var bookings []models.Booking
    for rows.Next() {
        var booking models.Booking
        if err := rows.Scan(
            &booking.ID,
            &booking.RoomID,
            &booking.UserID,
            &booking.StartTime,
            &booking.EndTime,
            &booking.Status,
            &booking.CreatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan booking: %w", err)
        }
        bookings = append(bookings, booking)
    }
    
    return bookings, rows.Err()
}

// GetByRoomID fetches bookings for a specific room
func (r *BookingRepository) GetByRoomID(ctx context.Context, roomID int) ([]models.Booking, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    query := `
        SELECT id, room_id, user_id, start_time, end_time, status, created_at
        FROM bookings
        WHERE room_id = $1
        ORDER BY start_time DESC
    `
    
    rows, err := r.db.DB.QueryContext(ctx, query, roomID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch bookings: %w", err)
    }
    defer rows.Close()
    
    var bookings []models.Booking
    for rows.Next() {
        var booking models.Booking
        if err := rows.Scan(
            &booking.ID,
            &booking.RoomID,
            &booking.UserID,
            &booking.StartTime,
            &booking.EndTime,
            &booking.Status,
            &booking.CreatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan booking: %w", err)
        }
        bookings = append(bookings, booking)
    }
    
    return bookings, rows.Err()
}

// GetMonthlyUsage aggregates approved bookings for current month
func (r *BookingRepository) GetMonthlyUsage(ctx context.Context) ([]models.MonthlyUsageReport, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    query := `
        SELECT 
            b.room_id,
            r.name as room_name,
            COUNT(*) as total_bookings,
            SUM(EXTRACT(EPOCH FROM (b.end_time - b.start_time))/3600) as total_hours
        FROM bookings b
        JOIN rooms r ON b.room_id = r.id
        WHERE b.status = 'approved'
          AND date_trunc('month', b.start_time) = date_trunc('month', NOW())
        GROUP BY b.room_id, r.name
        ORDER BY total_hours DESC
    `
    
    rows, err := r.db.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch monthly usage: %w", err)
    }
    defer rows.Close()
    
    var reports []models.MonthlyUsageReport
    for rows.Next() {
        var report models.MonthlyUsageReport
        if err := rows.Scan(
            &report.RoomID,
            &report.RoomName,
            &report.TotalBookings,
            &report.TotalHours,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan usage report: %w", err)
        }
        reports = append(reports, report)
    }
    
    return reports, rows.Err()
}

// PreemptBooking cancels existing approved bookings that conflict and approves the new one
func (r *BookingRepository) PreemptBooking(ctx context.Context, bookingID int) error {
    ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
    defer cancel()
    
    tx, err := r.db.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 1. Fetch the target booking
    var b models.Booking
    err = tx.QueryRowContext(ctx, "SELECT room_id, start_time, end_time FROM bookings WHERE id = $1 FOR UPDATE", bookingID).Scan(
        &b.RoomID, &b.StartTime, &b.EndTime,
    )
    if err != nil {
        return err
    }
    
    // 2. Reject all overlapping 'approved' bookings for this room
    rejectQuery := `
        UPDATE bookings 
        SET status = 'rejected' 
        WHERE room_id = $1 
          AND status = 'approved' 
          AND start_time < $3 
          AND end_time > $2
    `
    _, err = tx.ExecContext(ctx, rejectQuery, b.RoomID, b.StartTime, b.EndTime)
    if err != nil {
        return err
    }
    
    // 3. Approve the target booking
    _, err = tx.ExecContext(ctx, "UPDATE bookings SET status = 'approved' WHERE id = $1", bookingID)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}
