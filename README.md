# 📦 Lamina SaaS Backend

[![Go Version](https://img.shields.io/badge/Go-1.24.1-blue?logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue?logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Containerized-blue?logo=docker)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green)](https://opensource.org/licenses/MIT)

**Flight Crew Scheduling and Transport Management SaaS platform**, built for long-term scalability, maintainability, and low operational cost.

---

## 📑 Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Setup Instructions](#setup-instructions)
- [Authentication Strategy](#authentication-strategy)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Best Practices Followed](#best-practices-followed)

---

## ✨ Features
- 🚀 Fast and minimal Go backend
- 🔐 Secure authentication using JWT tokens
- 🛡️ Admin-only user creation (no public signup)
- 🗃️ PostgreSQL database (Dockerized)
- 🛠️ Structured project architecture: modular, clean, scalable
- 🧪 Manual and unit testing enforced (Test-Driven Development)
- 📈 Ready for CI/CD pipeline integration

---

## 🏛️ Architecture
- **Backend**: Golang (Gin framework)
- **Database**: PostgreSQL (inside Docker container)
- **Auth**: JWT (JSON Web Tokens) with custom claims (userID, email, role)
- **Deployment**: Docker & Docker Compose for local development and production readiness

```
Client → Nginx (future) → API Gateway (future) → Lamina API Server → PostgreSQL DB
```

---

## 🛠 Tech Stack
| Layer | Technology |
|------|-------------|
| Backend | Go 1.24.1 |
| Framework | Gin |
| Auth | JWT |
| Database | PostgreSQL 16 |
| DevOps | Docker, Docker Compose |
| Testing | Go built-in `testing` package |

---

## 📂 Project Structure

```
lamina/
├── cmd/
│   └── server/             # Entry point (main.go)
├── common/
│   ├── database/            # PostgreSQL connection management
│   └── utils/               # JWT, password utilities
├── config/                  # Environment config management
├── internal/
│   ├── auth/                # Authentication business logic
│   ├── user/                # User management logic
│   └── admin/               # Admin-only features (create users)
├── docker/
│   ├── docker-compose.yml   # Docker orchestration
│   └── Dockerfile           # Build Go app
├── .gitignore
├── .env.example
├── go.mod
├── go.sum
└── README.md                # Project documentation
```

---

## 🚀 Setup Instructions

### 1. Clone the project

```bash
git clone https://github.com/nomenarkt/lamina.git
cd lamina/docker
```

### 2. Configure environment variables

```bash
cp ../.env.example ../.env
```
- Fill `.env` with your local secrets (PostgreSQL credentials, JWT secret, etc.).

### 3. Launch Docker containers

```bash
docker-compose up --build
```

### 4. Initialize the database manually (if needed)

```bash
docker exec -it docker_db_1 psql -U postgres -d saasdb
```

Create `users` table if missing.

---

## 🔐 Authentication Strategy

| Endpoint             | Public/Protected | Description |
|----------------------|------------------|-------------|
| `/auth/login`         | Public | Users can login if already registered |
| `/auth/signup`        | ❌ Not exposed | No public signup allowed |
| `/admin/create-user`  | Protected (JWT + Admin role) | Admin creates users manually |

✅ JWT Tokens contain `userID`, `email`, and `role` claims for secured access.

---

## 🌐 API Endpoints Overview

| Method | Endpoint                | Access |
|--------|--------------------------|--------|
| POST   | `/api/v1/auth/login`      | Public |
| POST   | `/api/v1/admin/create-user` | Admin-only (JWT) |
| GET    | `/api/v1/user/me`          | Authenticated |
| GET    | `/api/v1/user/`            | Admin |

---

## 🧪 Testing

Unit tests (e.g., JWT generation and parsing) are located under:

```
/common/utils/jwt_test.go
```

Run tests locally:

```bash
go test ./common/utils/
```

---

## 📏 Best Practices Followed

- [x] Software Development Lifecycle (SDLC) applied properly
- [x] Test-Driven Development (TDD) enforced
- [x] Modular architecture (Clean Architecture Principles)
- [x] Secure role-based access control
- [x] Environment-specific configs (.env separation)
- [x] Full Dockerized environment (PostgreSQL + App)

---

# ✅ Status
Lamina backend is **production-ready** for internal deployment. Future improvements include:
- API Gateway
- Admin dashboard (frontend)
- Metrics/Monitoring setup (Prometheus + Grafana)

---

> Built with passion for quality software. ✈️🚛
