# Tasks: CruiseBooking Platform

**Feature**: 邮轮舱位预订平台 (CruiseBooking)  
**Total Tasks**: 120+  
**Estimated Duration**: 16 Sprints (9 Months)  
**Branch**: main  

---

## Progress Overview

| Phase | User Story | Priority | Status | Progress |
|-------|-----------|----------|--------|----------|
| 1 | Setup | - | ✅ **Complete** | 20/20 (100%) |
| 2 | Foundational | - | ✅ **Complete** | 16/16 (100%) |
| 3 | US1: Cruise Browsing | P1 | ✅ **Complete** | 20/20 (100%) |
| 4 | US4: Backend Management | P1 | ✅ **Complete** | 16/16 (100%) |
| 5 | US2: Booking & Payment | P1 | ✅ Complete | 27/27 (100%) |
| 6 | US3: Order Management | P1 | ✅ Complete | 13/13 (100%) |
| 7 | US5: User Auth | P2 | ✅ Complete | 14/14 (100%) |
| 8 | US6: Notifications | P2 | ✅ Complete | 10/10 (100%) |
| 9 | US7: Smart Features | P3 | ⏳ Pending | 0/12 (0%) |
| 10 | US8: Social Features | P3 | ⏳ Pending | 0/10 (0%) |
| 11 | Polish | - | ⏳ Pending | 0/8 (0%) |

**Total Progress**: 138/206 (67.0%)

### Phase 1-4 Completion Summary

**Phase 1: Setup** ✅
- Backend infrastructure: 8/8 tasks
- Frontend projects: 8/8 tasks
- DevOps: 4/4 tasks

**Phase 2: Foundational** ✅
- Auth & Authorization: 5/5 tasks
- Common Infrastructure: 7/7 tasks
- Shared Frontend: 4/4 tasks

**Phase 3: US1 - Cruise Browsing** ✅
- ✅ **Backend Complete** (16/16 tasks)
  - Database migrations: 5/5
  - Domain models: 4/4
  - Repositories: 3/3
  - Services: 3/3
  - Handlers: 3/3
  - Routes: 1/1
- ✅ **Frontend Web** (9/9 tasks)
  - Pages: 3/3
  - Components: 4/4
  - Tests: 2/2
- ✅ **Frontend Mini** (5/5 tasks)
  - Pages: 3/3
  - Components: 1/1
  - Tests: 1/1

**Phase 4: US4 - Backend Management** ✅
- ✅ **Backend APIs** (6/6 tasks)
  - Admin CRUD handlers for cruises, cabin types, facilities
  - MinIO file upload service
  - Admin routes with RBAC protection
  - Handler tests
- ✅ **Frontend Admin Panel** (10/10 tasks)
  - Admin layout with sidebar navigation
  - Login page
  - Cruise management (list, create/edit)
  - Image upload component with drag-drop
  - Cabin type management
  - Facility management
  - Component tests (Vitest)
  - E2E tests (Playwright)

**Phase 5: US2 - Online Booking & Payment** ✅
- ✅ **Backend Booking Domain** (17/17 tasks)
  - Route, Voyage, Cabin, Inventory, Price migrations and models
  - Repositories with optimistic locking for inventory
  - Inventory and Price services
- ✅ **Backend Order Domain** (12/12 tasks)
  - Order, OrderItem, Passenger, Payment migrations and models
  - Order repository with full CRUD
  - Order state machine with transitions
  - Order service and handlers
- ✅ **Payment Integration** (3/3 tasks)
  - WeChat Pay V3 SDK integration
  - Payment callback handler with idempotency
  - Order timeout job with NATS
- ✅ **Frontend Web** (7/7 tasks)
  - Booking wizard component with 4-step flow
  - Voyage selection, cabin selection, passenger form, order confirmation
  - Payment page with QR code and polling
  - Payment result page
- ✅ **Testing** (2/2 tasks)
  - Order service unit tests (100% coverage)
  - Payment integration tests (100% coverage)
- ✅ **E2E Tests** (1/1 task)
  - Playwright E2E tests for booking flow
- ✅ **Mini Program** (3/3 tasks)
  - Booking pages with 4-step wizard
  - WeChat SDK payment integration
  - Jest tests for booking components

**Phase 6: US3 - Order Management & Refunds** ✅
- ✅ **Backend** (6/6 tasks)
  - Refund requests table migration
  - RefundRequest domain model
  - Order query endpoints
  - Refund workflow service
  - Admin order handlers
  - Refund service tests
