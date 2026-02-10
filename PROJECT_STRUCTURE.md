# Project Structure Guide: CruiseBooking Platform

**Purpose**: This document provides a comprehensive overview of the project structure to facilitate code reviews and reduce token consumption during AI-assisted reviews.

**Last Updated**: 2026-02-10  
**Version**: 1.2.0  
**Phase**: Phase 1 & 2 100% Complete ✓  

---

## Repository Overview

```
open_cruise_sale_system/
├── .github/                    # CI/CD workflows
├── .specify.specify/           # Project constitution & templates
├── specs/                      # Specifications & documentation
│   └── 001-cruise-booking-system/
│       ├── spec.md            # Feature specification (37KB)
│       ├── tasks.md           # Task breakdown (206 tasks)
│       ├── plan/              # Implementation plans
│       ├── contracts/         # API contracts (OpenAPI)
│       ├── docs/              # Developer documentation
│       └── checklists/        # Quality checklists
├── backend/                   # Go monolith (Phase 1 & 2 ✓)
├── frontend-admin/           # Nuxt 4 - Management (Phase 1 & 2 ✓)
├── frontend-web/            # Nuxt 4 - Customer Web (Phase 1 ✓)
├── frontend-mini/          # uni-app - Mini Program (Phase 1 ✓)
└── shared/                # Shared types & utilities (Phase 1 ✓)
```

---

## Backend Structure (`backend/`)

**Technology**: Go 1.26 + Gin v1.11.0 + GORM v2.x

```
backend/
├── cmd/
│   ├── api/                 # Main API server entry point
│   │   └── main.go
│   └── migrate/            # Database migration tool
│       └── main.go
├── internal/
│   ├── config/            # Viper configuration
│   │   └── config.go
│   ├── database/         # GORM database connection
│   │   └── db.go
│   ├── cache/           # Redis client
│   │   └── redis.go
│   ├── storage/        # MinIO S3 client
│   │   └── minio.go
│   ├── messaging/     # NATS JetStream client
│   │   └── nats.go
│   ├── logger/       # Zap structured logging
│   │   └── logger.go
│   ├── domain/      # Business entities
│   │   ├── base.go
│   │   ├── cruise.go
│   │   ├── cabin_type.go
│   │   ├── facility.go
│   │   ├── route.go
│   │   ├── voyage.go
│   │   ├── cabin.go
│   │   ├── inventory.go
│   │   ├── price.go
│   │   ├── order.go
│   │   ├── user.go
│   │   └── ...
│   ├── repository/  # Data access layer
│   │   ├── cruise.go
│   │   ├── cabin_type.go
│   │   └── ...
│   ├── service/    # Business logic layer
│   │   ├── cruise.go
│   │   ├── order.go
│   │   └── ...
│   ├── handler/   # HTTP handlers (controllers)
│   │   ├── cruise.go
│   │   ├── order.go
│   │   └── ...
│   ├── middleware/  # HTTP middleware
│   │   ├── jwt.go
│   │   ├── auth.go
│   │   ├── error.go
│   │   ├── logger.go
│   │   └── cors.go
│   ├── auth/      # Authentication & RBAC
│   │   ├── rbac.go
│   │   └── hash.go
│   ├── validator/ # Request validation
│   │   └── validator.go
│   ├── response/ # API response helpers
│   │   └── response.go
│   └── pagination/ # Pagination utilities
│       └── pagination.go
├── pkg/            # Shared packages
│   └── utils/
├── migrations/    # Database migrations
│   ├── 001_cruise_companies.up.sql
│   ├── 001_cruise_companies.down.sql
│   └── ...
├── tests/        # Test suites
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── docs/        # Swagger/OpenAPI docs
│   └── swagger.json
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── .env.example
└── .air.toml   # Hot reload config
```

### Key Backend Patterns

1. **Clean Architecture**: Domain → Repository → Service → Handler
2. **Dependency Injection**: Services inject repositories, handlers inject services
3. **Middleware Chain**: JWT → RBAC → Logging → Error Recovery
4. **GORM Hooks**: Soft deletes, timestamps, optimistic locking
5. **Test Structure**: Table-driven tests with testify + gomock

---

## Frontend Admin Structure (`frontend-admin/`)

**Technology**: Nuxt 4.3.0 + Vue 3.5.28 + TypeScript 5.9 + Nuxt UI v3

