# Implementation Plan: CruiseBooking Platform

**Feature**: 邮轮舱位预订平台 (CruiseBooking)  
**Branch**: main  
**Spec File**: specs/001-cruise-booking-system/spec.md  
**Created**: 2026-02-10  

---

## Phase 0: Research & Design Decisions

**Status**: ✓ Complete  
**Output**: [research.md](./plan/research.md)  

### Key Decisions Documented

1. **Architecture**: Modular monolith (Go) + 3 frontend apps (Nuxt 4, Nuxt 4 SSR, uni-app)
2. **Inventory Management**: Redis distributed locks with 15-min timeout
3. **Payment Handling**: Idempotent callbacks + reconciliation jobs
4. **Database**: PostgreSQL with soft deletes, JSONB for flexible attributes
5. **Testing Strategy**: 100% coverage mandate with TDD approach

### Technology Stack (Constitution-Compliant)

| Layer | Technology | Version |
|-------|------------|---------|
| Backend | Go + Gin + GORM | 1.26, v1.11.0, v2.x |
| Database | PostgreSQL | 17.x |
| Cache | Redis | 7.4.x |
| Search | Meilisearch | 1.12.x |
| Queue | NATS JetStream | 2.11.x |
| Storage | MinIO | Latest |
| Web Frontend | Nuxt 4 + Vue 3 | 4.3.0, 3.5.28 |
| Mini Program | uni-app (Vue 3) | HBuilderX |

---

## Phase 1: Data Model & API Contracts

**Status**: Ready to Execute  

### 1.1 Data Model Design

**Entity Count**: 18 core entities  
**Output**: [data-model.md](./plan/data-model.md)  

#### Core Entities

**Cruise Domain**:
- `cruise_companies` - Cruise line operators
- `cruises` - Individual ships with specs
- `cabin_types` - Cabin categories (Inside, Oceanview, Balcony, Suite)
- `facilities` + `facility_categories` - Ship amenities

**Booking Domain**:
- `routes` - Cruise itineraries
- `voyages` - Specific sailing instances
- `cabins` - Individual bookable units
- `cabin_inventory` - Real-time stock tracking
- `cabin_prices` - Date-specific pricing

**Order Domain**:
- `orders` + `order_items` - Booking records
- `passengers` - Guest information (encrypted)
- `payments` - Transaction records
- `refund_requests` - Cancellation/refund workflow

**User Domain**:
- `users` - Customer accounts
- `frequent_passengers` - Saved guest profiles
- `staffs` + `roles` - Admin RBAC

**Content Domain**:
- `notifications` - User messages
- `travelogues` - Community content (V2.0)

#### Key Design Patterns

1. **Soft Deletes**: All business entities have `deleted_at` for audit trail
2. **Optimistic Locking**: `cabin_inventory.version` for concurrent updates
3. **State Machines**: Order status transitions with event logging
4. **JSONB Flexibility**: Dynamic attributes (images, tags, amenities)

### 1.2 API Contracts

**Output Directory**: [contracts/](./contracts/)  

#### RESTful API Structure

```
/api/v1/
├── /auth                    # Authentication
├── /cruises                 # Cruise management
├── /cabin-types             # Cabin categories
├── /facilities              # Ship amenities
├── /routes                  # Itineraries
├── /voyages                 # Sailing dates
├── /cabins                  # Inventory & booking
├── /orders                  # Booking workflow
├── /payments                # Payment processing
├── /users                   # Customer profiles
├── /staffs                  # Admin management
├── /notifications           # Messages
└── /analytics               # Reports & dashboards
```

#### Key Endpoints

**Booking Flow**:
```
POST   /api/v1/orders                    # Create order
POST   /api/v1/orders/{id}/lock          # Lock inventory
POST   /api/v1/payments                  # Initiate payment
POST   /api/v1/payments/callback         # Webhook handler
POST   /api/v1/orders/{id}/confirm       # Confirm booking
```

**Inventory Management**:
```
GET    /api/v1/cabins/availability       # Check availability
POST   /api/v1/cabins/inventory/adjust   # Manual adjustment
GET    /api/v1/cabins/prices/calendar    # Price calendar
```

**Admin Operations**:
```
CRUD   /api/v1/cruises                   # Cruise management
CRUD   /api/v1/cabin-types               # Cabin type management
GET    /api/v1/orders                    # Order list (filtered)
POST   /api/v1/refund-requests/{id}/approve  # Refund approval
```

### 1.3 Quick Start Guide

**Output**: [quickstart.md](./docs/quickstart.md)  

#### Local Development Setup

```bash
# 1. Clone and setup backend
cd backend
cp .env.example .env
docker-compose up -d postgres redis meilisearch minio nats
go mod tidy
go run cmd/api/main.go

# 2. Setup management frontend
cd frontend-admin
npm install
npm run dev

# 3. Setup customer web frontend
cd frontend-web
npm install
npm run dev

# 4. Setup mini program
cd frontend-mini
npm install
# Open in HBuilderX
```

