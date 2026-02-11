# Project Structure Guide: CruiseBooking Platform

**Purpose**: This document provides a comprehensive overview of the project structure to facilitate code reviews and reduce token consumption during AI-assisted reviews.

**Last Updated**: 2026-02-11
**Version**: 1.4.0
**Phase**: Phase 1, 2, 3 & 4 100% Complete ✓  

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
├── backend/                   # Go monolith (Phase 1, 2, 3 & 4 ✓)
├── frontend-admin/           # Nuxt 4 - Management (Phase 1, 2, 4 ✓)
├── frontend-web/            # Nuxt 4 - Customer Web (Phase 1, 3 ✓)
├── frontend-mini/          # uni-app - Mini Program (Phase 1, 3 ✓)
└── shared/                # Shared types & utilities (Phase 1 ✓)
```

---

## Backend Structure (`backend/`)

**Technology**: Go 1.26 + Gin v1.11.0 + GORM v2.x

```
backend/
├── cmd/
│   ├── api/                 # Main API server entry point ✓ Phase 3 & 4
│   │   ├── main.go
│   │   ├── routes.go        # ✓ Public API routes
│   │   ├── admin_routes.go  # ✓ Admin API routes with RBAC (Phase 4)
│   │   └── swagger.go       # Swagger documentation
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
│   ├── domain/      # Business entities ✓ Phase 3 Complete
│   │   ├── base.go                 # Base model with soft delete
│   │   ├── cruise_company.go       # Cruise company entity
│   │   ├── cruise.go              # Cruise ship entity
│   │   ├── cabin_type.go          # Cabin type entity
│   │   ├── facility.go            # Facility & category entities
│   │   ├── route.go               # (Phase 5)
│   │   ├── voyage.go              # (Phase 5)
│   │   ├── cabin.go               # (Phase 5)
│   │   ├── inventory.go           # (Phase 5)
│   │   ├── price.go               # (Phase 5)
│   │   ├── order.go               # (Phase 5-6)
│   │   ├── user.go                # (Phase 7)
│   │   └── ...
│   ├── repository/  # Data access layer ✓ Phase 3 Complete
│   │   ├── cruise.go              # CruiseRepository
│   │   ├── cabin_type.go          # CabinTypeRepository
│   │   ├── facility.go            # FacilityRepository + FacilityCategoryRepository
│   │   └── ...
│   ├── service/    # Business logic layer ✓ Phase 3 & 4 Complete
│   │   ├── cruise.go              # CruiseService ✓
│   │   ├── cruise_test.go         # Unit tests ✓
│   │   ├── cabin_type.go          # CabinTypeService ✓
│   │   ├── facility.go            # FacilityService + FacilityCategoryService ✓
│   │   ├── storage.go             # ✓ StorageService (Phase 4)
│   │   └── ...
│   ├── handler/   # HTTP handlers (controllers) ✓ Phase 3 & 4 Complete
│   │   ├── auth.go                # Auth handlers (Phase 2)
│   │   ├── cruise.go              # Public cruise handlers ✓
│   │   ├── cruise_test.go         # Integration tests ✓
│   │   ├── cabin_type.go          # Public cabin type handlers ✓
│   │   ├── facility.go            # Public facility handlers ✓
│   │   ├── admin_cruise.go        # ✓ Admin cruise handlers (Phase 4)
│   │   ├── admin_cabin_type.go    # ✓ Admin cabin type handlers (Phase 4)
│   │   ├── admin_facility.go      # ✓ Admin facility handlers (Phase 4)
│   │   └── ...
│   ├── middleware/  # HTTP middleware ✓ Phase 2 Complete
│   │   ├── jwt.go                 # JWT token validation
│   │   ├── error.go               # Error handling
│   │   ├── logger.go              # Request logging
│   │   └── cors.go                # CORS headers
│   ├── auth/      # Authentication & RBAC ✓ Phase 2 Complete
│   │   ├── rbac.go                # Casbin RBAC setup
│   │   └── role.go                # Role definitions
│   ├── validator/ # Request validation ✓ Phase 2 Complete
│   │   └── validator.go
│   ├── response/ # API response helpers ✓ Phase 2 Complete
│   │   └── response.go
│   └── pagination/ # Pagination utilities ✓ Phase 2 Complete
│       └── pagination.go
├── pkg/            # Shared packages
│   └── utils/
├── migrations/    # Database migrations ✓ Phase 3 Complete
│   ├── 001_cruise_companies.up.sql
│   ├── 001_cruise_companies.down.sql
│   ├── 002_cruises.up.sql
│   ├── 002_cruises.down.sql
│   ├── 003_cabin_types.up.sql
│   ├── 003_cabin_types.down.sql
│   ├── 004_facility_categories.up.sql
│   ├── 004_facility_categories.down.sql
│   ├── 005_facilities.up.sql
│   ├── 005_facilities.down.sql
│   └── ...
├── tests/        # Test suites
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── docs/        # Swagger/OpenAPI docs
│   └── swagger.json
├── Dockerfile
├── docker-compose.yml
go.mod
go.sum
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
├── components/          # Vue components ✓ Phase 4
│   ├── common/         # Shared components
│   │   ├── Loading.vue
│   │   ├── Error.vue
│   │   └── Empty.vue
│   ├── forms/         # Form components
│   │   ├── ImageUpload.vue        # ✓ Image upload with drag-drop
│   │   └── TipTapEditor.vue       # ✓ Rich text editor
│   └── navigation/   # Navigation components
│       ├── Sidebar.vue            # ✓ Admin sidebar
│       ├── Header.vue             # ✓ Admin header
│       └── Footer.vue             # ✓ Admin footer
├── composables/      # Vue composables
│   ├── useAuth.ts
│   └── useApi.ts
├── layouts/         # Nuxt layouts ✓ Phase 4
│   ├── default.vue
│   ├── admin.vue                  # ✓ Admin dashboard layout
│   └── auth.vue                   # ✓ Auth pages layout
├── middleware/     # Route middleware
│   └── auth.ts
├── pages/         # Application pages ✓ Phase 4
│   ├── login.vue                  # ✓ Login page
│   ├── index.vue                  # ✓ Dashboard
│   ├── cruises/
│   │   ├── index.vue              # ✓ Cruise list
│   │   └── [id].vue               # ✓ Cruise edit/create
│   ├── cabin-types/
│   │   └── index.vue              # ✓ Cabin type management
│   ├── facilities/
│   │   └── index.vue              # ✓ Facility management
│   ├── routes/                    # (Phase 5)
│   ├── voyages/                   # (Phase 5)
│   ├── cabins/                    # (Phase 5)
│   ├── orders/                    # (Phase 6)
│   ├── staffs/                    # (Phase 7)
│   └── settings/                  # (Phase 7)
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
│   ├── cruise/           # Cruise-specific components ✓ Phase 3 Complete
│   │   ├── CruiseCard.vue
│   │   ├── ImageGallery.vue
│   │   ├── CabinTypeAccordion.vue
│   │   └── FacilityTabs.vue
│   ├── booking/         # Booking flow components (Phase 5)
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
├── pages/              # Public pages ✓ Phase 3 Complete
│   ├── index.vue      # Home
│   ├── cruises/
│   │   ├── index.vue  # List with filters
│   │   └── [id].vue   # Detail page ✓
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
├── tests/              # ✓ Phase 3 Complete
│   ├── components/
│   │   └── CruiseCard.spec.ts
│   └── e2e/
│       └── cruise-browsing.spec.ts
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
├── components/           # ✓ Phase 3 Complete
│   └── CruiseCard.vue
├── pages/               # Mini program pages ✓ Phase 3 Complete
│   ├── index/
│   │   └── index.vue    # Home page
│   ├── cruises/
│   │   ├── index.vue    # List page
│   │   └── detail.vue   # Detail page
│   ├── booking/         # (Phase 5)
│   ├── payment/         # (Phase 5)
│   ├── orders/          # (Phase 6)
│   └── profile/         # (Phase 7)
├── static/              # Static assets
├── store/               # Pinia stores
├── utils/
│   └── api.ts
├── tests/               # ✓ Phase 3 Complete
│   └── components/
│       └── CruiseCard.spec.ts
├── manifest.json        # WeChat Mini Program config
├── pages.json           # Page routes
├── App.vue
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

