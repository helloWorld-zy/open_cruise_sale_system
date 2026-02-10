<!--
Sync Impact Report
===================
Version Change: v1.0.0 → v1.0.0 (Initial Creation)
Modified Principles: None (initial creation)
Added Sections:
  - 1. Project Overview
  - 2. Technical Stack Charter
  - 3. Testing Mandate
  - 4. Quality Assurance Principles
  - 5. Development Workflow
  - 6. Architecture Principles
  - 7. Governance
Removed Sections: None
Templates Requiring Updates: None (initial creation)
Follow-up TODOs: None
-->

# Constitution: 邮轮舱位预订平台 (CruiseBooking)

**Version:** v1.0.0  
**Ratification Date:** 2026-02-10  
**Last Amended Date:** 2026-02-10  

---

## 1. Project Overview

### 1.1 Identity

**Project Name:** 邮轮舱位预订平台 (CruiseBooking)  
**Type:** Full-stack e-commerce platform for cruise cabin sales  
**Scope:** End-to-end booking experience from discovery to post-trip engagement

### 1.2 Delivery Phases

| Phase | Timeline | Description |
|-------|----------|-------------|
| **MVP (Category II)** | Months 1-2 | Core booking flow: cruise display, cabin browsing, booking, payment |
| **V1.0** | Months 3-4 | Shore excursions, e-tickets, itinerary countdown, reviews |
| **V1.5** | Months 5-6 | Smart recommendations, VR preview, dynamic pricing, multi-currency |
| **V2.0** | Months 7-9 | UGC community, group booking, distribution channels, multi-language |

### 1.3 Frontend Applications

1. **Management Backend (Web)** - Nuxt 4 + Vue 3 + Nuxt UI v3
2. **Customer Frontend (Web)** - Nuxt 4 SSR for SEO
3. **Customer Frontend (Mini Program)** - uni-app Vue 3

---

## 2. Technical Stack Charter

### 2.1 Backend (Go Stack) - NON-NEGOTIABLE

| Component | Technology | Version | Purpose |
|-----------|------------|---------|---------|
| Language | Go | 1.26 | Primary backend language |
| Web Framework | Gin | v1.11.0 | HTTP routing and middleware |
| ORM | GORM | v2.x | Database abstraction |
| Database | PostgreSQL | 17.x | Primary data store |
| Cache | Redis | 7.4.x | Sessions, inventory locks, leaderboards |
| Search | Meilisearch | 1.12.x | Full-text cabin/route search |
| Message Queue | NATS JetStream | 2.11.x | Async processing, event-driven |
| Object Storage | MinIO | Latest | Media assets (S3-compatible) |
| API Spec | RESTful + OpenAPI 3.1 | — | Mandatory Swagger documentation |
| WebSocket | gorilla/websocket | Latest | Real-time cabin status |
| Auth/RBAC | JWT + Casbin | — | Role-based access control |
| Config | Viper | v1.19.x | Multi-environment config |
| Logging | Zap + Lumberjack | Latest | Structured logging with rotation |
| Payment | WeChat Pay V3 / Alipay | Latest | Online payments |
| Container | Kubernetes | Latest | Deployment orchestration |
| CI/CD | GitHub Actions | — | Automated build, test, deploy |
| Testing | Go testing + testify + gomock + httptest | — | **100% coverage required** |

**Principle T2.1:** All backend services MUST use the specified technology stack. No exceptions without explicit architectural review and constitution amendment.

### 2.2 Frontend Stack - NON-NEGOTIABLE

**Management & Web Frontend:**

| Component | Technology | Version |
|-----------|------------|---------|
| Framework | Nuxt | 4.3.0 |
| UI Core | Vue | 3.5.28 |
| Build Tool | Vite | 7.3.1 |
| UI Library | Nuxt UI v3 | Latest |
| CSS | Tailwind CSS | v4.x |
| State | Pinia | v3.x |
| HTTP | ofetch / useFetch | Built-in |
| Charts | ECharts | v5.6.x |
| Rich Text | TipTap | v2.x |
| Forms | VeeValidate + Zod | Latest |
| Language | TypeScript | 5.9.x (Strict Mode) |
| Testing | Vitest + Vue Test Utils + Playwright | **100% coverage** |

**Mini Program Frontend:**

| Component | Technology | Version |
|-----------|------------|---------|
| Framework | uni-app | HBuilderX (Vue 3 mode) |
| UI Core | Vue | 3.5.28 |
| UI Library | uni-ui + Custom | Latest |
| CSS | rpx + scss | — |
| State | Pinia | v3.x |
| HTTP | uni.request / luch-request | Latest |
| Testing | jest + @vue/test-utils + miniprogram-simulate | **100% coverage** |

**Principle T2.2:** All frontend applications MUST use the specified stack and achieve 100% test coverage. TypeScript strict mode is mandatory.

---

## 3. Testing Mandate - ABSOLUTE REQUIREMENT

### 3.1 Coverage Requirements

**ALL code (backend + all frontends) MUST achieve 100% test coverage. NO EXCEPTIONS.**

### 3.2 Backend Testing Discipline

