package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/indraprhmbd/allocra/internal/models"
)

type RoomRepository struct {
    db *Database
}

func NewRoomRepository(db *Database) *RoomRepository {
    return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(ctx context.Context, name string, capacity int, roomType string, status string) (*models.Room, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    query := `
        INSERT INTO rooms (name, capacity, type, status)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, capacity, type, status, created_at
    `
    
    var room models.Room
    err := r.db.DB.QueryRowContext(ctx, query, name, capacity, roomType, status).Scan(
        &room.ID,
        &room.Name,
        &room.Capacity,
        &room.Type,
        &room.Status,
        &room.CreatedAt,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to create room: %w", err)
    }
    
    return &room, nil
}

func (r *RoomRepository) GetAll(ctx context.Context) ([]models.Room, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    query := `SELECT id, name, capacity, type, status, created_at FROM rooms ORDER BY name`
    
    rows, err := r.db.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch rooms: %w", err)
    }
    defer rows.Close()
    
    var rooms []models.Room
    for rows.Next() {
        var room models.Room
        if err := rows.Scan(&room.ID, &room.Name, &room.Capacity, &room.Type, &room.Status, &room.CreatedAt); err != nil {
            return nil, fmt.Errorf("failed to scan room: %w", err)
        }
        rooms = append(rooms, room)
    }
    
    return rooms, rows.Err()
}

func (r *RoomRepository) Update(ctx context.Context, id int, name string, capacity int, roomType string, status string) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    query := `UPDATE rooms SET name = $1, capacity = $2, type = $3, status = $4 WHERE id = $5`
    result, err := r.db.DB.ExecContext(ctx, query, name, capacity, roomType, status, id)
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rows == 0 {
        return fmt.Errorf("room not found with id: %d", id)
    }
    
    return nil
}

func (r *RoomRepository) Delete(ctx context.Context, id int) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    query := `DELETE FROM rooms WHERE id = $1`
    _, err := r.db.DB.ExecContext(ctx, query, id)
    return err
}