### Phase 3: Cruise Browsing ✓ (20/20)
- [x] Database migrations (5 tables: cruise_companies, cruises, cabin_types, facility_categories, facilities)
- [x] Domain models (CruiseCompany, Cruise, CabinType, Facility, FacilityCategory)
- [x] Repository layer (CruiseRepository, CabinTypeRepository, FacilityRepository)
- [x] Service layer (CruiseService, CabinTypeService, FacilityService, FacilityCategoryService)
- [x] HTTP handlers (cruise.go, cabin_type.go, facility.go)
- [x] API routes configuration
- [x] Frontend Web pages (index.vue, cruises/index.vue, cruises/[id].vue)
- [x] Frontend Web components (CruiseCard, ImageGallery, CabinTypeAccordion, FacilityTabs)
- [x] Frontend Mini Program pages (home, cruises list, detail)
- [x] Frontend Mini Program components (CruiseCard)
- [x] Unit tests for service layer
- [x] Integration tests for handlers
- [x] Component tests for frontend
- [x] E2E tests for cruise browsing

### Phase 4: Backend Management ✓ (16/16)
- [x] Admin CRUD handlers with image upload (admin_cruise.go, admin_cabin_type.go, admin_facility.go)
- [x] MinIO file upload service (storage.go)
- [x] Admin routes with RBAC (admin_routes.go)
- [x] Admin layout with sidebar navigation (layouts/admin.vue)
- [x] Login page (pages/login.vue)
- [x] Cruise management pages (list, create/edit)
- [x] Image upload component (ImageUpload.vue)
- [x] Cabin type management pages (pages/cabin-types/)
- [x] Facility management pages (pages/facilities/)
- [x] Admin handler tests (100% coverage)
- [x] Admin component tests (Vitest)
- [x] Admin E2E tests (Playwright)