#### Environment Requirements

- **Go**: 1.26+
- **Node.js**: 20.x LTS
- **Docker & Docker Compose**: Latest
- **PostgreSQL**: 17.x
- **Redis**: 7.4.x

---

## Phase 2: Implementation Sprints

### Sprint 1 (Weeks 1-2): Foundation + Cruise Display

**Backend Tasks**:
- [ ] Initialize Go project with Gin, GORM, middleware
- [ ] Setup PostgreSQL schemas: cruises, cabin_types, facilities
- [ ] Implement Cruise CRUD API with image upload (MinIO)
- [ ] Implement Cabin Type CRUD with rich text (TipTap content)
- [ ] Implement Facility management
- [ ] Setup JWT + Casbin RBAC middleware
- [ ] Write unit + integration tests (100% coverage)
- [ ] Generate Swagger documentation

**Management Frontend Tasks**:
- [ ] Initialize Nuxt 4 with TypeScript, Tailwind, Nuxt UI
- [ ] Login page with JWT integration
- [ ] Layout: sidebar + header + content
- [ ] Cruise management pages (list, create, edit)
- [ ] TipTap editor integration for descriptions
- [ ] Cabin type management pages
- [ ] Facility management pages
- [ ] Write Vitest + Playwright tests (100% coverage)

**Customer Web Frontend Tasks**:
- [ ] Initialize Nuxt 4 SSR project
- [ ] Home page layout
- [ ] Cruise list page with filters
- [ ] Cruise detail page with galleries
- [ ] Cabin type display (accordion)
- [ ] Facilities display (tab navigation)

**Mini Program Tasks**:
- [ ] Initialize uni-app Vue 3 project
- [ ] WeChat login integration
- [ ] Home page
- [ ] Cruise list page
- [ ] Cruise detail page

**Sprint 1 Deliverable**: Browseable cruise catalog with full management backend

---

### Sprint 2 (Weeks 3-4): Cabin Management + Browsing

**Backend Tasks**:
- [ ] Setup tables: routes, voyages, cabins, inventory, prices
- [ ] Cabin CRUD API with SKU attributes
- [ ] Inventory management API with Redis locking
- [ ] Price calendar API with date-range pricing
- [ ] Availability query endpoint
- [ ] Inventory alert system
- [ ] Tests (100% coverage)

**Management Frontend Tasks**:
- [ ] Route management pages
- [ ] Voyage management pages
- [ ] Cabin management (list with filters)
- [ ] Inventory management interface
- [ ] Price calendar editor
- [ ] Inventory alert threshold settings
- [ ] Tests (100% coverage)

**Customer Frontend Tasks**:
- [ ] Cabin list page with filters (route, date, type)
- [ ] Cabin detail page
- [ ] Price calendar display
- [ ] Deck position visualization
- [ ] Tests (100% coverage)

**Sprint 2 Deliverable**: Complete cabin browsing with real-time inventory

---

### Sprint 3 (Weeks 5-6): Booking + Payment

**Backend Tasks**:
- [ ] Setup tables: orders, order_items, passengers, payments
- [ ] Order creation API with inventory locking
- [ ] Order state machine implementation
- [ ] WeChat Pay integration (JSAPI + Native)
- [ ] Payment callback handler with idempotency
- [ ] Order timeout handling (NATS delayed messages)
- [ ] Order query endpoints
- [ ] Tests (100% coverage)

**Management Frontend Tasks**:
- [ ] Order list page with status tabs
- [ ] Order detail page (info + passengers + payments)
- [ ] Order operations (confirm, cancel, add note)
- [ ] Tests (100% coverage)

**Customer Frontend Tasks**:
- [ ] Booking flow: voyage selection → cabin → passengers → confirmation
- [ ] Passenger information form
- [ ] Payment page with WeChat Pay
- [ ] Payment result page
- [ ] My orders list
- [ ] Order detail page
- [ ] Tests (100% coverage)

**Sprint 3 Deliverable**: End-to-end booking and payment workflow

---

### Sprint 4 (Weeks 7-8): Users + Refunds + Notifications + Analytics

**Backend Tasks**:
- [ ] Setup tables: users, staffs, roles, notifications, refund_requests
- [ ] User registration/login (WeChat, phone)
- [ ] Staff management API with Casbin
- [ ] Refund request workflow API
- [ ] Notification service (WeChat template, SMS)
- [ ] Financial reconciliation reports
- [ ] Dashboard analytics APIs
- [ ] Tests (100% coverage)

**Management Frontend Tasks**:
- [ ] Staff management pages
- [ ] Role & permission management
- [ ] Refund request list and approval
- [ ] Notification template configuration
- [ ] Financial reconciliation page
- [ ] Dashboard with ECharts
- [ ] Shop/brand info configuration
- [ ] Tests (100% coverage)

