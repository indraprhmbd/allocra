package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
    DB *sql.DB
}

// NewDatabase creates a new database connection with production settings
func NewDatabase(host, port, user, password, dbname string) (*Database, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Production connection pool settings for 1GB VPS
    db.SetMaxOpenConns(10)          // Limit concurrent connections
    db.SetMaxIdleConns(5)            // Keep some idle connections
    db.SetConnMaxLifetime(time.Hour) // Recycle connections hourly
    
    // Verify connection with retries
    var lastErr error
    for i := 0; i < 10; i++ {
        if err = db.Ping(); err == nil {
            log.Println("Database connection established")
            return &Database{DB: db}, nil
        }
        lastErr = err
        log.Printf("Attempt %d: Failed to ping database: %v. Retrying in 2s...", i+1, err)
        time.Sleep(2 * time.Second)
    }
    
    return nil, fmt.Errorf("failed to connect to database after retries: %w", lastErr)
}

func (d *Database) Close() error {
    return d.DB.Close()
}