- ✅ **Frontend Web** (4/4 tasks)
  - My Orders page with statistics
  - Order detail page
  - Refund request form
  - Order management E2E tests
- ✅ **Frontend Admin** (3/3 tasks)
  - Order management list
  - Order detail view with status update
  - Refund approval interface

**Phase 7: US5 - User Authentication & Account** ✅
- ✅ **Backend** (9/9 tasks)
  - Users and FrequentPassengers table migrations
  - User and FrequentPassenger domain models
  - User repository with CRUD operations
  - WeChat authentication service (code2session, phone login)
  - SMS verification service
  - User handlers (profile, frequent passengers)
  - User service tests
- ✅ **Frontend Web** (3/3 tasks)
  - Login modal with WeChat and SMS options
  - User profile page with statistics
  - Frequent passengers management
- ✅ **Mini Program** (1/1 task)
  - WeChat and SMS login page
- ✅ **E2E Tests** (1/1 task)
  - Auth flow E2E tests

**Next**: Phase 8 - Notifications & Reminders

---

## Summary

| Phase | User Story | Priority | Tasks | Independent Test |
|-------|-----------|----------|-------|------------------|
| 1 | Setup | - | 20/20 ✓ | Project initialized, all dependencies installed |
| 2 | Foundational | - | 16/16 ✓ | Infrastructure ready for all user stories |
| 3 | US1: Cruise Browsing | P1 | 20/20 ✓ | User can browse cruises and cabins without login |
| 4 | US4: Backend Management | P1 | 16/16 ✓ | Admin can manage all cruise data |
| 5 | US2: Booking & Payment | P1 | 27/27 ✓ | Complete booking flow with payment |
| 6 | US3: Order Management | P1 | 14 | Orders and refunds fully functional |
| 7 | US5: User Auth | P2 | 12 | Multi-platform login working |
| 8 | US6: Notifications | P2 | 10 | Multi-channel notifications working |
| 9 | US7: Smart Features | P3 | 12 | Recommendations and analytics |
| 10 | US8: Social Features | P3 | 10 | Community and sharing features |
| 11 | Polish | - | 8 | Production-ready system |

---

## Phase 1: Project Setup

**Goal**: Initialize all projects with correct structure and dependencies  
**Duration**: 2-3 days  
**Deliverable**: All 4 codebases (backend + 3 frontends) initialized and building

### Backend Setup

- [x] T001 [P] Initialize Go project structure in `backend/` with go.mod
- [x] T002 [P] Install Gin v1.11.0, GORM v2.x and configure router in `backend/cmd/api/main.go`
- [x] T003 [P] Setup PostgreSQL 17 connection with GORM in `backend/internal/database/db.go`
- [x] T004 [P] Setup Redis 7.4 client in `backend/internal/cache/redis.go`
- [x] T005 [P] Configure MinIO S3 client in `backend/internal/storage/minio.go`
- [x] T006 [P] Setup NATS JetStream connection in `backend/internal/messaging/nats.go`
- [x] T007 Create environment configuration with Viper in `backend/internal/config/config.go`
- [x] T008 Setup structured logging with Zap in `backend/internal/logger/logger.go`

### Frontend Setup

- [x] T009 [P] Initialize Nuxt 4.3.0 project in `frontend-admin/` with TypeScript 5.9
- [x] T010 [P] Initialize Nuxt 4.3.0 SSR project in `frontend-web/` with TypeScript 5.9
- [x] T011 [P] Initialize uni-app Vue 3 project in `frontend-mini/` with TypeScript
- [x] T012 Install and configure Tailwind CSS v4, Nuxt UI v3 in both Nuxt projects
- [x] T013 Install Pinia v3 in all frontend projects
- [x] T014 Setup Vitest + Vue Test Utils in Nuxt projects
- [x] T015 Setup Jest + @vue/test-utils in uni-app project
- [x] T016 Create shared types package in `shared/types/` for cross-project use

### DevOps Setup

- [x] T017 Create `docker-compose.yml` with PostgreSQL, Redis, Meilisearch, MinIO, NATS
- [x] T018 Create GitHub Actions CI workflow in `.github/workflows/ci.yml`
- [x] T019 Create `.env.example` files for all projects
- [x] T020 Setup Swagger generation config in `backend/`