**Customer Frontend Tasks**:
- [ ] User profile page
- [ ] Frequent passengers management
- [ ] Refund request submission
- [ ] Notification center
- [ ] Tests (100% coverage)

**Sprint 4 Deliverable**: MVP complete - usable booking system

---

### Phase 3: V1.0 Enhancement (Months 3-4)

**Sprints 5-6**: Shore excursions, E-tickets, Countdown, Customer Service
- Shore excursion product management
- PDF e-ticket generation
- Trip countdown widget
- AI + human customer service integration
- Frequent passenger quick-fill

**Sprints 7-8**: Reviews, Social Sharing, Membership
- Review system (ratings + photos/videos)
- Trip poster generation (Canvas/server-side)
- Membership levels + points system
- Referral program

---

### Phase 4: V1.5 Intelligence (Months 5-6)

**Sprints 9-10**: Smart Features
- Recommendation engine (collaborative filtering)
- Route calendar view with pricing
- Price trend analysis charts
- Cabin comparison tool
- 360° VR cabin preview
- Interactive deck maps (SVG/Canvas)

**Sprints 11-12**: Advanced Booking
- Dynamic pricing engine (rule-based)
- Multi-currency support (CNY/USD/HKD/JPY)
- ID document OCR integration
- Installment payments (deposit + final)
- Group booking with Excel upload

---

### Phase 5: V2.0 Ecosystem (Months 7-9)

**Sprints 13-14**: Community
- Travelogue community (UGC)
- Group buying / cabin sharing
- Points mall

**Sprints 15-16**: Operations + Global
- Multi-channel inventory distribution
- Revenue management dashboard
- CRM with lifecycle management
- Automated marketing engine
- Multi-language support (i18n)
- Real-time WebSocket updates

---

## Constitution Check

### Compliance Verification

| Principle | Status | Implementation |
|-----------|--------|----------------|
| T2.1 Technology Stack | ✓ Compliant | Go 1.26, Gin, GORM, PostgreSQL 17, Redis 7.4 |
| T2.2 Frontend Stack | ✓ Compliant | Nuxt 4.3, Vue 3.5.28, Pinia v3, TypeScript 5.9 |
| T3.1 100% Test Coverage | ✓ Enforced | CI gates, testify, gomock, Vitest, Playwright |
| Q4.1 API Documentation | ✓ Enforced | Swagger/OpenAPI 3.1, CI validation |
| Q4.2 Code Review | ✓ Required | All PRs require review, no direct main commits |
| Q4.3 Database Migrations | ✓ Required | GORM AutoMigrate + migration files |
| Q4.4 Lint Compliance | ✓ Enforced | golangci-lint, ESLint, Prettier in CI |
| A6.2 Auth/Authz | ✓ Implemented | JWT + Casbin RBAC |
| A6.3 Inventory Locks | ✓ Implemented | Redis distributed locks |
| A6.4 Payment Security | ✓ Implemented | Signature validation, idempotency |
| A6.5 Observability | ✓ Planned | Prometheus, Grafana, Loki integration |

### Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Concurrent overselling | Redis locks + inventory versioning |
| Payment callback failure | Idempotency keys + reconciliation jobs |
| Third-party service downtime | Circuit breakers + cached fallbacks |
| Performance degradation | Kubernetes HPA + caching layers |
| Security breaches | Encryption, RBAC, audit logs |

---

## Deliverables

### Documentation
- [x] [research.md](./plan/research.md) - Technical research & decisions
- [ ] data-model.md - Database schema (18 entities)
- [ ] contracts/openapi.yaml - REST API specification
- [ ] quickstart.md - Developer onboarding guide

### Code Structure
```
/
├── backend/               # Go monolith
│   ├── cmd/api/          # Entry point
│   ├── internal/         # Business logic
│   │   ├── domain/       # Entities
│   │   ├── service/      # Business services
│   │   ├── repository/   # Data access
│   │   ├── handler/      # HTTP handlers
│   │   └── middleware/   # Auth, logging
│   ├── pkg/              # Shared utilities
│   └── migrations/       # DB migrations
│
├── frontend-admin/       # Nuxt 4 - Management
├── frontend-web/         # Nuxt 4 - Customer Web
├── frontend-mini/        # uni-app - Mini Program
│
├── specs/001-cruise-booking-system/
│   ├── spec.md           # Feature specification
│   ├── plan/
│   │   ├── research.md   # Research document
│   │   ├── data-model.md # Database design
│   │   └── impl-plan.md  # This file
│   ├── contracts/        # API contracts
│   ├── docs/             # Documentation
│   └── .agent/           # Agent context
│
└── .github/workflows/    # CI/CD pipelines
```

---

## Next Steps

1. **Review Phase 0** research decisions
2. **Finalize Phase 1** data model and API contracts
3. **Setup development** environment per quickstart
4. **Begin Sprint 1** implementation
5. **Run CI/CD** validation after each sprint

---

**Ready for Development** ✓

*This plan aligns with Constitution v1.0.0 and Specification v1.0.0*
