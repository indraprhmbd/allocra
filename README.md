# Allocra

Resource allocation engine focusing on transactional safety and real-time observability. Built with a Go backend and a Vue 3 frontend, designed for deterministic stress testing and high-concurrency environments.

_Mesin alokasi sumber daya yang berfokus pada keamanan transaksi dan observabilitas real-time. Dibangun dengan backend Go dan frontend Vue 3, dirancang untuk pengujian stres deterministik dan lingkungan dengan konkurensi tinggi._

## Core Features / Fitur Utama

- **Transactional Integrity:** Row-level locking (`FOR UPDATE`) eliminates overbooking and race conditions.
  _(Integritas Transaksi: Penguncian tingkat baris (`FOR UPDATE`) menghilangkan risiko overbooking dan race conditions.)_
- **Conflict Persistence:** All allocation failures (conflicts) are logged as `REJECTED` state for auditability.
  _(Persistensi Konflik: Semua kegagalan alokasi (konflik) dicatat sebagai status `REJECTED` untuk kebutuhan audit.)_
- **Real-time Engine Load:** Live dashboard tracking occupancy rates, engine stress, and allocation velocity.
  _(Beban Mesin Real-time: Dashboard langsung untuk melacak tingkat hunian, stres mesin, dan kecepatan alokasi.)_
- **Deterministic Playground:** Integrated stress-testing unit to simulate sequential or parallel request storms.
  _(Playground Deterministik: Unit pengujian stres terintegrasi untuk mensimulasikan badai permintaan sekuensial atau paralel.)_
- **Timezone Standardization:** System-wide alignment to `Asia/Jakarta (WIB)` for precise temporal grids.
  _(Standarisasi Zona Waktu: Penyelarasan seluruh sistem ke `Asia/Jakarta (WIB)` untuk kisi temporal yang presisi.)_

## Project Structure

- `backend`: Go backend (Fiber, Postgres)
- `frontend`: Vue 3 frontend (Vite, TailwindCSS)
- `backend/migrations`: SQL schema and migration files
- `docker-compose.yml`: Full-stack orchestration (DB, API, Frontend, Nginx)

## Getting Started

Prerequisites:

- Docker and Docker Compose

Running the Application:

1. Start the services:

   ```bash
   docker-compose up -d --build
   ```

2. Access the application:
   - Dashboard: http://localhost:5173
   - Playground (Stress Test): http://localhost:5173/playground
   - API: http://localhost:8080

## API Endpoints

### Resources (Nodes)

- `GET /api/rooms` - List all registered resource nodes
- `POST /api/rooms` - Register a new resource
- `PUT /api/rooms/:id` - Update resource metadata and status (`online`/`offline`)
- `DELETE /api/rooms/:id` - Decommission a resource

### Allocations (Engine Logic)

- `GET /api/bookings/all` - Fetch all allocation history and conflicts
- `POST /api/bookings` - Submit new allocation (atomic conflict detection)
- `PATCH /api/bookings/:id/force` - Preempt existing allocations (Engine Override)
- `POST /api/allocations/reset` - Purge all allocations (Playground reset)

### Observability

- `GET /api/system/stats` - Fetch real-time engine load (CPU, Memory simulation)
- `GET /api/reports/monthly-usage` - Retrieve monthly utilization metrics

## Tech Stack / Tumpukan Teknologi

### Backend

- **Framework:** [Go Fiber v2](https://gofiber.io/) - High-performance, Express-inspired web framework for Go.
- **Runtime:** Go 1.21+
- **Database Driver:** `lib/pq` (Pure Go Postgres driver).
- **Configuration:** `godotenv` for environment-based configuration.
- **Security:** Standard library `crypto` for password hashing and deterministic UUID generation.

### Frontend

- **Framework:** [Vue 3](https://vuejs.org/) (Composition API) with Script Setup.
- **Build Tool:** [Vite 7](https://vitejs.dev/) - Next generation frontend tooling.
- **Language:** TypeScript for type-safe components.
- **State Management:** [Pinia](https://pinia.vuejs.org/) for reactive global state tracking.
- **Styling:** Tailwind CSS 4.0 with customized glassmorphism theme.
- **Icons:** [Tabler Icons](https://tabler.io/icons) for consistent UI iconography.
- **HTTP Client:** Axios with interceptors for global error handling.

### Infrastructure

- **Containerization:** Docker & Docker Compose for idempotent environments.
- **Proxy/Web Server:** Nginx (configured as reverse proxy for fallback SPA routing).
- **API Gateway:** Centralized handling of `/api` routing via Nginx load balancer simulation.

## Deterministic Safety / Keamanan Deterministik

The engine utilizes `READ COMMITTED` isolation levels combined with explicit row-level locking on resource nodes during the allocation window check. This ensures that even under parallel request storms (simulated in the Playground), the system maintains 100% allocation accuracy.

_Mesin ini menggunakan tingkat isolasi `READ COMMITTED` yang dikombinasikan dengan penguncian tingkat baris (row-level locking) eksplisit pada node sumber daya selama pemeriksaan jendela alokasi. Hal ini memastikan bahwa bahkan di bawah badai permintaan paralel (yang disimulasikan di Playground), sistem tetap mempertahankan akurasi alokasi 100%._