---

## Phase 2: Foundational Infrastructure

**Goal**: Build shared infrastructure usable by all user stories  
**Duration**: 1 week  
**Deliverable**: Auth, common middleware, base models, and utilities ready

### Authentication & Authorization

- [x] T021 Implement JWT middleware in `backend/internal/middleware/jwt.go`
- [x] T022 Setup Casbin RBAC with policy model in `backend/internal/auth/rbac.go`
- [x] T023 Create role definitions (Super Admin, Operations, Finance, Customer Service)
- [x] T024 Implement password hashing utilities in `backend/internal/auth/role.go`
- [x] T025 Create auth handlers in `backend/internal/handler/auth.go`

### Common Infrastructure

- [x] T026 Create base model with soft delete in `backend/internal/domain/base.go`
- [x] T027 Implement pagination utilities in `backend/internal/pagination/pagination.go`
- [x] T028 Create API response wrapper in `backend/internal/response/response.go`
- [x] T029 Implement error handling middleware in `backend/internal/middleware/error.go`
- [x] T030 Create request logging middleware in `backend/internal/middleware/logger.go`
- [x] T031 Implement CORS middleware in `backend/internal/middleware/cors.go`
- [x] T032 Create request validation utilities in `backend/internal/validator/validator.go`

### Shared Frontend Components

- [x] T033 Create authentication store with Pinia in `frontend-admin/stores/auth.ts`
- [x] T034 Create base API client with ofetch in `frontend-admin/utils/api.ts`
- [x] T035 Implement base layout components (Header, Sidebar, Footer)
- [x] T036 Create common UI components (Loading, Error, Empty states)
- [x] T037 Setup route guards for authentication in `frontend-admin/middleware/auth.ts`

---

## Phase 3: US1 - Cruise Browsing & Cabin Selection

**Goal**: Users can browse cruises and view cabin details without login  
**Priority**: P1 (Critical)  
**Duration**: 2 weeks  
**Independent Test**: User can browse cruise list, view details, see cabin types and facilities without logging in

### Backend - Cruise Domain

- [x] T038 Create `cruise_companies` table migration in `backend/migrations/001_cruise_companies.up.sql`
- [x] T039 Create `cruises` table migration in `backend/migrations/002_cruises.up.sql`
- [x] T040 Create `cabin_types` table migration in `backend/migrations/003_cabin_types.up.sql`
- [x] T041 Create `facility_categories` table migration in `backend/migrations/004_facility_categories.up.sql`
- [x] T042 Create `facilities` table migration in `backend/migrations/005_facilities.up.sql`
- [x] T043 [P] Implement CruiseCompany domain model in `backend/internal/domain/cruise_company.go`
- [x] T044 [P] Implement Cruise domain model in `backend/internal/domain/cruise.go`
- [x] T045 [P] Implement CabinType domain model in `backend/internal/domain/cabin_type.go`
- [x] T046 [P] Implement Facility domain models in `backend/internal/domain/facility.go`
- [x] T047 Implement Cruise repository with GORM in `backend/internal/repository/cruise.go`
- [x] T048 Implement CabinType repository in `backend/internal/repository/cabin_type.go`
- [x] T049 Implement Facility repository in `backend/internal/repository/facility.go`
- [x] T050 Implement Cruise service layer in `backend/internal/service/cruise.go`
- [x] T051 Implement CabinType service layer in `backend/internal/service/cabin_type.go`
- [x] T052 Implement Facility service layer in `backend/internal/service/facility.go`
- [x] T053 Implement Cruise handlers in `backend/internal/handler/cruise.go`
- [x] T054 Implement CabinType handlers in `backend/internal/handler/cabin_type.go`
- [x] T055 Implement Facility handlers in `backend/internal/handler/facility.go`
- [x] T056 Setup cruise routes in `backend/cmd/api/routes.go`
- [x] T057 Write unit tests for cruise service (100% coverage)
- [x] T058 Write integration tests for cruise handlers (100% coverage)

### Frontend - Customer Web

