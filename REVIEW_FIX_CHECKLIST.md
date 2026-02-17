# 修复清单（独立审查项）

更新时间：2026-02-14

## P0

- [x] 支付回调安全校验增强（签名元数据、时间戳窗口、nonce 防重放、serial 校验）
  - `backend/internal/handler/payment.go`
  - `backend/internal/payment/wechat.go`
- [x] 支付金额一致性强校验（回调金额与订单金额一致）
  - `backend/internal/payment/service.go`
- [x] 支付回调幂等键落地（Redis `payment:idempotent:{order_id}`，TTL 24h）
  - `backend/internal/payment/service.go`
- [x] 测试与编译阻塞项修复（`go test ./...` 通过）
  - `backend/internal/handler/cruise_test.go`
  - `backend/internal/payment/service_test.go`
  - `backend/internal/service/cruise_test.go`
  - `backend/internal/service/order_test.go`

## P1

- [x] 订单创建事务边界修复（库存锁定改为事务内仓储实例）
  - `backend/internal/repository/order.go`
  - `backend/internal/service/order.go`
- [x] 库存边界保护（避免 `locked/booked` 负数）
  - `backend/internal/repository/inventory.go`
- [x] 分布式锁补齐（下单链路 + 库存服务，Redis SETNX + Lua 解锁）
  - `backend/internal/service/order.go`
  - `backend/internal/service/inventory.go`
- [x] 超时规则统一为 15 分钟（任务回退逻辑、通知文案）
  - `backend/internal/jobs/order_timeout.go`
  - `backend/internal/service/notification.go`
- [x] 主路由补齐（auth/orders/payments/user）
  - `backend/cmd/api/routes.go`

## P2

- [x] 启动端口改为配置驱动（移除硬编码 `:8080`）
  - `backend/cmd/api/main.go`
- [x] 关键事件日志改为统一日志输出（移除 `fmt.Printf`）
  - `backend/internal/payment/service.go`
- [x] 订单详情权限 TODO 修复（补充管理员角色判定）
  - `backend/internal/handler/order_query.go`
- [x] 订单过期判断健壮性增强（RFC3339 解析）
  - `backend/internal/domain/order.go`

## 验证结果

- [x] 后端全量测试通过：`go test ./...`