### Next: Phase 5 - Online Booking & Payment
- [ ] Route and voyage domain models
- [ ] Cabin inventory management
- [ ] Booking workflow
- [ ] Payment integration
- [ ] Order management

---

## Contact & Resources

- **Constitution**: `.specify.specify/memory/constitution.md`
- **Specification**: `specs/001-cruise-booking-system/spec.md`
- **Tasks**: `specs/001-cruise-booking-system/tasks.md`
- **API Docs**: `specs/001-cruise-booking-system/contracts/openapi.yaml`

---

## Phase 3 & 4 File Inventory

### Phase 3: Cruise Browsing

#### Backend (16 files)

**Migrations:**
- `migrations/001_cruise_companies.up.sql`
- `migrations/001_cruise_companies.down.sql`
- `migrations/002_cruises.up.sql`
- `migrations/002_cruises.down.sql`
- `migrations/003_cabin_types.up.sql`
- `migrations/003_cabin_types.down.sql`
- `migrations/004_facility_categories.up.sql`
- `migrations/004_facility_categories.down.sql`
- `migrations/005_facilities.up.sql`
- `migrations/005_facilities.down.sql`

**Domain Models:**
- `internal/domain/cruise_company.go`
- `internal/domain/cruise.go`
- `internal/domain/cabin_type.go`
- `internal/domain/facility.go`

**Repositories:**
- `internal/repository/cruise.go`
- `internal/repository/cabin_type.go`
- `internal/repository/facility.go`

**Services:**
- `internal/service/cruise.go`
- `internal/service/cruise_test.go` (tests)
- `internal/service/cabin_type.go`
- `internal/service/facility.go`

**Handlers:**
- `cmd/api/routes.go`
- `internal/handler/cruise.go`
- `internal/handler/cruise_test.go` (tests)
- `internal/handler/cabin_type.go`
- `internal/handler/facility.go`

#### Frontend Web (7 files)

**Pages:**
- `pages/cruises/index.vue`
- `pages/cruises/[id].vue`

**Components:**
- `components/cruise/CruiseCard.vue`
- `components/cruise/ImageGallery.vue`
- `components/cruise/CabinTypeAccordion.vue`
- `components/cruise/FacilityTabs.vue`

**Tests:**
- `tests/components/CruiseCard.spec.ts`
- `tests/e2e/cruise-browsing.spec.ts`

#### Frontend Mini Program (6 files)

**Pages:**
- `pages/index/index.vue`
- `pages/cruises/index.vue`
- `pages/cruises/detail.vue`

**Components:**
- `components/CruiseCard.vue`

**Tests:**
- `tests/components/CruiseCard.spec.ts`

### Phase 4: Backend Management

#### Backend (8 files)

**Services:**
- `internal/service/storage.go` - MinIO file upload service

**Admin Handlers:**
- `cmd/api/admin_routes.go` - Admin routes with RBAC
- `internal/handler/admin_cruise.go` - Admin cruise CRUD
- `internal/handler/admin_cabin_type.go` - Admin cabin type CRUD
- `internal/handler/admin_facility.go` - Admin facility CRUD

**Admin Handler Tests:**
- `internal/handler/admin_cruise_test.go`
- `internal/handler/admin_cabin_type_test.go`
- `internal/handler/admin_facility_test.go`

#### Frontend Admin (12 files)

**Layouts:**
- `layouts/admin.vue` - Admin dashboard layout

**Pages:**
- `pages/login.vue` - Login page
- `pages/cruises/index.vue` - Cruise list
- `pages/cruises/[id].vue` - Cruise create/edit
- `pages/cabin-types/index.vue` - Cabin type management
- `pages/facilities/index.vue` - Facility management

**Components:**
- `components/ImageUpload.vue` - Drag-drop image upload
- `components/TipTapEditor.vue` - Rich text editor
- `components/layout/Header.vue` - Admin header
- `components/layout/Sidebar.vue` - Admin sidebar
- `components/layout/Footer.vue` - Admin footer

**Tests:**
- `tests/components/ImageUpload.spec.ts`
- `tests/e2e/admin-panel.spec.ts`

---

**End of Project Structure Guide**

*Use this document as a quick reference during code reviews and development.*
