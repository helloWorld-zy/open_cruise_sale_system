# Quick Start Guide: CruiseBooking Platform

**Version**: 1.0.0  
**Last Updated**: 2026-02-10  

---

## Prerequisites

- **Go**: 1.26 or later
- **Node.js**: 20.x LTS
- **Docker & Docker Compose**: Latest stable
- **Git**: Latest
- **HBuilderX**: For mini program development (uni-app)

---

## 1. Repository Setup

```bash
# Clone repository
git clone <repository-url>
cd open_cruise_sale_system

# Verify you're on main branch
git branch --show-current
# Should output: main
```

---

## 2. Infrastructure Setup (Docker)

### Start Infrastructure Services

```bash
# Navigate to backend
cd backend

# Copy environment file
cp .env.example .env

# Start all infrastructure services
docker-compose up -d

# Verify services are running
docker-compose ps
```

### Services Started

| Service | Port | Purpose |
|---------|------|---------|
| PostgreSQL | 5432 | Main database |
| Redis | 6379 | Cache & session |
| Meilisearch | 7700 | Search engine |
| MinIO | 9000/9001 | Object storage |
| NATS | 4222 | Message queue |

### Environment Variables (.env)

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=cruisebooking
DB_PASSWORD=your_secure_password
DB_NAME=cruisebooking

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Meilisearch
MEILI_HOST=http://localhost:7700
MEILI_API_KEY=your_meili_key

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=cruisebooking

# NATS
NATS_URL=nats://localhost:4222

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE_HOURS=24

# WeChat Pay
WECHAT_MCH_ID=your_merchant_id
WECHAT_API_V3_KEY=your_api_v3_key
WECHAT_APP_ID=your_app_id
WECHAT_CERT_PATH=./certs/apiclient_cert.pem
WECHAT_KEY_PATH=./certs/apiclient_key.pem
```

---

## 3. Backend Setup (Go)

### Install Dependencies

```bash
cd backend

# Download Go modules
go mod download

# Install development tools
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Database Migration

```bash
# Run migrations
go run cmd/migrate/main.go up

# Check migration status
go run cmd/migrate/main.go status
```

### Run Development Server

```bash
# Run with hot reload (using air)
air

# Or run directly
go run cmd/api/main.go

# Server will start on http://localhost:8080
```

### Verify Backend

```bash
# Health check
curl http://localhost:8080/health

# API documentation (Swagger)
open http://localhost:8080/swagger/index.html
```

---

## 4. Management Frontend Setup (Nuxt 4)

### Install Dependencies

```bash
cd frontend-admin

# Install dependencies
npm install

# Or use pnpm (recommended)
pnpm install
```

### Environment Configuration

```bash
# Create .env file
cp .env.example .env
```

```bash
# .env
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
NUXT_PUBLIC_STORAGE_BASE=http://localhost:9000/cruisebooking
```

### Run Development Server

```bash
# Start dev server
npm run dev

# Server will start on http://localhost:3000
```

---

## 5. Customer Web Frontend Setup (Nuxt 4 SSR)

### Install Dependencies

```bash
cd frontend-web

npm install
```

### Environment Configuration

```bash
cp .env.example .env
```

```bash
# .env
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
NUXT_PUBLIC_STORAGE_BASE=http://localhost:9000/cruisebooking
```

### Run Development Server

```bash
npm run dev

# Server will start on http://localhost:3001
```

---

## 6. Mini Program Setup (uni-app)

### Install Dependencies

```bash
cd frontend-mini

npm install
```

### HBuilderX Setup

1. Open HBuilderX
2. File → Open Directory → Select `frontend-mini`
3. Configure `manifest.json`:
   - Set WeChat App ID
   - Configure App name and description

### Run Development

```bash
# In HBuilderX, click:
# Run → Run to Mini Program Simulator → WeChat Developer Tools
```

### Configure WeChat Developer Tools

1. Open WeChat Developer Tools
2. Import project from `frontend-mini/dist/dev/mp-weixin`
3. Set App ID in project.config.json
4. Enable "Do not verify domain..." in settings

---

## 7. Testing

### Backend Tests

```bash
cd backend

# Run all tests
go test ./... -v

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/service/... -v
```

### Frontend Tests

```bash
# Management Frontend
cd frontend-admin

# Run Vitest unit tests
npm run test:unit

# Run Playwright E2E tests
npm run test:e2e

# Coverage report
npm run test:coverage

# Customer Web Frontend
cd frontend-web
npm run test:unit
npm run test:e2e
```

### Mini Program Tests

```bash
cd frontend-mini

# Run Jest tests
npm run test:unit

# WeChat DevTools automated testing
# Use WeChat DevTools → Test
```

---

## 8. Linting & Code Quality

### Backend

