# Lamina

Lamina is a modern SaaS platform scaffolded in Go, using clean architecture principles, dependency injection, and full TDD.

## 🚀 Features

- Signup and Login with hashed passwords and JWTs
- Role-based access control (RBAC)
- Auth middleware to protect routes
- Admin-only endpoint to create users
- TDD-first development (Go + Testify)
- Dockerized with PostgreSQL
- Modular architecture with clearly separated layers

## 📦 Tech Stack

- **Go 1.24**
- **Gin** (web framework)
- **PostgreSQL**
- **Docker + docker-compose**
- **Testify** for unit tests

---

## 🔐 Auth Flow

| Endpoint                      | Access Type       | Description                           |
|-------------------------------|-------------------|---------------------------------------|
| `POST /api/v1/auth/signup`    | Public            | Signup new user (role: user)          |
| `POST /api/v1/auth/login`     | Public            | Login with email + password           |
| `GET /api/v1/user/me`         | Authenticated     | Get current user info (JWT required)  |
| `POST /api/v1/admin/create-user` | Admin Only     | Create new user manually              |

JWT includes: `userID`, `email`, `role`  
Middleware extracts and injects claims into request context.

---

## 🧪 Testing

All business logic is tested using TDD and mocks.

### Unit Tests:
| File                                         | Tested Component                |
|----------------------------------------------|----------------------------------|
| `internal/auth/service_test.go`              | `Login()`, `SignupUser()` logic |
| `internal/auth/auth_middleware_test.go`      | JWT middleware                   |
| `common/utils/jwt_test.go`                   | Token generation and parsing     |
| `common/utils/password_test.go`              | Hashing and password checks      |

### Run all tests:
```bash
docker-compose exec app go test ./... -v
```

---

## 🛠 Development Setup

### Prerequisites:
- Docker + Docker Compose
- Go 1.20+

### Run app locally:
```bash
docker-compose up --build
```

- App available at: `http://localhost:8080`

---

## 🧾 Project Structure

```
.
├── cmd/server             # Entry point
├── internal/auth          # Auth service
├── internal/user          # User logic (in progress)
├── internal/admin         # Admin endpoints
├── common/utils           # JWT, hashing, helpers
├── docker/                # Dockerfile and compose config
├── README.md
└── go.mod
```

---

## ⚠️ Security Notes

- Do not commit `.env` files — these are excluded by `.gitignore`
- Tokens use `HMAC` signing (default). Rotate keys regularly.

---

## 📈 Next Steps

- Add refresh token flow
- Implement email confirmation
- Write integration tests for `/me` and admin endpoints

---

## 📄 License

MIT — use this scaffold for commercial or personal SaaS projects.
