# Go API

A simple **Go REST API** for managing users, backed by **PostgreSQL running in Docker**.  
The API runs locally on your machine and can be tested using **Postman** or `curl`.

---

## Tech Stack

- **Go** (standard library `net/http`)
- **PostgreSQL 16**
- **Docker + Docker Compose**
- **pgx** Postgres driver
- No authentication (by design)
- No frameworks

---

## Project Structure

```text
go-users-crud/
├── controllers/
│   └── user_controller.go   # HTTP handlers / routing
├── models/
│   └── user.go              # Data models & request DTOs
├── services/
│   └── user_service.go      # Business logic & DB access
├── docker-compose.yml       # Postgres container
├── go.mod
├── go.sum
└── main.go                  # App bootstrap (wiring only)