- [x] T059 [P] Create Home page layout in `frontend-web/pages/index.vue`
- [x] T060 Create Cruise list page with filters in `frontend-web/pages/cruises/index.vue`
- [x] T061 Create Cruise detail page in `frontend-web/pages/cruises/[id].vue`
- [x] T062 Create CruiseCard component in `frontend-web/components/cruise/CruiseCard.vue`
- [x] T063 Create ImageGallery component in `frontend-web/components/cruise/ImageGallery.vue`
- [x] T064 Create CabinTypeAccordion component in `frontend-web/components/cruise/CabinTypeAccordion.vue`
- [x] T065 Create FacilityTabs component in `frontend-web/components/cruise/FacilityTabs.vue`
- [x] T066 Write component tests for CruiseCard (100% coverage)
- [x] T067 Write Playwright E2E tests for cruise browsing (100% coverage)

### Frontend - Mini Program

- [x] T068 [P] Create home page in `frontend-mini/pages/index/index.vue`
- [x] T069 [P] Create cruise list page in `frontend-mini/pages/cruises/index.vue`
- [x] T070 [P] Create cruise detail page in `frontend-mini/pages/cruises/detail.vue`
- [x] T071 Create cruise card component in `frontend-mini/components/CruiseCard.vue`
- [x] T072 Write Jest tests for mini program components (100% coverage)

---

## Phase 4: US4 - Backend Management

**Goal**: Admin can manage cruise data (CRUD operations)  
**Priority**: P1 (Critical)  
**Duration**: 2 weeks (parallel with Phase 3)  
**Independent Test**: Admin can perform full CRUD on cruises, cabin types, facilities

### Backend - Admin APIs

- [x] T071 Implement Cruise CRUD handlers with image upload in `backend/internal/handler/admin_cruise.go`
- [x] T072 Implement CabinType CRUD handlers in `backend/internal/handler/admin_cabin_type.go`
- [x] T073 Implement Facility CRUD handlers in `backend/internal/handler/admin_facility.go`
- [x] T074 Implement MinIO file upload service in `backend/internal/service/storage.go`
- [x] T075 Setup admin routes with RBAC in `backend/cmd/api/admin_routes.go`
- [x] T076 Write admin handler tests (100% coverage)

### Frontend - Admin Panel

- [x] T077 Create AdminLayout with sidebar navigation in `frontend-admin/layouts/admin.vue`
- [x] T078 Create Login page in `frontend-admin/pages/login.vue`
- [x] T079 Create Cruise list page in `frontend-admin/pages/cruises/index.vue`
- [x] T080 Create Cruise create/edit form in `frontend-admin/pages/cruises/[id].vue`
- [x] T081 Integrate TipTap editor in `frontend-admin/components/TipTapEditor.vue`
- [x] T082 Create image upload component with drag-drop in `frontend-admin/components/ImageUpload.vue`
- [x] T083 Create CabinType management pages in `frontend-admin/pages/cabin-types/`
- [x] T084 Create Facility management pages in `frontend-admin/pages/facilities/`
- [x] T085 Write Vitest tests for admin components (100% coverage)
- [x] T086 Write Playwright E2E tests for admin flows (100% coverage)

---

## Phase 5: US2 - Online Booking & Payment

**Goal**: Complete booking flow with inventory locking and payment  
**Priority**: P1 (Critical)  
**Duration**: 2 weeks  
**Independent Test**: User can complete end-to-end booking with payment, inventory locks correctly

### Backend - Booking Domain

- [x] T087 Create `routes` table migration in `backend/migrations/006_routes.up.sql`
- [x] T088 Create `voyages` table migration in `backend/migrations/007_voyages.up.sql`
- [x] T089 Create `cabins` table migration in `backend/migrations/008_cabins.up.sql`
- [x] T090 Create `cabin_inventory` table migration in `backend/migrations/009_cabin_inventory.up.sql`
- [x] T091 Create `cabin_prices` table migration in `backend/migrations/010_cabin_prices.up.sql`
- [x] T092 [P] Implement Route, Voyage, Cabin, CabinInventory, CabinPrice domain models in `backend/internal/domain/booking.go`
- [x] T097 Implement Route repository in `backend/internal/repository/route.go`
- [x] T098 Implement Voyage repository in `backend/internal/repository/voyage.go`
- [x] T099 Implement Cabin repository in `backend/internal/repository/cabin.go`
- [x] T100 Implement Inventory repository with optimistic locking in `backend/internal/repository/inventory.go`
- [x] T101 Implement Price repository in `backend/internal/repository/price.go`
- [x] T102 Implement Inventory service with lock logic in `backend/internal/service/inventory.go`
- [x] T103 Implement Price service in `backend/internal/service/price.go`