```
frontend-admin/
├── .nuxt/                  # Nuxt build output (generated)
├── .output/               # Production build
├── assets/               # Static assets
│   ├── css/
│   └── images/
├── components/          # Vue components
│   ├── common/         # Shared components
│   │   ├── Loading.vue
│   │   ├── Error.vue
│   │   └── Empty.vue
│   ├── forms/         # Form components
│   │   ├── ImageUpload.vue
│   │   └── TipTapEditor.vue
│   └── navigation/   # Navigation components
│       ├── Sidebar.vue
│       └── Header.vue
├── composables/      # Vue composables
│   ├── useAuth.ts
│   └── useApi.ts
├── layouts/         # Nuxt layouts
│   ├── default.vue
│   ├── admin.vue   # Admin dashboard layout
│   └── auth.vue    # Auth pages layout
├── middleware/     # Route middleware
│   └── auth.ts
├── pages/         # Application pages
│   ├── login.vue
│   ├── index.vue           # Dashboard
│   ├── cruises/
│   │   ├── index.vue      # Cruise list
│   │   └── [id].vue      # Cruise edit/create
│   ├── cabin-types/
│   ├── facilities/
│   ├── routes/
│   ├── voyages/
│   ├── cabins/
│   ├── orders/
│   ├── staffs/
│   └── settings/
├── plugins/      # Nuxt plugins
│   ├── api.ts
│   └── auth.ts
├── stores/      # Pinia stores
│   ├── auth.ts
│   ├── cruise.ts
│   └── ...
├── utils/      # Utility functions
│   └── api.ts
├── types/     # TypeScript types
│   └── index.ts
├── tests/    # Test files
│   ├── unit/
│   └── e2e/
├── nuxt.config.ts
├── tailwind.config.ts
├── tsconfig.json
├── vitest.config.ts
├── playwright.config.ts
├── package.json
├── .env.example
└── Dockerfile
```

### Key Frontend Admin Patterns

1. **File-based Routing**: Pages auto-route based on file structure
2. **Pinia State Management**: Auth store, entity stores
3. **API Client**: Centralized ofetch instance with interceptors
4. **Component Structure**: Atomic design (common → forms → domain)
5. **Auth Guard**: Route middleware for RBAC

---

## Frontend Web Structure (`frontend-web/`)

**Technology**: Nuxt 4.3.0 SSR + Vue 3.5.28 + TypeScript 5.9

```
frontend-web/
├── .nuxt/
├── .output/
├── assets/
│   └── css/
├── components/
│   ├── common/
│   ├── cruise/           # Cruise-specific components
│   │   ├── CruiseCard.vue
│   │   ├── ImageGallery.vue
│   │   └── CabinTypeAccordion.vue
│   ├── booking/         # Booking flow components
│   │   ├── BookingWizard.vue
│   │   ├── SelectVoyage.vue
│   │   ├── SelectCabin.vue
│   │   ├── PassengerForm.vue
│   │   └── OrderConfirm.vue
│   └── layout/
├── composables/
│   ├── useAuth.ts
│   ├── useCruise.ts
│   └── useBooking.ts
├── layouts/
│   ├── default.vue      # Main layout with header/footer
│   └── blank.vue       # Clean layout for checkout
├── middleware/
│   └── auth.ts
├── pages/              # Public pages
│   ├── index.vue      # Home
│   ├── cruises/
│   │   ├── index.vue  # List with filters
│   │   └── [id].vue  # Detail page
│   ├── booking/
│   │   └── index.vue
│   ├── payment/
│   │   ├── [orderId].vue
│   │   └── result.vue
│   ├── orders/
│   │   └── index.vue
│   ├── profile/
│   │   └── index.vue
│   └── about.vue
├── plugins/
├── stores/
├── utils/
├── types/
├── tests/
├── nuxt.config.ts
├── tailwind.config.ts
├── tsconfig.json
├── vitest.config.ts
├── playwright.config.ts
├── package.json
├── .env.example
└── Dockerfile
```

### Key Frontend Web Patterns

1. **SSR for SEO**: All public pages use SSR
2. **Lazy Loading**: Dynamic imports for heavy components
3. **State Hydration**: Pinia state synced from server
4. **Mobile-first**: Responsive design with Tailwind

---

## Mini Program Structure (`frontend-mini/`)

**Technology**: uni-app (Vue 3) + TypeScript

```
frontend-mini/
├── src/
│   ├── components/
│   │   └── CruiseCard.vue
│   ├── pages/           # Mini program pages
│   │   ├── index/
│   │   │   └── index.vue
│   │   ├── cruises/
│   │   │   ├── index.vue
│   │   │   └── detail.vue
│   │   ├── booking/
│   │   ├── payment/
│   │   ├── orders/
│   │   └── profile/
│   ├── static/         # Static assets
│   ├── store/         # Pinia stores
│   ├── utils/
│   │   └── api.ts
│   ├── types/
│   └── App.vue
├── tests/
├── manifest.json      # WeChat Mini Program config
├── pages.json        # Page routes
├── uni.scss
├── vite.config.ts
├── tsconfig.json
├── jest.config.js
├── package.json
└── Dockerfile
```

