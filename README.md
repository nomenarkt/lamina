📚 Lamina - Flight & Transport SaaS Platform
Lamina is a modular, scalable SaaS backend designed for aviation and logistics companies to manage flight crew scheduling and transport operations efficiently.

✈️🚛 Features
Flight Crew Scheduling Management
Transport Jobs and Driver Assignments
Multi-tenant SaaS Architecture (Company Accounts)
Secure User Authentication (JWT + Bcrypt)
PostgreSQL 16 Database (Optimized for SaaS)
Dockerized Development and Deployment
Modular, Clean Golang Backend (Gin Framework)
Extensible Microservices-Ready Design
Professional CI/CD Ready (Future GitHub Actions Integration)

🏗️ Project Architecture
Layer | Technology
Backend | Go 1.24 (Gin Web Framework)
Database | PostgreSQL 16
Authentication | JWT Tokens, Bcrypt Passwords
DevOps | Docker, Docker Compose
Hosting (recommended) | AWS Lightsail (MVP) or AWS ECS (Scaling)
CI/CD (future) | GitHub Actions
✅ Clean Monolith Architecture (Microservices Ready Later)

✅ Multi-module setup (/internal/flight, /internal/transport, /internal/auth, etc.)

🧰 Tech Stack
Component | Technology
Language | Go 1.24
Framework | Gin
Database | PostgreSQL
Authentication | JWT + Bcrypt
Dev Environment | Docker, Docker Compose
Package Management | Go Modules

🚀 Getting Started
1. Clone the Repository
git clone https://github.com/nomenarkt/lamina.git
cd lamina
2. Setup Environment Variables
-Create your environment file:
cp .env.example .env
-Edit .env to configure:
PORT=8080
DATABASE_URL=postgres://postgres:postgres@db:5432/saasdb?sslmode=disable
JWT_SECRET=your_super_secret_key
3. Build and Run Locally
Use Docker Compose to start Postgres and the Go backend:
docker-compose up --build
Access API at:
http://localhost:8080/api/v1
4. API Endpoints Available
Method	Endpoint	Description
POST	/api/v1/auth/signup	Register a new user
POST	/api/v1/auth/login	Authenticate and receive access token
(More endpoints coming as we build: user management, flight scheduling, transport jobs.)

🧪 Testing the API
You can use tools like:
Postman (recommended)
cURL
Example (using cURL):
curl -X POST http://localhost:8080/api/v1/auth/signup \
-H "Content-Type: application/json" \
-d '{"email": "test@example.com", "password": "strongpassword"}'

🔥 Roadmap
Phase | Features
MVP	| Authentication, Basic User Management, Basic Flight Scheduling, Basic Transport Management
Post-MVP | Notifications System (Email/SMS), Advanced Reporting, Mobile App API
Scale | Multi-database support, Event Queues, Kubernetes ready deployments

📦 Folder Structure
plaintext
Copy
lamina/
├── cmd/server         # Main entry point (main.go)
├── config/             # Config loader
├── internal/           # Main business modules (auth, user, flight, transport)
├── common/             # Shared utilities (db, jwt, password)
├── docker/             # Docker and docker-compose files
├── tests/              # Unit and integration tests
├── .env.example        # Environment variable template
├── .gitignore          # Git ignored files
├── README.md           # Project documentation
├── go.mod, go.sum      # Go module dependencies

✍️ Contributing
Pull requests are welcome!
For major changes, please open an issue first to discuss what you would like to change.
✅ Always make sure to run:
go fmt ./...
before pushing your code!

📜 License
MIT License

🧠 Important Notes
Environment variables must be set properly to run backend.
Always use go mod tidy after adding any new package.
Make sure PostgreSQL is running locally (Docker Compose starts it automatically).

🏆 Author
Developed by Nomenarkt
In collaboration with [Lamina Project CTO Assistant] (powered by GPT technologies 🚀).

❤️ Final Line
Building SaaS the right way — clean, scalable, production-ready from Day 1.