```bash
cd backend

# Run linter
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix

# Generate Swagger docs
cd cmd/api && swag init
```

### Frontend

```bash
cd frontend-admin

# Run ESLint
npm run lint

# Fix auto-fixable issues
npm run lint:fix

# Format with Prettier
npm run format
```

---

## 9. CI/CD Pipeline

### GitHub Actions Workflow

The repository includes GitHub Actions workflows in `.github/workflows/`:

1. **CI Pipeline** (`ci.yml`)
   - Runs on every PR
   - Lint, test, build all components
   - Coverage threshold: 100%

2. **Deploy Pipeline** (`deploy.yml`)
   - Runs on main branch merge
   - Builds Docker images
   - Deploys to Kubernetes cluster

### Running CI Locally

```bash
# Use act to run GitHub Actions locally
act -j build-and-test
```

---

## 10. Common Development Tasks

### Add New Database Migration

```bash
cd backend

# Create new migration
go run cmd/migrate/main.go create add_user_preferences

# Edit generated files in migrations/
# Up migration: migrations/xxx_add_user_preferences.up.sql
# Down migration: migrations/xxx_add_user_preferences.down.sql
```

### Add New API Endpoint

1. Define model in `internal/domain/`
2. Create repository in `internal/repository/`
3. Implement service in `internal/service/`
4. Add handler in `internal/handler/`
5. Register route in `cmd/api/routes.go`
6. Write tests
7. Generate Swagger docs: `swag init`

### Add New Frontend Component

```bash
cd frontend-admin

# Create component
mkdir -p components/Cruise

# Write component
# components/Cruise/CruiseCard.vue

# Write tests
# components/Cruise/CruiseCard.spec.ts

# Run component tests
npm run test:unit CruiseCard
```

---

## 11. Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Reset database (WARNING: destroys data)
docker-compose down -v
docker-compose up -d postgres
go run cmd/migrate/main.go up
```

### Port Conflicts

```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Go Module Issues

```bash
cd backend

# Clean module cache
go clean -modcache

# Re-download
go mod download

# Tidy modules
go mod tidy
```

### Frontend Build Issues

```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear Nuxt cache
rm -rf .nuxt
npm run dev
```

---

## 12. Production Deployment

### Build Docker Images

```bash
# Backend
docker build -t cruisebooking/backend:latest -f backend/Dockerfile .

# Frontend Admin
docker build -t cruisebooking/frontend-admin:latest -f frontend-admin/Dockerfile .

# Frontend Web
docker build -t cruisebooking/frontend-web:latest -f frontend-web/Dockerfile .
```

### Kubernetes Deployment

```bash
# Apply manifests
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend-admin.yaml
kubectl apply -f k8s/frontend-web.yaml
kubectl apply -f k8s/ingress.yaml
```

---

## 13. Useful Commands

### Database

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U cruisebooking -d cruisebooking

# Backup database
docker-compose exec postgres pg_dump -U cruisebooking cruisebooking > backup.sql

# Restore database
docker-compose exec -T postgres psql -U cruisebooking -d cruisebooking < backup.sql
```

### Redis

```bash
# Connect to Redis CLI
docker-compose exec redis redis-cli

# Flush all data (WARNING)
docker-compose exec redis redis-cli FLUSHALL
```

### MinIO

```bash
# Create bucket
mc alias set local http://localhost:9000 minioadmin minioadmin
mc mb local/cruisebooking

# List objects
mc ls local/cruisebooking
```

---

## 14. Contributing

### Branch Strategy

- `main` - Production-ready code
- `feature/*` - Feature branches
- `bugfix/*` - Bug fix branches

### Commit Convention

```
feat: add new feature
fix: fix bug
docs: update documentation
style: formatting changes
refactor: code refactoring
test: add tests
chore: maintenance tasks
```

### Pull Request Process

1. Create feature branch from main
2. Make changes with tests
3. Ensure 100% test coverage
4. Run linting and fix issues
5. Create PR with description
6. Request review from maintainers
7. Merge after approval and CI pass

---

## 15. Resources

### Documentation
- [Constitution](../../.specify.specify/memory/constitution.md) - Project constitution
- [Specification](../spec.md) - Feature specification
- [Implementation Plan](./impl-plan.md) - Development plan
- [Data Model](./data-model.md) - Database schema

### External Links
- [Go Documentation](https://go.dev/doc/)
- [Nuxt 4 Documentation](https://nuxt.com/docs)
- [Vue 3 Documentation](https://vuejs.org/guide/)
- [uni-app Documentation](https://uniapp.dcloud.net.cn/)
- [GORM Documentation](https://gorm.io/docs/)
- [Gin Documentation](https://gin-gonic.com/docs/)

---

**Happy Coding!** 🚢

For questions or issues, refer to the Constitution or contact the maintainers.
