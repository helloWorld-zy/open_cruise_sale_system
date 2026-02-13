# Project Structure Guide: CruiseBooking Platform

**Purpose**: This document provides a comprehensive overview of the project structure to facilitate code reviews and reduce token consumption during AI-assisted reviews.

**Last Updated**: 2026-02-12
**Version**: 2.0.0
**Phase**: All 11 Phases Complete ✓ (206/206 tasks - 100%)
**Status**: Production Ready  

---

## Repository Overview

```
open_cruise_sale_system/
├── .github/                    # CI/CD workflows
├── .specify.specify/           # Project constitution & templates
├── specs/                      # Specifications & documentation
│   └── 001-cruise-booking-system/
│       ├── spec.md            # Feature specification (50KB+, 含技术规范章节)
│       ├── tasks.md           # Task breakdown (206 tasks)
│       ├── plan/              # Implementation plans
│       ├── contracts/         # API contracts (OpenAPI)
│       ├── docs/              # Developer documentation
│       └── checklists/        # Quality checklists
├── backend/                   # Go monolith (Phase 1, 2, 3 & 4 ✓)
├── frontend-admin/           # Nuxt 4 - Management (Phase 1, 2, 4 ✓)
├── frontend-web/            # Nuxt 4 - Customer Web (Phase 1, 3, 5 ✓)
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
│   │   ├── booking.go             # Route, Voyage, Cabin, Inventory, Price domains ✓ Phase 5
│   │   ├── order.go               # Order, OrderItem, Passenger, Payment domains ✓ Phase 5
│   │   ├── user.go                # (Phase 7)
│   │   └── ...
│   ├── repository/  # Data access layer ✓ Phase 3 Complete
│   │   ├── cruise.go              # CruiseRepository
│   │   ├── cabin_type.go          # CabinTypeRepository
│   │   ├── facility.go            # FacilityRepository + FacilityCategoryRepository
│   │   ├── voyage.go              # VoyageRepository ✓ Phase 5
│   │   ├── cabin.go               # CabinRepository ✓ Phase 5
│   │   ├── inventory.go           # InventoryRepository ✓ Phase 5
│   │   ├── price.go               # PriceRepository ✓ Phase 5
│   │   ├── order.go               # OrderRepository ✓ Phase 5
│   │   └── ...
│   ├── service/    # Business logic layer ✓ Phase 3 & 4 Complete
│   │   ├── cruise.go              # CruiseService ✓
│   │   ├── cruise_test.go         # Unit tests ✓
│   │   ├── cabin_type.go          # CabinTypeService ✓
│   │   ├── facility.go            # FacilityService + FacilityCategoryService ✓
│   │   ├── storage.go             # ✓ StorageService (Phase 4)
│   │   ├── inventory.go           # InventoryService ✓ Phase 5
│   │   ├── price.go               # PriceService ✓ Phase 5
│   │   ├── order_state.go         # Order state machine ✓ Phase 5
│   │   ├── order.go               # OrderService ✓ Phase 5
│   │   └── ...
│   ├── handler/   # HTTP handlers (controllers) ✓ Phase 3 & 4 Complete
│   │   ├── auth.go                # Auth handlers (Phase 2)
│   │   ├── cruise.go              # Public cruise handlers ✓
│   │   ├── cruise_test.go         # Integration tests ✓
│   │   ├── cabin_type.go          # Public cabin type handlers ✓
│   │   ├── facility.go            # Public facility handlers ✓
│   │   ├── order.go               # Order handlers ✓ Phase 5
│   │   ├── payment.go             # Payment handlers ✓ Phase 5
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
│   ├── payment/    # Payment providers ✓ Phase 5
│   │   ├── wechat.go     # WeChat Pay V3 integration
│   │   └── service.go    # Payment service
│   ├── jobs/       # Background jobs ✓ Phase 5
│   │   └── order_timeout.go  # Order timeout cleanup
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
│   ├── booking/         # Booking flow components ✓ Phase 5
│   │   ├── BookingWizard.vue    # Main booking wizard
│   │   ├── SelectVoyage.vue     # Step 1: Voyage selection
│   │   ├── SelectCabin.vue      # Step 2: Cabin selection
│   │   ├── PassengerForm.vue    # Step 3: Passenger info
│   │   └── OrderConfirm.vue     # Step 4: Order confirmation
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
│   │   └── index.vue  # Booking wizard page ✓ Phase 5
│   ├── payment/
│   │   ├── [orderId].vue  # Payment page ✓ Phase 5
│   │   └── result.vue     # Payment result (Phase 6)
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

### Phase 5: Online Booking & Payment ✓ (Complete)
- [x] Database migrations (9 tables: routes, voyages, cabins, cabin_inventory, cabin_prices, orders, order_items, passengers, payments)
- [x] Domain models (Route, Voyage, Cabin, CabinInventory, CabinPrice, Order, OrderItem, Passenger, Payment)
- [x] Repository layer (VoyageRepository, CabinRepository, InventoryRepository, PriceRepository, OrderRepository)
- [x] Service layer (InventoryService, PriceService, OrderStateService, OrderService)
- [x] Service unit tests (order_test.go, payment/service_test.go)
- [x] Payment integration (WeChat Pay V3 SDK, PaymentService)
- [x] Background jobs (OrderTimeoutJob with NATS)
- [x] HTTP handlers (order.go, payment.go)
- [x] Booking wizard component (4-step flow: SelectVoyage, SelectCabin, PassengerForm, OrderConfirm)
- [x] Payment page with QR code and polling
- [x] Payment result page
- [x] Mini Program booking pages (pages/booking/, pages/payment/)
- [x] Mini Program WeChat SDK payment integration
- [x] Jest tests for Mini Program booking components
- [x] Playwright E2E tests for booking flow (tests/e2e/booking-flow.spec.ts)

### Phase 6: Order Management & Refunds ✓ (Complete)
- [x] Refund requests table migration
- [x] RefundRequest domain model
- [x] Order query endpoints (order_query.go)
- [x] Refund workflow service (refund.go)
- [x] Admin order handlers (admin_order.go)
- [x] My Orders page (pages/orders/index.vue)
- [x] Order detail page (pages/orders/[id].vue)
- [x] Refund request modal (RefundRequestModal.vue)
- [x] Admin order list (pages/orders/index.vue)
- [x] Admin order detail (pages/orders/[id].vue)
- [x] Admin refund management (pages/refunds/index.vue)

### Phase 7: User Authentication & Account ✓ (Complete)
- [x] Users table migration
- [x] FrequentPassengers table migration
- [x] User domain model
- [x] FrequentPassenger domain model
- [x] User repository
- [x] WeChat login service (wechat_auth.go)
- [x] SMS verification service (sms.go)
- [x] User handlers
- [x] Login modal component
- [x] User profile page
- [x] Frequent passengers management
- [x] Mini program login page
- [x] Auth flow E2E tests

### Phase 8: Notifications & Reminders ✓ (Complete)
- [x] Notifications table migration (018_notifications.up/down.sql)
- [x] NotificationSettings table migration (019_notification_settings.up/down.sql)
- [x] Notification domain model (notification.go)
- [x] Notification repository (notification.go)
- [x] Notification service (notification.go)
- [x] WeChat template message sender (wechat.go)
- [x] SMS sender (sms.go)
- [x] Inventory alert job (inventory_alert.go)
- [x] Notification handler (notification.go)
- [x] NotificationCenter component
- [x] Notification settings page

### Phase 9: Smart Recommendations & Analytics ✓ (Complete)
- [x] User behavior tracking (analytics/tracking.go)
- [x] Recommendation engine (recommendation/engine.go)
- [x] Price trend analysis (analytics/price_trends.go)
- [x] Analytics API endpoints (handler/analytics.go)
- [x] Recommendation carousel component
- [x] Price calendar view component
- [x] Cabin comparison tool
- [x] Analytics dashboard (admin)

### Phase 10: Social Sharing & Community ✓ (Complete)
- [x] Travelogues table migration (020_travelogues.up/down.sql)
- [x] Travelogue domain model (travelogue.go)
- [x] Review system (service/review.go)
- [x] Poster generation service (service/poster.go)
- [x] Invitation system (service/invitation.go)
- [x] Review form component (ReviewForm.vue)
- [x] Travelogue editor component (TravelogueEditor.vue)
- [x] Community page (community/index.vue)
- [x] Share poster modal (SharePoster.vue)

### Phase 11: Polish & Production ✓ (Complete)
- [x] Redis caching layer (cache/manager.go)
- [x] Meilisearch indexing for search
- [x] Database query optimization
- [x] API rate limiting middleware
- [x] Frontend code splitting and lazy loading
- [x] Prometheus metrics collection
- [x] Grafana dashboards
- [x] Loki log aggregation
- [x] Distributed tracing
- [x] Health check endpoints (/health, /health/ready, /health/live)
- [x] API request signing
- [x] SQL injection prevention tests
- [x] Security audit
- [x] GDPR data export/deletion
- [x] Complete API documentation
- [x] Deployment guide
- [x] Operation runbook
- [x] User manual

### Project Status: Production Ready ✓
**All 11 Phases Complete** (206/206 tasks - 100%)

### Latest Updates (2026-02-12)
- ✅ **Technical Specifications** added to spec.md (300+ lines)
  - 并发控制与库存保护规范（双层锁机制、CAS操作、分布式锁）
  - 支付安全技术规范（幂等性、重复支付检测、敏感数据加密）
  - 监控与告警规范（42个量化指标）
  - 代码质量要求（审查清单、测试要求）
- ✅ **Code Quality Checklists** completed (70 items)
  - 需求完整性、清晰度、一致性全覆盖
  - 并发控制与超卖防护详细规范
  - 支付安全与幂等性实现要求
- ✅ **Deployment Guide** published (production-ready)

---

## Contact & Resources

- **Constitution**: `.specify.specify/memory/constitution.md`
- **Specification**: `specs/001-cruise-booking-system/spec.md` (含技术规范章节)
- **Tasks**: `specs/001-cruise-booking-system/tasks.md`
- **API Docs**: `specs/001-cruise-booking-system/contracts/openapi.yaml`
- **Quality Checklists**: 
  - `specs/001-cruise-booking-system/checklists/requirements.md` (需求质量检查)
  - `specs/001-cruise-booking-system/checklists/code-quality.md` (代码质量检查 - 70项)
  - `specs/001-cruise-booking-system/checklists/code-quality-report.md` (检查完成报告)
- **Deployment Guide**: `DEPLOYMENT_GUIDE.md`

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

### Phase 5: Online Booking & Payment

#### Backend (45 files)

**Migrations:**
- `migrations/006_routes.up.sql` / `006_routes.down.sql`
- `migrations/007_voyages.up.sql` / `007_voyages.down.sql`
- `migrations/008_cabins.up.sql` / `008_cabins.down.sql`
- `migrations/009_cabin_inventory.up.sql` / `009_cabin_inventory.down.sql`
- `migrations/010_cabin_prices.up.sql` / `010_cabin_prices.down.sql`
- `migrations/011_orders.up.sql` / `011_orders.down.sql`
- `migrations/012_order_items.up.sql` / `012_order_items.down.sql`
- `migrations/013_passengers.up.sql` / `013_passengers.down.sql`
- `migrations/014_payments.up.sql` / `014_payments.down.sql`

**Domain Models:**
- `internal/domain/booking.go` - Route, Voyage, Cabin, CabinInventory, CabinPrice
- `internal/domain/order.go` - Order, OrderItem, Passenger, Payment

**Repositories:**
- `internal/repository/voyage.go` - VoyageRepository
- `internal/repository/cabin.go` - CabinRepository
- `internal/repository/inventory.go` - InventoryRepository with optimistic locking
- `internal/repository/price.go` - PriceRepository
- `internal/repository/order.go` - OrderRepository with payments

**Services:**
- `internal/service/inventory.go` - InventoryService with lock/unlock/confirm
- `internal/service/price.go` - PriceService with calculations
- `internal/service/order_state.go` - Order state machine
- `internal/service/order.go` - OrderService with full lifecycle

**Service Tests:**
- `internal/service/order_test.go` - Order service unit tests (100% coverage)

**Payment:**
- `internal/payment/wechat.go` - WeChat Pay V3 SDK
- `internal/payment/service.go` - Payment service
- `internal/payment/service_test.go` - Payment integration tests (100% coverage)

**Jobs:**
- `internal/jobs/order_timeout.go` - Order timeout cleanup with NATS

**Handlers:**
- `internal/handler/order.go` - Order handlers (create, list, cancel)
- `internal/handler/payment.go` - Payment handlers (create, callback, refund)

#### Frontend Web (9 files)

**Components:**
- `components/booking/BookingWizard.vue` - Main booking wizard
- `components/booking/SelectVoyage.vue` - Step 1: Voyage selection
- `components/booking/SelectCabin.vue` - Step 2: Cabin selection
- `components/booking/PassengerForm.vue` - Step 3: Passenger form
- `components/booking/OrderConfirm.vue` - Step 4: Order confirmation

**Pages:**
- `pages/booking/index.vue` - Booking page
- `pages/payment/[orderId].vue` - Payment page with QR code
- `pages/payment/result.vue` - Payment result page

**Tests:**
- `tests/e2e/booking-flow.spec.ts` - Playwright E2E tests for booking flow

#### Frontend Mini Program (4 files)

**Pages:**
- `pages/booking/index.vue` - Booking wizard with 4-step flow
- `pages/payment/result.vue` - Payment result page

**Tests:**
- `tests/components/BookingPage.spec.js` - Jest tests for booking page
- `tests/components/PaymentResultPage.spec.js` - Jest tests for payment result

### Phase 6: Order Management & Refunds

#### Backend (8 files)

**Migrations:**
- `migrations/015_refund_requests.up.sql` / `015_refund_requests.down.sql`

**Domain Models:**
- `internal/domain/refund.go` - RefundRequest domain model

**Services:**
- `internal/service/refund.go` - RefundService with workflow

**Handlers:**
- `internal/handler/order_query.go` - Order query endpoints
- `internal/handler/admin_order.go` - Admin order and refund handlers

#### Frontend Web (Customer) (3 files)

**Pages:**
- `pages/orders/index.vue` - My orders page with statistics
- `pages/orders/[id].vue` - Order detail page

**Components:**
- `components/RefundRequestModal.vue` - Refund request form modal

#### Frontend Admin (3 files)

**Pages:**
- `pages/orders/index.vue` - Admin order list with filters
- `pages/orders/[id].vue` - Admin order detail view
- `pages/refunds/index.vue` - Refund approval interface

### Phase 7: User Authentication & Account

#### Backend (8 files)

**Migrations:**
- `migrations/016_users.up.sql` / `016_users.down.sql`
- `migrations/017_frequent_passengers.up.sql` / `017_frequent_passengers.down.sql`

**Domain Models:**
- `internal/domain/user.go` - User domain model
- `internal/domain/frequent_passenger.go` - FrequentPassenger domain model

**Repositories:**
- `internal/repository/user.go` - UserRepository with frequent passengers

**Services:**
- `internal/service/wechat_auth.go` - WeChat authentication service
- `internal/service/sms.go` - SMS verification service

**Handlers:**
- `internal/handler/user.go` - User handlers (profile, passengers)

#### Frontend Web (3 files)

**Components:**
- `components/LoginModal.vue` - Login modal with WeChat and SMS

**Pages:**
- `pages/profile/index.vue` - User profile page
- `pages/profile/passengers.vue` - Frequent passengers management

#### Frontend Mini Program (1 file)

**Pages:**
- `pages/login/index.vue` - WeChat and SMS login page

---

**End of Project Structure Guide**

*Use this document as a quick reference during code reviews and development.*
