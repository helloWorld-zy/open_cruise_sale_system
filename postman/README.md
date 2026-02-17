# Postman 验收资产说明

目录内容：
- `CruiseBooking-Prelaunch.postman_collection.json`
- `CruiseBooking-Local.postman_environment.json`

## 导入步骤

1. 打开 Postman，导入 Collection 与 Environment 两个 JSON。
2. 选择环境 `CruiseBooking Local`。
3. 填写必要变量：
   - `token`
   - `voyageId` / `cruiseId` / `cabinId` / `cabinTypeId`
   - 微信回调相关：`wechatpaySignature` `wechatpayTimestamp` `wechatpayNonce` `wechatpaySerial` `wechatCallbackBody`

## 推荐执行顺序

1. `0. Health / GET Health`
2. `1. Order / POST Create Order`（自动回填 `orderId`）
3. `2. Payment / POST Create Payment`（自动回填 `paymentId`、`paymentNo`）
4. `3. WeChat Callback / POST Callback First Time`
5. `3. WeChat Callback / POST Callback Replay (Same Packet)`
6. `1. Order / GET Order Detail` 或 `2. Payment / GET Payment By Order` 校验状态

## 注意

- 回调验签链路依赖真实且匹配的签名上下文，若使用伪造值，预期返回 `fail`。
- 重放请求用于验证幂等与防重放，不应产生重复业务副作用。

## Newman 本地批跑

- 脚本：`postman/run-newman-local.ps1`
- 先安装：`npm i -g newman`

示例（完整执行）：

```powershell
pwsh ./postman/run-newman-local.ps1 `
   -BaseUrl "http://localhost:8080/api/v1" `
   -Token "<JWT_TOKEN>" `
   -VoyageId "<VOYAGE_ID>" `
   -CruiseId "<CRUISE_ID>" `
   -CabinId "<CABIN_ID>" `
   -CabinTypeId "<CABIN_TYPE_ID>" `
   -WechatpaySignature "<SIGNATURE>" `
   -WechatpayTimestamp "<TIMESTAMP>" `
   -WechatpayNonce "<NONCE>" `
   -WechatpaySerial "<SERIAL>" `
   -WechatCallbackBody '{"id":"...","resource":{...}}'
```

示例（跳过回调用例）：

```powershell
pwsh ./postman/run-newman-local.ps1 -Token "<JWT_TOKEN>" -SkipCallback
```

报告输出：`postman/newman-report.xml`

## GitHub Actions 批跑

- Workflow：`.github/workflows/postman-acceptance.yml`
- 触发方式：`workflow_dispatch`
- 输入参数：
   - `base_url`（必填，形如 `http://host:port/api/v1`）
   - `skip_callback`（可选）

需要在仓库 Secrets 中配置：

- `POSTMAN_TOKEN`
- `POSTMAN_VOYAGE_ID`
- `POSTMAN_CRUISE_ID`
- `POSTMAN_CABIN_ID`
- `POSTMAN_CABIN_TYPE_ID`
- `POSTMAN_WECHATPAY_SIGNATURE`
- `POSTMAN_WECHATPAY_TIMESTAMP`
- `POSTMAN_WECHATPAY_NONCE`
- `POSTMAN_WECHATPAY_SERIAL`
- `POSTMAN_WECHAT_CALLBACK_BODY`

## 主 CI 自动触发（已接入）

- 文件：`.github/workflows/ci.yml`
- 触发条件：
   - push 到 `release/**` 分支
   - push `v*` 标签
- 依赖关系：`postman-acceptance-auto` 依赖 `backend` job 先通过（`needs: [backend]`）
- 执行策略：默认仅跑稳定链路（`Health + Order + Payment + Order Cancel`），跳过回调文件夹
- 执行前检查：对 `${CI_POSTMAN_BASE_URL}/health` 进行最多 12 次探测（每次间隔 5 秒），健康检查通过后才执行 newman

自动批跑额外需要：

- `CI_POSTMAN_BASE_URL`（例如 `https://api.example.com/api/v1`）

说明：若 `CI_POSTMAN_BASE_URL` 或 `POSTMAN_TOKEN` 为空，自动 Job 会在 `Guard required secrets` 步骤输出 skip 并结束后续执行。
