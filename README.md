# ✈️ Lamina

**Lamina** is a modular, production-ready SaaS platform for aviation and logistics teams to manage crew operations and secure user workflows — built with clean architecture, test-driven development, and scalable cloud-native practices.

---

## 📂 Module READMEs

- [`backend/README.md`](./backend/README.md) — Go services, PostgreSQL, auth flows, Makefile CLI
- [`frontend-admin/README.md`](./frontend-admin/README.md) — Next.js-based admin panel with auth and redirect UX

---

## 📦 System Architecture

```text
+-------------+       +------------------+
|  Frontend   | <---> | Gin HTTP Server  |
+-------------+       +------------------+
                            |
                            v
                   +---------------------+
                   | Middleware (JWT)    |
                   +---------------------+
                            |
                            v
                  +----------------------+
                  | Domain Services      |
                  | (Admin/User/Auth)    |
                  +----------------------+
                            |
                            v
                   +---------------------+
                   | PostgreSQL (sqlx)   |
                   +---------------------+


🛠️ Development Quick Start
make dev-up       # Start DB, run migrations, launch backend with Air
make frontend-dev # Start frontend Vite dev server
make app-logs     # Tail backend logs
make dev-reset    # Full reset (rebuild containers, migrations)


🧪 CI, TDD & Engineering Practices
-Full unit tests with mocks (make test)
-Code linting via golangci-lint
-TDD-first auth and user flows
-Separation of concerns via domain packages
-Frontend UX designed for token-based redirects and email flows


🔐 Auth Flow Summary
-Email signup with domain restrictions
-Admin-invite flow for internal/external users
-Confirmation token with expiration (24h)
-Frontend redirected to:
 -/email-confirmed (success)
 -/confirm-error?reason=expired|invalid|already-confirmed (errors)


 🧱 Project Layout
 lamina/
├── backend/                # Go backend and SQL migrations
├── frontend-admin/         # React frontend built with Next.js and Tailwind CSS
├── Makefile                # Unified CLI to manage backend + frontend
└── README.md               # Root overview (you are here)


🧠 Philosophy
Lamina is built on principles from Clean Architecture, Refactoring, and Software Engineering at Google. Every feature is tested, every decision justified, and every component designed for clarity and long-term maintainability.

Own your architecture. Ship with confidence.