### Backend - Order Domain

- [x] T104 Create `orders` table migration in `backend/migrations/011_orders.up.sql`
- [x] T105 Create `order_items` table migration in `backend/migrations/012_order_items.up.sql`
- [x] T106 Create `passengers` table migration in `backend/migrations/013_passengers.up.sql`
- [x] T107 Create `payments` table migration in `backend/migrations/014_payments.up.sql`
- [x] T108-111 [P] Implement Order, OrderItem, Passenger, Payment domain models in `backend/internal/domain/order.go`
- [x] T112 Implement Order repository in `backend/internal/repository/order.go`
- [x] T113 Implement Order state machine in `backend/internal/service/order_state.go`
- [x] T114 Implement Order service in `backend/internal/service/order.go`
- [x] T115 Implement Order handlers in `backend/internal/handler/order.go`

### Backend - Payment Integration

- [x] T116 Implement WeChat Pay V3 SDK integration in `backend/internal/payment/wechat.go`
- [x] T117 Implement payment callback handler with idempotency in `backend/internal/handler/payment.go`
- [x] T118 Implement order timeout job with NATS in `backend/internal/jobs/order_timeout.go`
- [x] T119 Write order service unit tests (100% coverage)
- [x] T120 Write payment integration tests (100% coverage)

### Frontend - Booking Flow

- [x] T121 Create booking wizard component in `frontend-web/components/booking/BookingWizard.vue`
- [x] T122 Create voyage selection step in `frontend-web/components/booking/SelectVoyage.vue`
- [x] T123 Create cabin selection step in `frontend-web/components/booking/SelectCabin.vue`
- [x] T124 Create passenger info step in `frontend-web/components/booking/PassengerForm.vue`
- [x] T125 Create order confirmation step in `frontend-web/components/booking/OrderConfirm.vue`
- [x] T126 Create payment page in `frontend-web/pages/payment/[orderId].vue`
- [x] T127 Create payment result page in `frontend-web/pages/payment/result.vue`
- [x] T128 Write Playwright E2E tests for booking flow (100% coverage)

### Frontend - Mini Program

- [x] T129 Create booking pages in `frontend-mini/pages/booking/`
- [x] T130 Create payment integration with WeChat SDK
- [x] T131 Write Jest tests for booking components (100% coverage)

---

## Phase 6: US3 - Order Management & Refunds

**Goal**: Order lifecycle management and refund processing  
**Priority**: P1 (Critical)  
**Duration**: 1 week  
**Independent Test**: Users can view orders and request refunds; admins can process them

### Backend - Order Management

- [x] T132 Create `refund_requests` table migration in `backend/migrations/015_refund_requests.up.sql`
- [x] T133 Implement RefundRequest domain model in `backend/internal/domain/refund.go`
- [x] T134 Implement order query endpoints in `backend/internal/handler/order_query.go`
- [x] T135 Implement refund workflow service in `backend/internal/service/refund.go`
- [x] T136 Implement admin order handlers in `backend/internal/handler/admin_order.go`
- [x] T137 Write refund service tests (100% coverage)

### Frontend - Customer

- [x] T138 Create My Orders page in `frontend-web/pages/orders/index.vue`
- [x] T139 Create Order detail page in `frontend-web/pages/orders/[id].vue`
- [x] T140 Create refund request form in `frontend-web/components/RefundRequest.vue`
- [x] T141 Write order management E2E tests (100% coverage)

### Frontend - Admin

- [x] T142 Create Order management list in `frontend-admin/pages/orders/index.vue`
- [x] T143 Create Order detail view in `frontend-admin/pages/orders/[id].vue`
- [x] T144 Create Refund approval interface in `frontend-admin/pages/refunds/index.vue`

---

## Phase 7: US5 - User Authentication & Account

**Goal**: Multi-platform authentication and user profiles  
**Priority**: P2 (Important)  
**Duration**: 1 week  
**Independent Test**: Users can login via WeChat/Phone and manage profiles

### Backend - User Domain

- [x] T145 Create `users` table migration in `backend/migrations/016_users.up.sql`
- [x] T146 Create `frequent_passengers` table migration in `backend/migrations/017_frequent_passengers.up.sql`
- [x] T147 [P] Implement User domain model in `backend/internal/domain/user.go`
- [x] T148 [P] Implement FrequentPassenger domain model in `backend/internal/domain/frequent_passenger.go`
- [x] T149 Implement User repository in `backend/internal/repository/user.go`
- [x] T150 Implement WeChat login service in `backend/internal/service/wechat_auth.go`
- [x] T151 Implement SMS verification service in `backend/internal/service/sms.go`
- [x] T152 Implement user handlers in `backend/internal/handler/user.go`
- [x] T153 Write user service tests (100% coverage)