- **Unit Tests:** Every handler, service, repository, and middleware MUST have unit tests using testify + gomock
- **Integration Tests:** Critical business flows (booking, payment, refunds) MUST have integration tests using httptest + testcontainers-go
- **E2E Tests:** All API endpoints MUST have end-to-end test coverage

### 3.3 Frontend Testing Discipline

**Nuxt Applications:**
- All components MUST have Vitest + Vue Test Utils unit tests
- All pages MUST have Playwright E2E tests
- Nuxt server routes MUST have Vitest integration tests

**Mini Program:**
- All components MUST have jest + @vue/test-utils unit tests
- Component behavior MUST be tested with miniprogram-simulate
- Critical flows MUST have WeChat DevTools automated tests

### 3.4 CI Enforcement

**Principle T3.1:** GitHub Actions MUST enforce 100% coverage threshold. Pull requests failing coverage checks CANNOT be merged.

---

## 4. Quality Assurance Principles

### 4.1 API Documentation

**Principle Q4.1:** Every API endpoint MUST have corresponding Swagger/OpenAPI 3.1 documentation. CI MUST validate documentation-code consistency.

### 4.2 Code Review

**Principle Q4.2:** ALL code MUST pass pull request review before merging to main branch. No direct commits to main.

### 4.3 Database Changes

**Principle Q4.3:** ALL database schema changes MUST use migration files. Manual database modifications are STRICTLY PROHIBITED.

### 4.4 Code Standards

**Principle Q4.4:** ALL code MUST pass lint checks in CI:
- Backend: golangci-lint
- Frontend: ESLint + Prettier

---

## 5. Development Workflow

### 5.1 Sprint Structure

- **Sprint Duration:** 2 weeks
- **Sprint Goals:** Each sprint MUST deliver tested, documented, deployable features
- **Sprint Closure:** 100% coverage verification mandatory before sprint completion

### 5.2 CI/CD Pipeline

The following pipeline MUST be enforced for ALL changes:

```
Lint Check → Unit Tests → Integration Tests → E2E Tests 
→ Coverage Threshold (100%) → Build → Deploy
```

### 5.3 Feature Delivery Checklist

Before any feature is considered complete:

- [ ] Implementation follows technology stack charter
- [ ] Unit tests written with 100% coverage
- [ ] Integration/E2E tests written where applicable
- [ ] API documentation updated (if applicable)
- [ ] Database migrations created (if applicable)
- [ ] Code review approved
- [ ] CI pipeline passing

---

## 6. Architecture Principles

### 6.1 Service Boundaries

**Principle A6.1:** The backend MUST expose RESTful APIs consumed by all three frontend applications. No direct database access from frontends.

### 6.2 Authentication & Authorization

**Principle A6.2:** ALL protected endpoints MUST use JWT authentication with Casbin RBAC authorization. Role hierarchy: Super Admin → Operations → Finance → Customer Service.

### 6.3 Inventory Management

**Principle A6.3:** Cabin inventory MUST use Redis distributed locks for concurrent booking protection. Lock timeout: 15 minutes for unpaid orders.

### 6.4 Payment Security

**Principle A6.4:** ALL payment callbacks MUST validate signatures. Order amounts MUST be server-side verified against payment amounts. Tampering attempts MUST be logged and rejected.

### 6.5 Observability

**Principle A6.5:** ALL services MUST have:
- Prometheus metrics collection
- Grafana dashboards
- Loki centralized logging
- Structured logs with trace IDs

---

## 7. Governance

### 7.1 Amendment Procedure

Constitution amendments require:

1. **Minor Changes** (clarifications, wording): PR review + approval from 2 maintainers
2. **Major Changes** (principle additions/removals, stack changes): Team discussion + unanimous agreement + version bump

### 7.2 Versioning Policy

Constitution versions follow semantic versioning:

- **MAJOR (X.0.0):** Backward incompatible governance changes, principle removals, stack redefinitions
- **MINOR (x.Y.0):** New principles added, sections materially expanded
- **PATCH (x.y.Z):** Clarifications, wording improvements, typo fixes

### 7.3 Compliance Review

- **Quarterly Review:** All principles reviewed for relevance and adherence
- **Sprint Audit:** Random compliance checks on test coverage and documentation
- **Violation Escalation:** Repeated violations require team retrospective and process improvement

### 7.4 Technical Debt

**Principle G7.1:** Technical debt items MUST be tracked explicitly. Debt accumulation without remediation plan is NOT permitted across sprint boundaries.

---

## 8. Summary of Non-Negotiables

| ID | Principle | Violation Consequence |
|----|-----------|---------------------|
| T2.1 | Technology stack compliance | Architecture review + potential rewrite |
| T2.2 | Frontend stack + TS strict mode | PR rejection |
| T3.1 | 100% test coverage | PR cannot merge |
| Q4.1 | API documentation | PR rejection |
| Q4.2 | Code review requirement | Revert + disciplinary |
| Q4.3 | Migration-only DB changes | Database recovery + process review |
| Q4.4 | Lint compliance | Auto-fix or PR rejection |
| A6.2 | Auth/Authz enforcement | Security audit + immediate fix |
| A6.4 | Payment validation | Security incident response |

---

**End of Constitution v1.0.0**

*This document defines the foundational rules and principles for the CruiseBooking project. All team members are expected to understand and adhere to these guidelines.*
