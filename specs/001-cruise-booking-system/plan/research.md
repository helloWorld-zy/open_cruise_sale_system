# Research Document: CruiseBooking Platform

**Feature**: 邮轮舱位预订平台 (CruiseBooking)  
**Date**: 2026-02-10  
**Status**: Complete  

---

## Technical Context

### Architecture Overview

This is a full-stack e-commerce platform for cruise cabin sales with the following architecture:

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend Layer                            │
├──────────────────┬──────────────────┬──────────────────────────┤
│  Management      │  Customer Web    │  Customer Mini Program   │
│  (Nuxt 4 SSR)    │  (Nuxt 4 SSR)    │  (uni-app Vue 3)         │
│  - Admin Dashboard│  - Public Site   │  - WeChat Mini Program   │
└──────────────────┴──────────────────┴──────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway (K8s Ingress)                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Backend Layer (Go)                        │
│  ┌──────────────┬──────────────┬──────────────┬──────────────┐  │
│  │   Cruise     │    Cabin     │    Order     │    User      │  │
│  │   Service    │   Service    │   Service    │   Service    │  │
│  └──────────────┴──────────────┴──────────────┴──────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Data & Infrastructure                       │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────────┐   │
│  │PostgreSQL│  Redis   │ Meilisearch│ MinIO  │NATS JetStream│   │
│  │  (Main)  │ (Cache)  │ (Search)  │(Storage)│  (Queue)     │   │
│  └──────────┴──────────┴──────────┴──────────┴──────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### Technology Stack Decisions

All technology choices are mandated by Constitution v1.0.0 and are non-negotiable.

#### Backend Stack (Go)

| Component | Technology | Version | Rationale |
|-----------|------------|---------|-----------|
| Language | Go | 1.26 | High performance, excellent concurrency support for inventory management |
| Web Framework | Gin | v1.11.0 | Fast HTTP router, rich middleware ecosystem |
| ORM | GORM | v2.x | Auto-migration, hooks, preload support for complex relationships |
| Database | PostgreSQL | 17.x | ACID compliance, JSON support, excellent for financial data |
| Cache | Redis | 7.4.x | Distributed locking for inventory, session storage |
| Search | Meilisearch | 1.12.x | Lightweight alternative to ES, fast full-text search |
| Message Queue | NATS JetStream | 2.11.x | Event-driven architecture, delayed messages for order timeout |
| Object Storage | MinIO | Latest | S3-compatible, self-hosted media storage |
| Auth | JWT + Casbin | — | Industry-standard RBAC implementation |

#### Frontend Stack

**Management & Customer Web (Nuxt 4)**

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Framework | Nuxt 4.3.0 | Full-stack Vue framework, SSR for SEO |
| UI Library | Nuxt UI v3 | Tailwind-based, consistent design system |
| State | Pinia v3 | Vue official state management |
| Forms | VeeValidate + Zod | Type-safe form validation |
| Rich Text | TipTap v2 | Extensible editor for cruise/cabin descriptions |
| Charts | ECharts v5.6 | Industry standard for data visualization |

**Mini Program (uni-app)**
| Component | Technology | Rationale |
|-----------|------------|-----------|
| Framework | uni-app (Vue 3) | Compile to WeChat Mini Program |
| State | Pinia v3 | Shared state management pattern |
| Testing | miniprogram-simulate | WeChat official testing tool |

### Critical Technical Challenges Resolved

#### 1. Concurrent Inventory Management

**Challenge**: Preventing overselling when multiple users book the same cabin simultaneously.

**Decision**: Use Redis distributed locks with the following strategy:

```
1. User submits order → Try to acquire Redis lock (cabin_id + voyage_id)
2. Lock acquired → Check inventory → Deduct if available
3. Create order with 15-min expiration
4. Release lock after order creation
5. If payment timeout → Release inventory + delete lock
```

**Alternative Considered**: 
- Database pessimistic locking (too slow, causes contention)
- Optimistic locking with versioning (complex retry logic)
- **Selected**: Redis distributed locks (best balance of performance and reliability)

#### 2. Payment Callback Reliability

**Challenge**: Ensuring order state consistency when payment callbacks are delayed or lost.