### Key Mini Program Patterns

1. **WeChat SDK**: Native WeChat API integration
2. **Conditional Compilation**: Platform-specific code
3. **Subpackages**: Code splitting for performance

---

## Shared Types (`shared/`)

**Purpose**: TypeScript types shared across all frontend projects

```
shared/
└── types/
    ├── package.json
    ├── tsconfig.json
    ├── src/
    │   ├── index.ts
    │   ├── cruise.ts
    │   ├── cabin.ts
    │   ├── order.ts
    │   ├── user.ts
    │   └── api.ts
    └── dist/
```

---

## Configuration Files

### Root Level

| File | Purpose |
|------|---------|
| `.gitignore` | Git ignore patterns |
| `docker-compose.yml` | Local infrastructure (Postgres, Redis, MinIO, NATS) |
| `README.md` | Project documentation |
| `Makefile` | Common development commands |

### Backend Config

| File | Purpose |
|------|---------|
| `go.mod` | Go module definition |
| `.env.example` | Environment variables template |
| `Dockerfile` | Multi-stage build for production |
| `.air.toml` | Hot reload configuration |
| `golangci.yml` | Linter configuration |

### Frontend Config

| File | Purpose |
|------|---------|
| `nuxt.config.ts` | Nuxt configuration |
| `tailwind.config.ts` | Tailwind CSS configuration |
| `vitest.config.ts` | Unit test configuration |
| `playwright.config.ts` | E2E test configuration |
| `tsconfig.json` | TypeScript configuration |
| `eslint.config.js` | ESLint configuration |
| `prettier.config.js` | Prettier formatting |
| `package.json` | NPM dependencies & scripts |

---

## Naming Conventions

### Backend (Go)

```go
// Files: snake_case.go
cruise_repository.go
cruise_service.go
cruise_handler.go

// Types: PascalCase
type CruiseService struct {}
type CruiseRepository struct {}

// Interfaces: PascalCase with -er suffix
type CruiseServicer interface {}
type CruiseRepositoryer interface {}

// Functions: PascalCase (exported) / camelCase (internal)
func (s *CruiseService) GetByID(id uuid.UUID) (*Cruise, error)
func validateCruise(c *Cruise) error

// Variables: camelCase
var cruiseRepository *CruiseRepository

// Constants: UPPER_SNAKE_CASE
const MAX_PAGE_SIZE = 100
```

### Frontend (Vue/TypeScript)

```typescript
// Files: PascalCase.vue for components
CruiseCard.vue
BookingWizard.vue

// Files: camelCase.ts for utilities
useAuth.ts
apiClient.ts

// Components: PascalCase
export default defineComponent({
  name: 'CruiseCard'
})

// Composables: camelCase with use- prefix
export function useCruise() {}
export function useAuth() {}

// Stores: camelCase with Store suffix
export const useAuthStore = defineStore('auth', {})

// Types: PascalCase
interface Cruise {}
type OrderStatus = 'pending' | 'paid' | 'confirmed'

// Variables: camelCase
const cruiseList = ref<Cruise[]>([])

// Constants: UPPER_SNAKE_CASE
const API_BASE_URL = '/api/v1'
```

---

## Import Patterns

### Backend

```go
// Standard library first
import (
    "context"
    "time"
    
    // Third-party packages
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    // Internal packages (project-specific)
    "backend/internal/domain"
    "backend/internal/repository"
    "backend/internal/service"
)
```

### Frontend

```typescript
// Vue/Nuxt imports
import { defineComponent, ref, computed } from 'vue'
import { useRoute, useRouter } from 'nuxt/app'

// Third-party
import { useFetch } from '#app'

// Local imports
import { useAuthStore } from '~/stores/auth'
import type { Cruise } from '~/types'

// Component imports (auto-imported by Nuxt, but can be explicit)
import CruiseCard from '~/components/cruise/CruiseCard.vue'
```

---

## Testing Structure

### Backend Tests

```
backend/
├── internal/
│   ├── service/
│   │   ├── cruise_test.go           # Unit tests
│   │   └── cruise_integration_test.go
│   └── handler/
│       └── cruise_test.go
└── tests/
    ├── e2e/
    │   └── booking_flow_test.go
    └── fixtures/
        └── cruises.yml
```

**Test Command**: `go test ./... -coverprofile=coverage.out`

### Frontend Tests

```
frontend-admin/
├── components/
│   └── CruiseCard.spec.ts    # Vitest unit tests
└── tests/
    └── e2e/
        └── cruise-management.spec.ts  # Playwright E2E
```

**Test Commands**:
- Unit: `npm run test:unit`
- E2E: `npm run test:e2e`
- Coverage: `npm run test:coverage`

