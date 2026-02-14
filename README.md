# Allocra

Resource allocation engine focusing on transactional safety and real-time observability. Built with a Go backend and a Vue 3 frontend.

## Project Structure

- backend: Go backend (Fiber, Postgres)
- frontend: Vue 3 frontend (Vite, Tailwind)
- backend/migrations: SQL migration files
- docker-compose.yml: Orchestration for local development

## Getting Started

Prerequisites:

- Docker and Docker Compose

Running the Application:

1. Start the services:

   ```bash
   docker-compose up -d --build
   ```

   This command starts the PostgreSQL database, the Go API, and the Vue frontend.

2. Access the application:
   - Frontend: http://localhost:5173
   - API: http://localhost:8080

## API Endpoints

### Resources

- GET /api/rooms - List all registered resources
- POST /api/rooms - Register a new resource
- PUT /api/rooms/:id - Update resource metadata
- DELETE /api/rooms/:id - Decommission a resource

### Allocations (Bookings)

- GET /api/bookings/all - Fetch all allocation requests
- POST /api/bookings - Submit a new allocation (automatic conflict detection)
- PATCH /api/bookings/:id/approve - Approve a pending request
- PATCH /api/bookings/:id/reject - Deny a request
- PATCH /api/bookings/:id/force - Preempt existing allocations (Force Allocate)

### System & Metrics

- GET /api/system/stats - Fetch real-time engine load and system health
- GET /api/reports/monthly-usage - Retrieve monthly utilization data

## Architecture and Tech Stack

Backend:

- Language: Go 1.21
- Framework: Fiber v2
- Database: PostgreSQL 16
- Core Logic: 3-layer architecture (Handler, Service, Repository) with strict transactional isolation

Frontend:

- Framework: Vue 3 (Composition API)
- Build Tool: Vite
- Styling: Tailwind CSS
- Icons: Tabler Icons

Transactional Safety:
The system uses Read Committed isolation and row-level locking (FOR UPDATE) to prevent overbooking and race conditions during high-concurrency resource allocation.
