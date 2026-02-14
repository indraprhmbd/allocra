# Allocra

Experimental resource allocation engine focusing on SQL correctness and transactional safety.

## Project Structure

- `backend/`: Go backend (Fiber, Postgres)
- `backend/migrations/`: SQL migration files
- `docker-compose.yml`: Local development setup

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Running the Application

1. **Start the services:**

   ```bash
   docker-compose up --build
   ```

   This will start the PostgreSQL database and the Go backend.

2. **Access the API:**
   The API will be available at `http://localhost:8080`.

### Database & Migrations

The database is automatically initialized with the schema and seed data defined in `backend/migrations` when the container starts for the first time.

## API Endpoints

### Rooms

- `POST /api/rooms` - Create a new room
- `GET /api/rooms` - List all rooms

### Bookings

- `POST /api/bookings` - Create a booking (checks for conflicts)
- `GET /api/bookings?room_id=<id>` - List bookings for a room
- `PATCH /api/bookings/:id/approve` - Approve a booking (re-checks conflicts)
- `PATCH /api/bookings/:id/reject` - Reject a booking

### Reports

- `GET /api/reports/monthly-usage` - Get monthly usage statistics

## Architecture

- **Language:** Go 1.21
- **Framework:** Fiber v2
- **Database:** PostgreSQL 15
- **Concurrency:** Optimistic concurrency control via database transactions and isolation levels.