### Frontend

- [x] T154 Create login modal in `frontend-web/components/LoginModal.vue`
- [x] T155 Create user profile page in `frontend-web/pages/profile/index.vue`
- [x] T156 Create frequent passengers management in `frontend-web/pages/profile/passengers.vue`
- [x] T157 Implement WeChat login in mini program
- [x] T158 Write auth flow E2E tests (100% coverage)

---

## Phase 8: US6 - Notifications & Reminders

**Goal**: Multi-channel notification system  
**Priority**: P2 (Important)  
**Duration**: 1 week  
**Independent Test**: Notifications sent on order events, inventory alerts work

### Backend

- [x] T159 Create `notifications` table migration in `backend/migrations/018_notifications.up.sql`
- [x] T160 Implement Notification domain model in `backend/internal/domain/notification.go`
- [x] T161 Implement notification service in `backend/internal/service/notification.go`
- [x] T162 Implement WeChat template message sender in `backend/internal/notification/wechat.go`
- [x] T163 Implement SMS sender in `backend/internal/notification/sms.go`
- [x] T164 Implement inventory alert job in `backend/internal/jobs/inventory_alert.go`
- [x] T165 Write notification tests (100% coverage)

### Frontend

- [x] T166 Create notification center in `frontend-web/components/NotificationCenter.vue`
- [x] T167 Create notification settings in `frontend-web/pages/profile/notifications.vue`

---

## Phase 9: US7 - Smart Recommendations & Analytics

**Goal**: AI recommendations and business analytics  
**Priority**: P3 (Enhancement)  
**Duration**: 2 weeks  
**Independent Test**: Recommendations displayed, analytics dashboards functional

### Backend

- [ ] T168 Implement user behavior tracking in `backend/internal/analytics/tracking.go`
- [ ] T169 Implement recommendation engine in `backend/internal/recommendation/engine.go`
- [ ] T170 Implement price trend analysis in `backend/internal/analytics/price_trends.go`
- [ ] T171 Create analytics API endpoints in `backend/internal/handler/analytics.go`
- [ ] T172 Write recommendation tests (100% coverage)

### Frontend

- [ ] T173 Create recommendation carousel in `frontend-web/components/RecommendationCarousel.vue`
- [ ] T174 Create price calendar view in `frontend-web/components/PriceCalendar.vue`
- [ ] T175 Create cabin comparison tool in `frontend-web/components/CabinComparison.vue`
- [ ] T176 Create analytics dashboard in `frontend-admin/pages/dashboard.vue`
- [ ] T177 Write E2E tests for smart features (100% coverage)

---

## Phase 10: US8 - Social Sharing & Community

**Goal**: Social features and community engagement  
**Priority**: P3 (Enhancement)  
**Duration**: 2 weeks  
**Independent Test**: Users can share, review, and post travelogues

### Backend

- [ ] T178 Create `travelogues` table migration in `backend/migrations/019_travelogues.up.sql`
- [ ] T179 Implement Travelogue domain model in `backend/internal/domain/travelogue.go`
- [ ] T180 Implement review system in `backend/internal/service/review.go`
- [ ] T181 Implement poster generation service in `backend/internal/service/poster.go`
- [ ] T182 Implement invitation system in `backend/internal/service/invitation.go`
- [ ] T183 Write social feature tests (100% coverage)

### Frontend

- [ ] T184 Create review form in `frontend-web/components/ReviewForm.vue`
- [ ] T185 Create travelogue editor in `frontend-web/components/TravelogueEditor.vue`
- [ ] T186 Create community page in `frontend-web/pages/community/index.vue`
- [ ] T187 Create share poster modal in `frontend-web/components/SharePoster.vue`
- [ ] T188 Write social features E2E tests (100% coverage)

---

## Phase 11: Polish & Production

**Goal**: Production-ready system with monitoring and optimization  
**Duration**: 1 week

### Performance & Optimization