---

## CI/CD Pipeline

```
.github/
└── workflows/
    ├── ci.yml           # Pull request validation
    └── deploy.yml       # Production deployment
```

### CI Steps

1. **Lint**: golangci-lint, ESLint, Prettier
2. **Test**: Unit + Integration (100% coverage gate)
3. **Build**: Docker images
4. **Security**: Dependency audit, secret scan
5. **Deploy**: Kubernetes rollout (staging/production)

---

## Database Migrations

### Migration Naming

```
migrations/
├── 001_cruise_companies.up.sql
├── 001_cruise_companies.down.sql
├── 002_cruises.up.sql
├── 002_cruises.down.sql
└── ...
```

### Migration Commands

```bash
# Create new migration
go run cmd/migrate/main.go create add_user_table

# Run migrations
go run cmd/migrate/main.go up

# Rollback
go run cmd/migrate/main.go down 1

# Status
go run cmd/migrate/main.go status
```

---

## API Versioning

**Base URL**: `/api/v1`

**Structure**:
```
/api/v1/
├── /auth          # Authentication
├── /cruises       # Cruise management
├── /cabin-types   # Cabin categories
├── /facilities    # Ship amenities
├── /routes        # Itineraries
├── /voyages       # Sailing dates
├── /cabins        # Inventory & booking
├── /orders        # Booking workflow
├── /payments      # Payment processing
├── /users         # Customer profiles
├── /staffs        # Admin management
└── /notifications # Messages
```

---

## Environment Variables

### Backend (.env)

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=cruisebooking
DB_PASSWORD=secret
DB_NAME=cruisebooking

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin

# WeChat Pay
WECHAT_MCH_ID=xxx
WECHAT_API_V3_KEY=xxx
```

### Frontend (.env)

```bash
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
NUXT_PUBLIC_STORAGE_BASE=http://localhost:9000/cruisebooking
```

---

## Quick Reference

### Start Development

```bash
# Infrastructure
docker-compose up -d

# Backend
cd backend
go run cmd/api/main.go

# Frontend Admin
cd frontend-admin
npm run dev

# Frontend Web
cd frontend-web
npm run dev
```

### Common Commands

```bash
# Backend tests
go test ./... -v

# Frontend tests
npm run test:unit
npm run test:e2e

# Lint
golangci-lint run
npm run lint

# Swagger generation
cd cmd/api && swag init
```

---

## Architecture Decisions

1. **Monolith First**: Single Go backend, split later if needed
2. **Domain-driven Design**: Clear domain boundaries (Cruise, Order, User)
3. **CQRS-lite**: Separate read/write models where beneficial
4. **Event-driven**: NATS for async operations (order timeout, notifications)
5. **API-first**: OpenAPI specification drives implementation

---

## Implementation Progress

### Phase 1: Setup ✓ (20/20)
- [x] Go project structure with go.mod
- [x] Gin + GORM configuration
- [x] PostgreSQL, Redis, MinIO, NATS clients
- [x] Viper configuration management
- [x] Zap structured logging
- [x] Nuxt 4.3.0 frontend projects (admin, web)
- [x] uni-app mini program project
- [x] docker-compose.yml infrastructure
- [x] GitHub Actions CI workflow
- [x] Environment configuration files
- [x] Shared TypeScript types package
- [x] Swagger generation config

### Phase 2: Foundational Infrastructure ✓ (16/16)
- [x] JWT middleware with token generation/validation
- [x] Casbin RBAC with role hierarchy
- [x] Role definitions (Super Admin, Operations, Finance, CS)
- [x] Password hashing utilities (bcrypt)
- [x] Pagination utilities
- [x] API response wrapper (standard JSON format)
- [x] Error handling middleware
- [x] Request logging middleware
- [x] CORS middleware
- [x] Request validation utilities
- [x] Pinia auth store
- [x] API client with ofetch
- [x] Auth handlers (login, refresh, logout, me)
- [x] Layout components (Header, Sidebar, Footer)
- [x] Common UI components (Loading, Error, Empty)
- [x] Route guards for authentication

### Next: Phase 3 - Cruise Browsing
- [ ] Database migrations (cruises, cabin_types, facilities)
- [ ] Domain models
- [ ] Repository layer
- [ ] Service layer
- [ ] HTTP handlers
- [ ] Frontend pages and components

---

## Contact & Resources

- **Constitution**: `.specify.specify/memory/constitution.md`
- **Specification**: `specs/001-cruise-booking-system/spec.md`
- **Tasks**: `specs/001-cruise-booking-system/tasks.md`
- **API Docs**: `specs/001-cruise-booking-system/contracts/openapi.yaml`

---

**End of Project Structure Guide**

*Use this document as a quick reference during code reviews and development.*
