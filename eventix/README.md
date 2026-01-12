# Eventix Backend

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![Gin Framework](https://img.shields.io/badge/Gin-1.11-00ADD8?style=flat)
![GORM](https://img.shields.io/badge/GORM-1.31-blue?style=flat)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

A scalable, production-ready backend for an event ticketing platform. Built with **Go** featuring concurrency-safe ticket booking with mutex locks and asynchronous notification processing via goroutines and channels.

---

## Key Features

### 🔒 Concurrency-Safe Booking
Prevents overselling with **sync.Mutex** locking during the critical section of ticket reservation. Combined with **GORM transactions** to ensure atomicity when decrementing available tickets and creating orders.

```go
s.bookingMutex.Lock()
defer s.bookingMutex.Unlock()
```

### ⚡ Asynchronous Processing
Fire-and-forget email notifications using **buffered channels** and **goroutine workers**. Payment confirmation triggers async job dispatch, ensuring fast API response times.

```go
emailChan := make(chan worker.EmailJob, 100)
worker.StartEmailWorker(emailChan)
```

### 🔐 Security
- **JWT Authentication** with role-based claims (user/admin)
- **Bcrypt** password hashing
- Protected routes with middleware chain

---

## Project Architecture

```
eventix/
├── cmd/
│   └── api/
│       └── main.go              
├── internal/
│   ├── entity/                  
│   ├── repository/              
│   ├── service/                 
│   ├── handler/                 
│   └── middleware/              
├── pkg/
│   ├── database/                
│   ├── utils/                   
│   └── worker/                  
├── .env                         
├── go.mod
└── go.sum
```

---

## Getting Started

### Prerequisites

- **Go** 1.22 or higher
- **PostgreSQL** 12 or higher
- **Git**

### Installation

```bash
cd eventix

# Install dependencies
go mod download

# Setup environment variables
cp .env.example .env
```

### Environment Configuration

Edit a `.env` file in the project root:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_pass
DB_NAME=eventix
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key

# Server
SERVER_PORT=8080
```

### Database Setup

```bash
# Create database
createdb eventix

# Tables are auto-migrated on startup
```

### Run the Application

```bash
# Development
go run cmd/api/main.go

# Build and run
go build -o eventix cmd/api/main.go
./eventix
```

Server starts at `http://localhost:8080`

---

## API Documentation

### Health Check

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/health` | No | Service health status |

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/auth/register` | No | Register new user |
| POST | `/api/auth/login` | No | Login, returns JWT token |

### Users

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/users/profile` | User | Get authenticated user profile |

### Events

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/events` | No | List events (supports `search`, `location`, `date_from`, `date_to`, `page`, `page_size` query params) |
| GET | `/api/events/:id` | No | Get event details |
| POST | `/api/events` | Admin | Create new event |
| PUT | `/api/events/:id` | Admin | Update event |
| DELETE | `/api/events/:id` | Admin | Delete event |
| POST | `/api/events/:id/book` | User | Book tickets for event |

### Orders

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/orders` | User | List user's orders |
| GET | `/api/orders/:id` | User | Get order details with tickets |
| POST | `/api/orders/:id/pay` | User | Process payment, generates tickets |
| POST | `/api/orders/:id/cancel` | User | Cancel pending order |

### Request/Response Examples

**Register User:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "secret123"}'
```

**Book Tickets:**
```bash
curl -X POST http://localhost:8080/api/events/1/book \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"qty": 2}'
```

**Process Payment:**
```bash
curl -X POST http://localhost:8080/api/orders/1/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"payment_method": "credit_card"}'
```

---

## License

This project is licensed under the MIT License.
