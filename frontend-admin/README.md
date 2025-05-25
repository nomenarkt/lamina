# ✈️ Lamina Frontend (Admin Portal)

This is the React-based frontend for Lamina — a scalable SaaS platform for aviation logistics and crew scheduling.

It’s built for Madagascar Airlines staff and admin users to access protected resources and onboard securely.

---

## 🧱 Tech Stack

| Area              | Tech                            |
|-------------------|---------------------------------|
| Framework         | Next.js App Router (React 18)   |
| Styling           | Tailwind CSS (token-based)      |
| Language          | TypeScript                      |
| Testing           | Jest + React Testing Library    |
| Forms & Auth      | LocalStorage tokens (MVP phase) |
| Build & Deploy    | Vercel-ready                    |

---

## 📦 Key Features

- ✅ Email confirmation flows with 24h expiry
- ✅ JWT-based login with inline error feedback
- ✅ Admin and self-service signup (internal-only domain)
- ✅ ARIA-accessible forms and alerts
- ✅ Error handling and redirect logic from backend token flows
- ✅ Unit-tested Card components and confirmation screens

---

## 🔧 Local Setup

```bash
cd frontend-admin
npm install
cp .env.local.example .env.local
Update .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080

🚀 Development Commands
| Action                          | Command              |
| ------------------------------- | -------------------- |
| Run dev server (localhost:3000) | `npm run dev`        |
| Run tests                       | `make frontend-test` |
| Run lint check                  | `make frontend-lint` |
| Run tests directly              | `npm test`           |

📁 Directory Structure
frontend-admin/
├── src/
│   ├── app/                    # App Router pages
│   ├── components/auth/        # LoginCard, SignupCard, Layout
│   ├── app/confirm-error/      # Dynamic redirect + messaging
│   ├── lib/api/                # API calls (e.g. signup, login)
│   └── tailwind.config.js      # Theme tokens (green, gold)
├── public/
│   └── logo.webp               # Madagascar Airlines branding
├── __tests__/                  # Jest unit tests (co-located soon)
└── README.md                   # You're here

🧪 Test Coverage
make frontend-test
-Includes tests for:
-Field validation
-Email domain enforcement
-Backend error display
-Confirmation screen behavior
-Dynamic error routing

🌐 Planned Features
 Replace localStorage with HttpOnly secure cookies
 Multi-language support (English + French)
 Resend confirmation link from /check-email
 Session-aware route guards (/dashboard)
 End-to-end tests (Cypress or Playwright)

🧠 Frontend Philosophy
Lamina’s frontend is:
-Token-driven for theme and brand consistency
-Built with ARIA-first accessibility patterns
-Designed for testability and DX
-Scalable across teams and modules