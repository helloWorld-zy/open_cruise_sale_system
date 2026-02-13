# 技术债务审查报告 (REVIEW.md)

**项目**: 邮轮舱位预订平台 (CruiseBooking)
**审查日期**: 2026-02-13
**审查范围**: 全栈（Backend Go + Frontend Vue/TypeScript）

---

## 摘要

| 维度 | 严重等级 | 问题数 |
|------|---------|-------|
| 🔴 安全债务 | **Critical** | 7 |
| 🔴 代码异味 | **High** | 8 |
| 🟠 设计债务 | **Medium-High** | 7 |
| 🟡 测试债务 | **High** | 6 |
| 🟡 文档债务 | **Medium** | 4 |

---

## 1. 安全债务 🔴

### SEC-001 [CRITICAL] 认证系统完全使用 Mock 数据

**文件**: `backend/internal/handler/auth.go`
**行号**: L56-L89, L101-L126, L137-L139, L174-L183, L200-L222

Login、Refresh、Logout、ChangePassword、WeChatLogin 全部使用 TODO 占位及硬编码 mock 响应，**未接入数据库验证**。任何用户名/密码组合均可登录并获取有效 JWT token。

```go
// L63-64: Login 函数
// TODO: Validate credentials against database
// For now, return mock response

// L67: 硬编码用户ID
token, err := middleware.GenerateToken("user-id", req.Username, "operations", h.jwtConfig)
```

**风险**: 生产环境零认证保护，任何人可获取管理后台完整权限。
**建议**: 实现基于数据库的用户验证、bcrypt 密码校验、refresh token 存储与轮换。

---

### SEC-002 [CRITICAL] JWT Token 无黑名单/撤销机制

**文件**: `backend/internal/handler/auth.go` L137-L139

```go
func (h *AuthHandler) Logout(c *gin.Context) {
    // TODO: Add token to blacklist
    response.Success(c, gin.H{"message": "Logged out successfully"})
}
```

用户登出后 token 仍然有效直到过期，无法主动撤销已泄露的 token。

**建议**: 使用 Redis 维护 token 黑名单，或改用短有效期 access token + refresh token 轮换方案。

---

### SEC-003 [CRITICAL] 微信支付 AES-GCM 解密未实现

**文件**: `backend/internal/payment/wechat.go` L415-L429

```go
func (w *wechatPay) decrypt(ciphertext, associatedData, nonce string) ([]byte, error) {
    // This is a simplified implementation - real implementation should use proper AES-GCM
    _ = associatedData
    _ = nonce
    _ = cipherBytes
    return nil, errors.New("decryption not implemented")
}
```

支付回调通知解密直接返回错误，意味着 **所有微信支付回调均会失败**。

**建议**: 使用 `crypto/aes` + `crypto/cipher` 实现标准 AES-256-GCM 解密。

---

### SEC-004 [HIGH] 支付回调缺少幂等性保护

**文件**: `backend/internal/payment/service.go` L88-L127

spec 中明确要求支付回调幂等性（FR-035A/B/C），但当前实现：
- 无幂等性键生成与存储
- 无重复支付检测
- 无重复回调防护

**建议**: 按照 spec 规范在 Redis 中实现幂等性键（`payment:idempotent:{order_id}`，TTL=24h），处理前先检查。

---

### SEC-005 [HIGH] `.env.example` 包含可推测的默认凭据

**文件**: `backend/.env.example` L17, L29-L30, L38

