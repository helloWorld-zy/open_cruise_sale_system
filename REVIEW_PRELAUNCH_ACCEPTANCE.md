# 上线前验收脚本清单（重点链路）

更新时间：2026-02-14
适用范围：`backend` 服务（支付、订单、库存）

---

## 0. 预检

- 确认依赖服务可用：PostgreSQL、Redis、NATS
- 进入目录：`cd backend`
- 回归基线：`go test ./...`
- 启动服务（按你当前运行方式）：`go run ./cmd/api`

通过标准：
- 测试全绿
- 服务可正常启动，健康检查 `GET /api/v1/health` 返回 200

---

## 1. 支付回调重放与签名安全验收

目标：验证“签名上下文 + 时间窗 + nonce 防重放 + 幂等键”同时生效。

### 1.1 正常回调（首包）

步骤：
1. 准备一个待支付订单并发起支付，拿到订单/支付号。
2. 构造一份合法微信回调报文（含 `Wechatpay-Signature/Timestamp/Nonce/Serial`）。
3. POST 到：`/api/v1/payments/callback/wechat`。

预期：
- HTTP 200，响应 `success`
- 支付状态更新为成功
- 订单支付状态更新为已支付

### 1.2 重放同一回调（同 nonce 或同业务幂等）

步骤：
1. 在首包成功后，原样再次发送同一报文。

预期：
- 接口不应重复变更业务状态（无二次记账/无重复事件副作用）
- 订单与支付最终状态保持正确且单次生效

### 1.3 时间窗与签名非法场景

步骤：
- 场景 A：将 `Wechatpay-Timestamp` 改为超出 ±5 分钟
- 场景 B：篡改 body 内容但保留旧签名
- 场景 C：使用错误 `Wechatpay-Serial`

预期：
- HTTP 400，响应 `fail`
- 数据库业务状态不发生变化

---

## 2. 库存并发锁冲突验收

目标：验证“分布式锁 + 库存下界保护 + 事务内库存锁定”在并发下无超卖、无负数。

### 2.1 并发抢同一房型

步骤：
1. 准备同一航次同一房型，库存设为 `N`（建议 1 或 2，便于观察）。
2. 使用并发压测工具对创建订单接口并发请求（如 10~50 并发）。

建议命令（示例，需按真实 token/body 调整）：
- `hey -n 50 -c 20 -m POST -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" -d "<ORDER_JSON>" http://localhost:8080/api/v1/orders`

预期：
- 成功单数 <= 可售库存
- 失败请求返回“库存不足/系统繁忙重试”类错误，不出现 500 泛化错误
- `available_cabins/locked_cabins/booked_cabins` 均不出现负数

### 2.2 取消/确认交错并发

步骤：
1. 对同一订单或同一库存资源，交错触发取消与确认动作。
2. 重复执行多轮。

预期：
- 不出现库存负数
- 订单状态机保持合法流转（无非法回跳）

---

## 3. 超时取消联调验收

目标：验证“15 分钟规则”在下单、通知、超时任务三处一致。

### 3.1 新建待支付订单

步骤：
1. 下单后读取订单 `expires_at`。
2. 校验通知文案是否为“15 分钟内支付”。

预期：
- `expires_at` 与创建时间相差约 15 分钟
- 通知文案与规则一致

### 3.2 超时任务触发

步骤：
1. 构造一个已过期且待支付订单（测试环境可手工改 `expires_at` 为过去时间）。
2. 触发超时任务（按你们当前 job 启动方式）。

预期：
- 订单自动变更为已取消
- 对应锁定库存被释放
- 重复触发任务时不应二次扣改或产生异常

---

## 4. 数据核验 SQL（验收后执行）

> 下列 SQL 作为验收核对模板，按实际库名执行。

```sql
-- A. 检查库存是否出现负数
SELECT voyage_id, cabin_type_id, available_cabins, locked_cabins, booked_cabins
FROM cabin_inventory
WHERE available_cabins < 0 OR locked_cabins < 0 OR booked_cabins < 0;

-- B. 检查待支付但已过期订单
SELECT id, order_number, status, payment_status, expires_at
FROM orders
WHERE status = 'pending' AND expires_at < NOW();

-- C. 检查支付成功的订单状态一致性
SELECT o.id AS order_id, o.status AS order_status, o.payment_status, p.status AS payment_status_detail
FROM orders o
JOIN payments p ON p.order_id = o.id
WHERE p.status = 'success' AND o.payment_status <> 'paid';
```

通过标准：
- A 返回 0 行
- B 在超时任务执行后仅保留“刚过期尚未扫描”窗口内数据
- C 返回 0 行

---

## 5. 上线阻断条件（任一命中则禁止发布）

- 支付回调可被重放并造成重复状态变更
- 并发下出现超卖或库存负数
- 15 分钟超时规则在订单/通知/job 任一处不一致
- `go test ./...` 非全绿

---

## 6. 验收记录模板

- 验收时间：
- 验收环境：
- 执行人：
- 用例通过率：
- 阻断项：
- 结论（通过/不通过）：