**Decision**: Implement idempotent payment processing with dual-write strategy:

```
1. User pays → Third-party gateway
2. Gateway callback → API endpoint
3. Idempotency check (payment_idempotency_key)
4. Update order status + create payment record (transaction)
5. Background job: Periodic reconciliation + compensation
6. Query API: Frontend can actively query payment status
```

**Alternative Considered**:
- Pure async processing (too slow for UX)
- Synchronous waiting (blocks resources)
- **Selected**: Async callback + active query + reconciliation (eventual consistency)

#### 3. Multi-Frontend Code Sharing

**Challenge**: Maximizing code reuse between Web and Mini Program while maintaining platform-specific optimizations.

**Decision**: 
- API layer: Shared OpenAPI spec generates clients for all platforms
- UI components: Platform-specific implementations (different design systems)
- Business logic: Shared TypeScript types via monorepo
- State management: Same Pinia patterns, adapted for each platform

**Alternative Considered**:
- Taro (React-based, not Vue-native)
- Uni-app for all platforms (limited Web flexibility)
- **Selected**: Separate Nuxt for Web + uni-app for Mini Program (best UX per platform)

#### 4. Full Test Coverage Strategy

**Challenge**: Achieving 100% test coverage across 4 codebases (Go backend + 3 frontends) without slowing development.

**Decision**: 
- Backend: Table-driven tests for handlers, mock repositories with gomock
- Frontend: Component tests with Vitest/Vue Test Utils, E2E with Playwright (critical flows)
- Mini Program: Component tests with jest, E2E with WeChat DevTools automation
- CI: Parallel test execution, coverage gates at 100%

**Alternative Considered**:
- Lower coverage threshold (violates Constitution)
- Manual testing only (not scalable)
- **Selected**: Automated testing at all levels with strict CI gates

---

## Database Schema Research

### Key Design Decisions

#### 1. Inventory Tracking Strategy

**Approach**: Separate `cabin_inventory` table with real-time calculations

```sql
-- Available inventory = total - sold - locked
-- Locked inventory tracked separately for unpaid orders
-- Inventory log for audit trail
```

**Rationale**: 
- Avoids row-level locking on main cabin table
- Supports inventory history and reconciliation
- Enables complex queries (low stock alerts, etc.)

#### 2. Pricing Model

**Approach**: Flexible price matrix with date-range support

```sql
-- Base price per cabin per voyage per date
-- Override prices for: child, single supplement, holidays, early bird
-- Currency support for multi-currency (V1.5)
```

**Rationale**:
- Cruise pricing is highly variable by date
- Supports promotional pricing without code changes
- Efficient querying for calendar views

#### 3. Order State Machine

**Approach**: Explicit state transitions with event sourcing

```
created → pending_payment → paid → confirmed → pending_departure → departed → completed
   ↓           ↓              ↓
cancelled   timeout      refund_requested → refund_processing → refunded
```

**Rationale**:
- Clear business rules per state
- Audit trail for all transitions
- Supports complex refund workflows

---

## Integration Research

### Third-Party Services

| Service | Purpose | Integration Pattern | Risk Mitigation |
|---------|---------|---------------------|-----------------|
| WeChat Pay | Payment (Mini Program) | REST API + Webhooks | Webhook retry + reconciliation |
| Alipay | Payment (Web) | REST API + Webhooks | Same as above |
| Meilisearch | Search | Go client library | Local cache fallback |
| MinIO | File storage | S3-compatible SDK | Multi-region replication |
| NATS | Messaging | Go NATS client | Cluster deployment |
| Weather API | Port weather | REST API | Cached responses, graceful degradation |
| OCR Service | ID document recognition | REST API | Client-side validation fallback |

### Internal Service Communication

**Pattern**: REST APIs for synchronous, NATS for asynchronous

- **Synchronous**: User auth, CRUD operations, queries
- **Asynchronous**: Order timeout handling, notification sending, report generation

---

## Performance Considerations

### Caching Strategy

