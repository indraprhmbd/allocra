package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/indraprhmbd/allocra/internal/handlers"
	"github.com/indraprhmbd/allocra/internal/repository"
	"github.com/indraprhmbd/allocra/internal/services"
)

func main() {
    // Set default timezone to WIB (Asia/Jakarta)
    loc, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        log.Printf("Warning: failed to load Asia/Jakarta timezone: %v", err)
        // Fallback or just continue (system might handle it via TZ env)
    } else {
        time.Local = loc
        log.Println("Timezone configured to Asia/Jakarta (WIB)")
    }

    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Database connection
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    
    // Default values if not set
    if dbHost == "" { dbHost = "localhost" }
    if dbPort == "" { dbPort = "5432" }
    if dbUser == "" { dbUser = "postgres" }
    if dbPassword == "" { dbPassword = "password" }
    if dbName == "" { dbName = "allocra" }
    
    db, err := repository.NewDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // Run migrations
    if err := runMigrations(db.DB); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
    
    // Wire up dependencies
    roomRepo := repository.NewRoomRepository(db)
    bookingRepo := repository.NewBookingRepository(db)
    
    roomService := services.NewRoomService(roomRepo)
    bookingService := services.NewBookingService(bookingRepo)
    
    roomHandler := handlers.NewRoomHandler(roomService)
    bookingHandler := handlers.NewBookingHandler(bookingService)
    systemHandler := handlers.NewSystemHandler(bookingService)
    
    // Initialize Fiber
    app := fiber.New()
    
    // Middleware
    app.Use(logger.New())
    app.Use(recover.New())
    
    // Routes
    api := app.Group("/api")
    
    // Room routes
    api.Post("/rooms", roomHandler.CreateRoom)
    api.Get("/rooms", roomHandler.GetRooms)
    api.Put("/rooms/:id", roomHandler.UpdateRoom)
    api.Delete("/rooms/:id", roomHandler.DeleteRoom)
    
    // Booking routes
    api.Get("/bookings/all", bookingHandler.GetAllBookings)
    api.Post("/bookings", bookingHandler.CreateBooking)
    api.Patch("/bookings/:id/approve", bookingHandler.ApproveBooking)
    api.Patch("/bookings/:id/reject", bookingHandler.RejectBooking)
    api.Patch("/bookings/:id/force", bookingHandler.ForceAllocate)
    api.Get("/reports/monthly-usage", bookingHandler.GetMonthlyReport)
    
    // System routes
    api.Get("/system/stats", systemHandler.GetStats)
    api.Post("/allocations/reset", systemHandler.ResetAllocations)
    
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

func runMigrations(db *sql.DB) error {
    // Create migrations table if it doesn't exist
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS _migrations (name TEXT PRIMARY KEY)`)
    if err != nil {
        return fmt.Errorf("failed to create migrations table: %w", err)
    }

    log.Println("Checking for pending migrations...")
    
    // Retroactively seed migrations if tables exist but tracking is missing
    var roomsExist bool
    db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'rooms')").Scan(&roomsExist)
    if roomsExist {
        db.Exec("INSERT INTO _migrations (name) VALUES ('migrations/001_initial_schema.sql') ON CONFLICT DO NOTHING")
        db.Exec("INSERT INTO _migrations (name) VALUES ('migrations/002_seed_data.sql') ON CONFLICT DO NOTHING")
    }

    files := []string{
        "migrations/001_initial_schema.sql", 
        "migrations/002_seed_data.sql",
        "migrations/003_extend_rooms.sql",
    }

    for _, file := range files {
        var applied bool
        err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM _migrations WHERE name = $1)", file).Scan(&applied)
        if err != nil {
            return fmt.Errorf("failed to check migration status: %w", err)
        }

        if applied {
            continue
        }

        log.Printf("Applying migration: %s", file)
        content, err := os.ReadFile(file)
        if err != nil {
            return fmt.Errorf("failed to read migration file %s: %w", file, err)
        }
        
        _, err = db.Exec(string(content))
        if err != nil {
            return fmt.Errorf("failed to execute migration %s: %w", file, err)
        }

        _, err = db.Exec("INSERT INTO _migrations (name) VALUES ($1)", file)
        if err != nil {
            return fmt.Errorf("failed to record migration %s: %w", file, err)
        }
        log.Printf("Successfully applied migration: %s", file)
    }
    return nil
}