- [ ] T189 Implement Redis caching layer in `backend/internal/cache/`
- [ ] T190 Setup Meilisearch indexing for cruises and cabins
- [ ] T191 Implement database query optimization
- [ ] T192 Add API rate limiting middleware
- [ ] T193 Implement frontend code splitting and lazy loading

### Monitoring & Observability

- [ ] T194 Setup Prometheus metrics collection
- [ ] T195 Create Grafana dashboards
- [ ] T196 Configure Loki log aggregation
- [ ] T197 Implement distributed tracing
- [ ] T198 Setup health check endpoints

### Security & Compliance

- [ ] T199 Implement API request signing
- [ ] T200 Add SQL injection prevention tests
- [ ] T201 Conduct security audit
- [ ] T202 Implement GDPR data export/deletion

### Documentation

- [ ] T203 Complete API documentation with examples
- [ ] T204 Create deployment guide
- [ ] T205 Write operation runbook
- [ ] T206 Create user manual

---

## Dependencies Graph

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundation)
    ↓
┌─────────────────────────────────────────────────────────────┐
│  Phase 3 (US1: Browse) ──────┐                              │
│  Phase 4 (US4: Admin) ───────┤ Parallel                     │
└─────────────────────────────────────────────────────────────┘
    ↓ (Both complete)
Phase 5 (US2: Booking) ──→ Phase 6 (US3: Orders)
    ↓
Phase 7 (US5: Auth) ─────→ Phase 8 (US6: Notifications)
    ↓
Phase 9 (US7: Smart) ────→ Phase 10 (US8: Social)
    ↓
Phase 11 (Polish)
```

## Parallel Execution Opportunities

### Maximum Parallelism Per Phase

**Phase 1**: All 8 backend + 8 frontend tasks can run in parallel (T001-T016)

**Phase 2**: Authentication and infrastructure tasks can parallelize (T021-T036)

**Phase 3 & 4**: These two phases can run simultaneously:
- Phase 3 focuses on customer-facing browsing
- Phase 4 focuses on admin management
- Both need backend models but can be developed in parallel

**Phase 5+**: Sequential dependency on booking flow completion

### Team Parallelization Strategy

**Team A (Backend)**: Focus on domain models and business logic
**Team B (Frontend Web)**: Focus on customer web experience  
**Team C (Frontend Admin)**: Focus on admin panel
**Team D (Mini Program)**: Focus on WeChat mini program

Each team can work on their respective tasks within each phase.

---

## MVP Scope Recommendation

**MVP = Phase 1-6** (First 8 weeks)

This includes:
- ✓ Project setup and infrastructure
- ✓ Cruise browsing (US1)
- ✓ Admin management (US4)
- ✓ Booking and payment (US2)
- ✓ Order management (US3)
- ✓ Basic user auth (US5 - core only)

**Excluded from MVP**:
- Smart recommendations (US7)
- Social features (US8)
- Advanced notifications (US6 - basic only)
- Full analytics dashboard

**MVP Deliverable**: A fully functional cruise booking platform where users can browse, book, pay, and manage orders.

---

## Testing Strategy

Per Constitution v1.0.0 requirements:

**Backend (Go)**:
- Unit tests: testify + gomock
- Integration tests: httptest + testcontainers-go
- E2E tests: Full API workflows
- **Coverage Requirement**: 100%

**Frontend (Nuxt)**:
- Unit tests: Vitest + Vue Test Utils
- E2E tests: Playwright
- **Coverage Requirement**: 100%

**Mini Program**:
- Unit tests: Jest + @vue/test-utils
- Component tests: miniprogram-simulate
- **Coverage Requirement**: 100%

---

## Success Criteria Per Phase

| Phase | Success Criteria |
|-------|-----------------|
| 1 | All projects build without errors, CI pipeline passes |
| 2 | Auth middleware working, tests passing |
| 3 | User can browse cruises end-to-end, all tests pass |
| 4 | Admin can CRUD all cruise data, tests pass |
| 5 | Complete booking flow works, payment processes, tests pass |
| 6 | Order lifecycle complete, refunds work, tests pass |
| 7 | Multi-platform login works, profiles functional |
| 8 | Notifications sent correctly |
| 9 | Recommendations generated, analytics visible |
| 10 | Social features functional |
| 11 | System production-ready, monitoring active |

---

**End of Tasks Document**

*Total Tasks: 206 | Estimated Effort: 16 Sprints (9 Months) | Team Size: 4-6 developers*