```
DB_PASSWORD=cruisebooking_secret
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

默认密码过于可猜测，且 JWT 密钥虽标注了 "change-this"，但生产环境可能被直接拷贝使用。

**建议**: `.env.example` 中密码/密钥值留空或使用 `<CHANGE_ME>` 占位符；增加启动配置校验拒绝默认值。

---

### SEC-006 [HIGH] 支付处理器 Handler 中变量遮蔽

**文件**: `backend/internal/handler/payment.go` L100-L114

```go
func (h *PaymentHandler) Query(c *gin.Context) {
    id := c.Param("id")
    payment, err := h.service.QueryPayment(c.Request.Context(), id) // payment 遮蔽了包名
    if err != nil {
        if err == payment.ErrPaymentNotFound {  // 此处 payment 已经是变量而非包名，编译可能出错或逻辑错误
```

局部变量 `payment` 遮蔽了导入的 `payment` 包，导致 `payment.ErrPaymentNotFound` 引用的是返回的对象而非包级错误常量。

**建议**: 重命名局部变量为 `result` 或 `paymentResult`，避免遮蔽包名。

---

### SEC-007 [MEDIUM] 签名验证未校验 JWT 签名算法

**文件**: `backend/internal/middleware/jwt.go` L43-L45

```go
token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(cfg.Secret), nil
})
```

未验证 `token.Method` 是否为预期的 `HS256`，存在算法混淆攻击风险（如伪造 `none` 算法）。

**建议**: 添加 `if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { return nil, errors.New("unexpected signing method") }`。

---

## 2. 代码异味 🔴

### CS-001 [CRITICAL] DSN 生成存在严重 Bug

**文件**: `backend/internal/config/config.go` L115-L122

```go
func (c *DatabaseConfig) GetDSN() string {
    return "host=" + c.Host +
        " port=" + string(rune(c.Port)) +   // ❌ 将 int 转为 rune 再转 string，得到的是 Unicode 字符而非数字字符串
        " user=" + c.User +
        " password=" + c.Password +
        " dbname=" + c.Name +
        " sslmode=" + c.SSLMode
}
```

`string(rune(5432))` 不会得到 `"5432"`，而是 Unicode 码位 5432 对应的字符。数据库连接将 **必然失败**。

**建议**: 使用 `fmt.Sprintf("%d", c.Port)` 或 `strconv.Itoa(c.Port)`。

---

### CS-002 [HIGH] 金额计算使用 `float64`

**文件**: 涉及多个文件
- `backend/internal/service/order.go` L198, L251-L256（价格计算）
- `backend/internal/payment/service.go` L26（退款金额参数）
- `backend/internal/payment/wechat.go` L148, L268, L308-L309（元/分转换）
- `backend/internal/domain/order.go`（TotalAmount 等字段）

浮点数精度问题在金融计算中会导致舍入错误。例如 `0.1 + 0.2 != 0.3`。

**建议**: 使用 `shopspring/decimal` 库或改为以 **分** 为单位的 `int64` 存储金额。

---

### CS-003 [HIGH] 价格计算逻辑重复

**文件**: `backend/internal/service/order.go`
- `Create` 方法 L251-L256
- `CalculateTotal` 方法 L480-L485

完全相同的价格计算逻辑在两个方法中重复出现：

```go
adultTotal := price.AdultPrice * float64(itemReq.AdultCount)
childTotal := price.ChildPrice * float64(itemReq.ChildCount)
infantTotal := price.InfantPrice * float64(itemReq.InfantCount)
portFee := price.PortFee * float64(itemReq.AdultCount+itemReq.ChildCount)
serviceFee := price.ServiceFee * float64(itemReq.AdultCount+itemReq.ChildCount)
subtotal := adultTotal + childTotal + infantTotal + portFee + serviceFee
```

**建议**: 抽取为独立的 `calculateItemSubtotal(price, item)` 方法。

---

### CS-004 [HIGH] `getCurrentTimestamp()` 函数返回空字符串

**文件**: `backend/internal/repository/order.go` L436-L438

```go
func getCurrentTimestamp() string {
    return "" // This should be implemented properly in real code
}
```

该函数被 `UpdateStatus` 和 `UpdatePaymentStatus` 调用，导致 `confirmed_at`、`cancelled_at`、`paid_at` 等时间戳字段始终为空。

**建议**: `return time.Now().UTC().Format(time.RFC3339)`。

---

### CS-005 [HIGH] `IsExpired()` 方法始终返回 false

**文件**: `backend/internal/service/order_state.go` L307-L311

```go
func (o *domain.Order) IsExpired() bool {
    // This should be implemented on the domain.Order struct
    // For now, return false as placeholder
    return false
}
```

订单永远不会被判定为过期，支付超时机制完全失效。

**建议**: 实现实际的过期判断逻辑（比较 `ExpiresAt` 与当前时间）。

---

### CS-006 [MEDIUM] 使用 `fmt.Printf` 而非结构化日志

**文件**: `backend/internal/service/order_state.go` L274, L278

```go
fmt.Printf("Failed to cancel booking for cabin %s: %v\n", item.CabinTypeID, err)
fmt.Printf("Failed to unlock cabin %s: %v\n", item.CabinTypeID, err)
```

项目已集成 Zap 结构化日志，但关键的库存操作错误仍使用 `fmt.Printf`，无法追踪和告警。

**建议**: 注入 logger 依赖，使用 `logger.Error("failed to cancel booking", zap.String("cabinTypeID", ...))`。

---

### CS-007 [MEDIUM] NATS 事件发布忽略错误

**文件**: `backend/internal/payment/service.go` L245-L246

```go
data, _ := json.Marshal(event)     // 忽略序列化错误
s.natsConn.Publish(eventType, data) // 忽略发布错误
```

关键的支付成功/退款事件发布结果被完全忽略。

**建议**: 记录错误日志并考虑重试机制或补偿队列。

---

### CS-008 [MEDIUM] 订单锁定超时与 spec 不一致

**文件**:
- `backend/internal/service/order.go` L201: `expiresAt := now.Add(30 * time.Minute)`
- `specs/001-cruise-booking-system/spec.md` L46-L48: 规定锁定 **15 分钟**

代码中设置的订单过期时间为 30 分钟，与 spec 要求的 15 分钟不一致。

**建议**: 改为 `15 * time.Minute` 或将超时时间抽取为可配置项。

---

## 3. 设计债务 🟠

### DD-001 [HIGH] 订单创建缺少数据库事务

**文件**: `backend/internal/service/order.go` L176-L329（`Create` 方法）

订单创建涉及多步操作（创建订单 → 遍历锁定库存 → 创建订单项 → 更新订单总额 → 批量创建乘客），全部操作 **未包裹在数据库事务中**。任何中间步骤失败会导致数据不一致（如库存被锁定但订单项创建失败）。

当前的回滚逻辑仅在单个项失败时尝试解锁 (L279)，但如果有多个 items 且后续 item 失败，前面已成功的 item 不会回滚。

**建议**: 使用 GORM 事务 `db.Transaction()` 包裹整个创建流程。

---

### DD-002 [HIGH] Repository 接口违反单一职责原则（SRP）

**文件**: `backend/internal/repository/order.go` L12-L58

`OrderRepository` 接口包含 **27 个方法**，同时管理 Order、OrderItem、Passenger、Payment、RefundRequest 五种实体的 CRUD。

**建议**: 拆分为独立的 `OrderRepository`、`OrderItemRepository`、`PassengerRepository`、`PaymentRepository`、`RefundRepository`。

---

### DD-003 [HIGH] 跨包添加方法违反封装原则

**文件**: `backend/internal/service/order_state.go` L307-L311

```go
func (o *domain.Order) IsExpired() bool { ... }
```

在 `service` 包中为 `domain.Order` 添加方法（通过非局部类型扩展），违反 Go 的类型系统约定。仅在同一包内才可为类型添加方法。

> **注意**: 这在 Go 中实际上应该会编译失败。如果能编译，说明代码可能未真正运行过。

**建议**: 将 `IsExpired()` 移至 `domain/order.go` 中定义。

---

### DD-004 [MEDIUM] 缺少分布式锁机制

**文件**: `backend/internal/service/inventory.go`

Spec 明确要求双层锁机制（Redis 分布式锁 + 数据库乐观锁），但当前 `InventoryService` 仅使用数据库级操作，**完全没有 Redis 分布式锁**。

**建议**: 在库存扣减前获取 Redis 分布式锁，按照 spec FR-028A/B/C 规范实现。

---

### DD-005 [MEDIUM] 错误比较使用 `errors.New()` 新实例

**文件**: `backend/internal/service/inventory.go` L173

```go
if errors.Is(err, errors.New("record not found")) {
```

`errors.Is` 比较的是指针/值相等，每次 `errors.New()` 创建新实例，两个不同实例永远不相等，此分支永远不会命中。

**建议**: 使用 `errors.Is(err, gorm.ErrRecordNotFound)` 或预定义的 sentinel error。

---

### DD-006 [MEDIUM] `notification.go` ID 类型不一致

**文件**: `backend/internal/service/notification.go`

`NotificationService` 接口使用 `uint64` 作为 ID 类型，但其他模块（Order、User 等）使用 `string`（UUID）。类型不一致增加集成难度。

**建议**: 统一使用 `string` (UUID) 或 `uint64`，保持一致。

---

### DD-007 [MEDIUM] 硬编码管理员用户 ID

**文件**: `backend/internal/service/notification.go` L414

```go
req := CreateNotificationRequest{
    UserID: 1, // Admin user ID
```

库存预警通知硬编码发送给 ID=1 的用户。

**建议**: 查询具有管理员角色的用户列表，或使用通知组/频道机制。

---

## 4. 测试债务 🟡

### TD-001 [HIGH] 后端测试覆盖率严重不足

**当前状态**: 仅存在 **4 个测试文件**：
- `internal/handler/cruise_test.go`
- `internal/payment/service_test.go`
- `internal/service/cruise_test.go`
- `internal/service/order_test.go`

**缺失测试**（按严重程度排序）:
| 模块 | 文件数 | 测试文件 |
|------|-------|---------|
| handler/auth.go | 1 | ❌ |
| handler/order.go | 1 | ❌ |
| handler/payment.go | 1 | ❌ |
| handler/admin_*.go | 3 | ❌ |
| handler/user.go | 1 | ❌ |
| service/inventory.go | 1 | ❌ |
| service/order_state.go | 1 | ❌ |
| service/facility.go | 1 | ❌ |
| service/notification.go | 1 | ❌ |
| service/refund.go | 1 | ❌ |
| repository/*.go | 10 | ❌ |
| middleware/jwt.go | 1 | ❌ |
| payment/wechat.go | 1 | ❌ |

Spec 要求 **100% 测试覆盖率**（SC-018），当前实际覆盖率估计 < 15%。

---

### TD-002 [HIGH] 前端测试覆盖率不足

**当前状态**: 仅存在 **6 个测试文件**:
- `frontend-admin/tests/components/ImageUpload.spec.ts`
- `frontend-admin/tests/e2e/admin-panel.spec.ts`
- `frontend-mini/tests/components/CruiseCard.spec.ts`
- `frontend-web/tests/components/CruiseCard.spec.ts`
- `frontend-web/tests/e2e/booking-flow.spec.ts`
- `frontend-web/tests/e2e/cruise-browsing.spec.ts`

**缺失测试**:
- 登录页面组件测试
- Pinia stores 单元测试（auth、cruise 等）
- 前台页面组件测试（订单、用户中心等）
- 小程序页面测试
- API composables 测试
- 中间件（auth.ts）测试

---

### TD-003 [MEDIUM] 缺少集成测试

项目 `tests/` 目录下的 `unit/`、`integration/`、`e2e/` 子目录均为空或未包含实际测试文件。Spec 要求所有 API 端点 100% 集成测试覆盖。

---

### TD-004 [MEDIUM] 缺少并发安全测试

Spec 要求 1000 并发用户同时订同一舱位的测试（SC-003），当前无任何并发竞争测试。

---

### TD-005 [MEDIUM] 缺少支付流程测试

微信支付的核心流程（创建支付、回调处理、退款）虽有 `service_test.go`，但缺少：
- 幂等性测试
- 重复支付测试
- 签名验证测试
- 超时场景测试

---

### TD-006 [LOW] 缺少 Race Condition 检测配置

Spec 要求使用 Go race detector 无竞态条件（SC-010A），但项目中未配置 `go test -race` 的 CI 步骤。

---

## 5. 文档债务 🟡

### DOC-001 [MEDIUM] Swagger/OpenAPI 文档可能过期

`backend/docs/swagger.json` 存在但未确认是否与最新 handler 代码同步（新增的 user、notification、analytics 等 handler 可能未包含在内）。

**建议**: 在 CI 中增加 `swag init` 并检查 diff 以确保文档同步。

---

### DOC-002 [MEDIUM] 10+ TODO 注释未清理

全项目存在至少 10 处 TODO 注释：
- `handler/auth.go`: 6 处（登录验证、token 刷新、token 黑名单、密码修改、微信登录）
- `handler/admin_order.go`: 1 处（分页计数）
- `handler/order_query.go`: 1 处（管理员检查）
- `cmd/api/main.go`: 3 处（中间件、路由、服务器启动）

**建议**: 将 TODO 转化为 issue 跟踪，并清理已临时完成的 TODO。

---

### DOC-003 [LOW] 项目宪法文件未填写

**文件**: `.specify/memory/constitution.md`

文件仍为模板状态，所有原则、治理规则均为占位符（`[PRINCIPLE_1_NAME]`、`[PRINCIPLE_1_DESCRIPTION]` 等）。

---

### DOC-004 [LOW] `notification.go` 存在结构体标签语法问题

**文件**: `backend/internal/service/notification.go` L79-L90

```go
type CreateNotificationRequest struct {
    UserID  uint64  `json:"user_id" validate:"required"    // ← 缺少闭合反引号
    Type    string  `json:"type" validate:"required,oneof=..."
    ...
}
```

多个字段的 struct tag 缺少闭合反引号，Go 编译器应该会报错。这暗示代码可能从未被编译运行。

---

## 修复优先级建议

### P0 — 立即修复（阻塞生产部署）

| ID | 问题 | 预估工时 |
|-----|------|---------|
| SEC-001 | 实现真实认证系统 | 2-3天 |
| SEC-003 | 实现 AES-GCM 解密 | 0.5天 |
| CS-001 | 修复 DSN 生成 Bug | 10分钟 |
| CS-004 | 实现 `getCurrentTimestamp()` | 10分钟 |
| CS-005 | 实现 `IsExpired()` | 30分钟 |
| DD-003 | 移动 `IsExpired()` 至 domain 包 | 15分钟 |
| DOC-004 | 修复 struct tag 语法 | 15分钟 |

### P1 — 尽快修复（严重影响业务正确性）

| ID | 问题 | 预估工时 |
|-----|------|---------|
| SEC-004 | 实现支付幂等性 | 1-2天 |
| CS-002 | 金额改用整数或 Decimal | 2-3天 |
| DD-001 | 订单创建添加事务 | 1天 |
| DD-004 | 实现 Redis 分布式锁 | 1-2天 |
| CS-008 | 统一锁定超时为 15 分钟 | 10分钟 |
| DD-005 | 修复 `errors.Is` 比较 | 10分钟 |

### P2 — 计划修复（改善质量与可维护性）

| ID | 问题 | 预估工时 |
|-----|------|---------|
| SEC-002 | Token 黑名单机制 | 1天 |
| SEC-005 | 清理 `.env.example` | 15分钟 |
| SEC-006 | 修复变量遮蔽 | 15分钟 |
| SEC-007 | JWT 签名算法验证 | 15分钟 |
| CS-003 | 抽取价格计算方法 | 30分钟 |
| CS-006 | 替换 fmt.Printf 为 Zap | 30分钟 |
| CS-007 | 处理 NATS 发布错误 | 30分钟 |
| DD-002 | 拆分 OrderRepository | 2天 |
| DD-006 | 统一 ID 类型 | 1天 |
| DD-007 | 移除硬编码管理员 ID | 30分钟 |
| TD-001~006 | 补充测试 | 10-15天 |
| DOC-001~003 | 文档更新 | 1天 |
