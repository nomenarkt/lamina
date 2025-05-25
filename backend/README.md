# Lamina Backend

🧱 Backend services for Lamina — an aviation-grade crew and flight management system. Built with Go, PostgreSQL, and Gin, using modular clean architecture and test-driven development.

---

## 🚀 Features

- Secure JWT authentication
- Admin user invites and email confirmation
- Role-based access control (RBAC)
- Token expiration handling and confirmation redirect flows
- Full Docker + migrate integration
- Live reload via Air during development
- Unit tests with coverage enforcement

---

## 🔧 Development Setup

### Prerequisites

- Docker & Docker Compose
- Go 1.24+

### Start Development Environment

```bash
make dev-up
This starts:
-PostgreSQL (via docker-compose)
-Backend app with live reload (via air)
-Auto-run database migrations
Useful Commands
| Command          | Description                                  |
| ---------------- | -------------------------------------------- |
| `make dev-up`    | Start DB, run migrations, boot backend (air) |
| `make app-logs`  | Tail logs from backend container             |
| `make dev-reset` | Full clean, rebuild images, restart services |
| `make test`      | Run all backend unit tests with coverage     |
| `make lint`      | Run Go linters via golangci-lint             |
| `make migrate`   | Apply latest DB migrations                   |
| `make rebuild`   | Force rebuild backend and migrate images     |
| `make down`      | Shutdown and cleanup containers/volumes      |

🔐 Auth Flow Overview
| Endpoint                     | Description                                 |
| ---------------------------- | ------------------------------------------- |
| `POST /auth/signup`          | Register internal user                      |
| `POST /auth/login`           | Login and return access/refresh JWTs        |
| `GET /auth/confirm/:token`   | Email token confirmation redirect           |
| `POST /auth/complete-invite` | Set password after admin invite             |
| `POST /admin/create-user`    | Admin-only: invite user (internal/external) |
| `GET /user/me`               | Return authenticated user details           |
Confirmation Redirect Logic
| Scenario                      | Redirect Target                           |
| ----------------------------- | ----------------------------------------- |
| First-time confirmation       | `/email-confirmed`                        |
| Reused token (already active) | `/confirm-error?reason=already-confirmed` |
| Token expired (24h+)          | `/confirm-error?reason=expired`           |
| Invalid/malformed token       | `/confirm-error?reason=invalid`           |

🧪 Testing
Unit tests are enforced with coverage for critical flows:
make test
Key test files:
-internal/auth/service_test.go
-internal/tests/auth_test.go
-common/utils/password_test.go

🧠 Project Structure
backend/
├── cmd/server               # Gin HTTP entrypoint
├── common/utils             # Shared password, token, and helper logic
├── internal/
│   ├── auth                 # Signup, login, email tokens
│   ├── admin                # Admin invite flow
│   ├── user                 # Profile & info
│   ├── crew                 # Crew assignments
│   └── middleware           # JWT middleware
├── migrations/              # Golang Migrate SQL scripts
├── docker/                  # App + migrate Dockerfiles
├── .air.toml                # Air config
├── Makefile                 # Backend scoped makefile

🐳 Docker Stack
-App: live reload using cosmtrek/air
-PostgreSQL: v16
-Migrations: executed via golang-migrate
make dev-up      # Full launch
make migrate     # Rerun migrations only
make down        # Cleanup

🛡 Security Highlights
-Passwords hashed with bcrypt
-JWTs signed with .env secrets (not committed)
-Confirmation tokens expire in 24h
-Role-based endpoint protection
-Admin creation is email + duration controlled

📄 License
MIT — freely usable for commercial or internal SaaS backend development.