| Data Type | Cache Layer | TTL | Invalidation |
|-----------|-------------|-----|--------------|
| Cruise basic info | Redis | 1 hour | On update |
| Cabin availability | Redis | 5 minutes | On booking |
| Price data | Redis | 15 minutes | On price change |
| User session | Redis | 24 hours | On logout |
| Search results | Meilisearch | Real-time | Index update |
| Static assets | CDN/MinIO | Long-term | Versioned URLs |

### Scaling Strategy

- **Horizontal**: Kubernetes HPA based on CPU/memory
- **Database**: Read replicas for queries, connection pooling
- **Cache**: Redis Cluster for high availability
- **Search**: Meilisearch horizontal scaling

---

## Security Research

### Authentication Flow

1. **Mini Program**: WeChat OAuth → Code → OpenID/UnionID → JWT token
2. **Web WeChat**: OAuth redirect → Code → User info → JWT token
3. **Web Mobile**: SMS验证码 → Phone number → JWT token

### Authorization (RBAC)

**Roles**: Super Admin > Operations > Finance > Customer Service

**Permissions managed via Casbin policies**:
```
p, super_admin, *, *
p, operations, cruise, (read|write)
p, finance, order_refund, (read|write)
```

### Data Protection

- **Encryption at rest**: PostgreSQL transparent data encryption
- **Encryption in transit**: TLS 1.3 for all communications
- **PII handling**: Masked in logs, encrypted in database
- **Payment data**: Never stored, only tokens

---

## Development Workflow

### Code Organization

```
/
├── backend/                    # Go monolith (can split later)
│   ├── cmd/api/               # Main application
│   ├── internal/
│   │   ├── domain/            # Business entities
│   │   ├── service/           # Business logic
│   │   ├── repository/        # Data access
│   │   ├── handler/           # HTTP handlers
│   │   └── middleware/        # Auth, logging, etc.
│   ├── pkg/                   # Shared packages
│   └── migrations/            # Database migrations
│
├── frontend-admin/            # Nuxt 4 - Management
├── frontend-web/              # Nuxt 4 - Customer Web
├── frontend-mini/             # uni-app - Mini Program
│
├── specs/                     # This directory
└── .github/workflows/         # CI/CD pipelines
```

### Testing Strategy

**Backend (Go)**:
- Unit: `*_test.go` alongside source files
- Integration: `tests/integration/*_test.go`
- E2E: `tests/e2e/*_test.go`
- Tools: testify, gomock, httptest, testcontainers-go

**Frontend (Nuxt)**:
- Unit: `*.spec.ts` with Vitest
- Component: Vue Test Utils
- E2E: Playwright tests in `tests/e2e/`

**Mini Program (uni-app)**:
- Unit: jest + @vue/test-utils
- Component: miniprogram-simulate

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Third-party payment API changes | Medium | High | Abstract payment interface, multi-provider support |
| WeChat Mini Program policy changes | Low | High | Stay updated, have Web fallback |
| Performance under high load | Medium | Medium | Load testing, auto-scaling, caching |
| Data inconsistency in distributed locks | Low | High | Reconciliation jobs, audit logs |
| Test coverage enforcement slowing delivery | Medium | Medium | TDD discipline, test scaffolding tools |

---

## Open Questions Resolved

### Q1: Should we use microservices or monolith?

**Decision**: Start with modular monolith, split when needed.

**Rationale**: 
- Team size is small (typical for MVP)
- Simpler deployment and testing
- Clear module boundaries for future extraction
- Constitution allows either, but simplicity is preferred

### Q2: How to handle timezone issues?

**Decision**: Store all dates in UTC, convert to local timezone for display.

**Rationale**:
- Cruises span multiple timezones
- UTC is unambiguous for server-side logic
- Frontend handles user-local display

### Q3: Soft delete or hard delete?

**Decision**: Soft delete for all business entities (Cruise, Cabin, Order, etc.)

**Rationale**:
- Auditing requirements
- Referential integrity with historical data
- Can hard delete after retention period if needed

---

## Conclusion

All major technical decisions have been researched and documented. The architecture follows industry best practices while adhering to Constitution constraints. The modular design allows for iterative development across the 4 delivery phases (MVP → V1.0 → V1.5 → V2.0).

**Next Step**: Proceed to Phase 1 - Data Model and API Contract Design.

---

**End of Research Document**
