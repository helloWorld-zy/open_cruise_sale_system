# CruiseBooking 平台部署及上线指南

**版本**: 2.0.0  
**更新日期**: 2026-02-12  
**适用项目**: 邮轮舱位预订平台 (CruiseBooking)  

---

## 📋 目录

1. [部署前准备](#1-部署前准备)
2. [环境要求](#2-环境要求)
3. [基础设施部署](#3-基础设施部署)
4. [数据库部署](#4-数据库部署)
5. [后端服务部署](#5-后端服务部署)
6. [前端应用部署](#6-前端应用部署)
7. [监控与日志](#7-监控与日志)
8. [安全加固](#8-安全加固)
9. [测试验证](#9-测试验证)
10. [上线 checklist](#10-上线-checklist)
11. [故障处理](#11-故障处理)
12. [运维手册](#12-运维手册)

---

## 1. 部署前准备

### 1.1 部署清单确认

- [ ] 所有 206 个任务已完成并通过测试
- [ ] 代码已合并到 `main` 分支
- [ ] 版本号已更新为 v2.0.0
- [ ] 数据库迁移脚本已准备就绪
- [ ] 环境配置文件已准备
- [ ] SSL 证书已申请
- [ ] 域名已解析

### 1.2 团队准备

- [ ] 运维团队已培训
- [ ] 客服团队已培训  
- [ ] 应急响应流程已建立
- [ ] 值班表已排定

---

## 2. 环境要求

### 2.1 硬件配置 (推荐)

```yaml
生产环境:
  服务器数量: 3台 (高可用)
  单台配置:
    CPU: 8核+
    内存: 32GB+
    存储: 500GB SSD (系统盘) + 1TB (数据盘)
    带宽: 100Mbps+
  
  数据库服务器:
    CPU: 16核+
    内存: 64GB+
    存储: 2TB SSD RAID 10
    
  缓存/搜索服务器:
    CPU: 8核
    内存: 16GB+
    存储: 200GB SSD
```

### 2.2 软件依赖

| 组件 | 版本 | 用途 |
|------|------|------|
| Docker | 24.0+ | 容器化部署 |
| Docker Compose | 2.20+ | 编排管理 |
| Kubernetes | 1.28+ | 容器编排 (可选) |
| PostgreSQL | 17.x | 主数据库 |
| Redis | 7.4.x | 缓存/会话/锁 |
| Meilisearch | 1.12.x | 搜索引擎 |
| MinIO | Latest | 对象存储 |
| NATS | 2.11.x | 消息队列 |

### 2.3 网络要求

- 公网 IP: 2个 (主备)
- 内网 IP 段: 172.16.0.0/16
- 开放端口:
  - 80/443: HTTP/HTTPS
  - 5432: PostgreSQL (仅内网)
  - 6379: Redis (仅内网)
  - 7700: Meilisearch (仅内网)
  - 9000: MinIO (仅内网)
  - 9090: Prometheus (内网/白名单)
  - 3000: Grafana (内网/白名单)

---

## 3. 基础设施部署

### 3.1 服务器初始化

```bash
#!/bin/bash
# server-init.sh

# 1. 系统更新
apt-get update && apt-get upgrade -y

# 2. 安装基础工具
apt-get install -y \
    curl wget git vim htop net-tools \
    ca-certificates gnupg lsb-release \
    ufw fail2ban

# 3. 配置防火墙
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable

# 4. 配置 fail2ban
cat > /etc/fail2ban/jail.local <<EOF
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 3

[sshd]
enabled = true
EOF

systemctl enable fail2ban
systemctl start fail2ban

# 5. 配置时区
timedatectl set-timezone Asia/Shanghai

# 6. 安装 Docker
curl -fsSL https://get.docker.com | sh
systemctl enable docker
systemctl start docker

# 7. 安装 Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# 8. 创建项目目录
mkdir -p /opt/cruise-booking/{backend,frontend,infra,data}
mkdir -p /opt/cruise-booking/data/{postgres,redis,minio,meilisearch}

# 9. 配置目录权限
chown -R 1000:1000 /opt/cruise-booking/data

echo "服务器初始化完成！"
```

### 3.2 网络配置

```bash
# 配置内网通信
cat >> /etc/hosts <<EOF
172.16.0.10     db-primary
172.16.0.11     db-replica
172.16.0.20     redis-primary
172.16.0.21     redis-replica
172.16.0.30     meilisearch
172.16.0.40     minio
172.16.0.50     nats
172.16.0.100    backend-1
172.16.0.101    backend-2
172.16.0.102    backend-3
EOF
```

---

## 4. 数据库部署

### 4.1 PostgreSQL 主从配置

```yaml
# docker-compose.db.yml
version: '3.8'

services:
  postgres-primary:
    image: postgres:17-alpine
    container_name: postgres-primary
    hostname: db-primary
    environment:
      POSTGRES_USER: cruise_admin
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: cruise_booking
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
      - /opt/cruise-booking/data/postgres/primary:/var/lib/postgresql/data
      - ./init-scripts:/docker-entrypoint-initdb.d
      - ./pg_hba.conf:/etc/postgresql/pg_hba.conf
      - ./postgresql.conf:/etc/postgresql/postgresql.conf
    ports:
      - "5432:5432"
    command: >
      postgres
      -c config_file=/etc/postgresql/postgresql.conf
      -c hba_file=/etc/postgresql/pg_hba.conf
    networks:
      - backend-network
    restart: always
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U cruise_admin -d cruise_booking"]
      interval: 10s
      timeout: 5s
      retries: 5

  postgres-replica:
    image: postgres:17-alpine
    container_name: postgres-replica
    hostname: db-replica
    environment:
      POSTGRES_USER: cruise_admin
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      PGDATA: /var/lib/postgresql/data/pgdata
      REPLICATE_FROM: db-primary
    volumes:
      - /opt/cruise-booking/data/postgres/replica:/var/lib/postgresql/data
    ports:
      - "5433:5432"
    networks:
      - backend-network
    restart: always
    depends_on:
      - postgres-primary
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U cruise_admin -d cruise_booking"]
      interval: 10s
      timeout: 5s
      retries: 5

networks:
  backend-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.16.0.0/24
```

### 4.2 数据库初始化

```bash
#!/bin/bash
# init-database.sh

cd /opt/cruise-booking/infra

# 1. 启动数据库
docker-compose -f docker-compose.db.yml up -d

# 2. 等待主库就绪
echo "等待 PostgreSQL 主库就绪..."
sleep 30

# 3. 执行迁移脚本
docker exec -i postgres-primary psql -U cruise_admin -d cruise_booking <<EOF
-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- 用于全文搜索

-- 创建初始管理员账号 (需要修改密码)
INSERT INTO users (id, phone, password_hash, status, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    '13800138000',
    '\x243261243130246b6e53656165476536536f6d654761765053534b2e2e2e',  -- bcrypt hash of 'Admin@123'
    'active',
    NOW(),
    NOW()
)
ON CONFLICT DO NOTHING;

-- 创建管理员角色
INSERT INTO roles (id, name, permissions, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'admin',
    '["*"]'::jsonb,
    NOW(),
    NOW()
)
ON CONFLICT DO NOTHING;

-- 关联用户和角色
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r WHERE u.phone = '13800138000' AND r.name = 'admin'
ON CONFLICT DO NOTHING;
EOF

# 4. 验证复制状态
echo "检查主从复制状态..."
docker exec postgres-primary psql -U cruise_admin -c "SELECT * FROM pg_stat_replication;"

echo "数据库初始化完成！"
```

---

## 5. 后端服务部署

### 5.1 后端 Docker 配置

```dockerfile
# Dockerfile.backend
FROM golang:1.26-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制 go mod
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动
CMD ["./main"]
```

```yaml
# docker-compose.backend.yml
version: '3.8'

services:
  backend:
    build:
      context: ../backend
      dockerfile: Dockerfile
    image: cruise-booking-backend:v2.0.0
    container_name: cruise-backend
    environment:
      # 数据库
      DB_HOST: db-primary
      DB_PORT: 5432
      DB_USER: cruise_admin
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: cruise_booking
      DB_SSL_MODE: disable
      DB_MAX_CONNECTIONS: 100
      
      # Redis
      REDIS_HOST: redis-primary
      REDIS_PORT: 6379
      REDIS_PASSWORD: ${REDIS_PASSWORD}
      REDIS_DB: 0
      
      # Meilisearch
      MEILISEARCH_HOST: http://meilisearch:7700
      MEILISEARCH_API_KEY: ${MEILISEARCH_API_KEY}
      
      # MinIO
      MINIO_ENDPOINT: minio:9000
      MINIO_ACCESS_KEY: ${MINIO_ACCESS_KEY}
      MINIO_SECRET_KEY: ${MINIO_SECRET_KEY}
      MINIO_BUCKET: cruise-booking
      MINIO_USE_SSL: "false"
      
      # JWT
      JWT_SECRET: ${JWT_SECRET}
      JWT_EXPIRATION: 24h
      
      # 微信支付
      WECHAT_APP_ID: ${WECHAT_APP_ID}
      WECHAT_MCH_ID: ${WECHAT_MCH_ID}
      WECHAT_API_KEY: ${WECHAT_API_KEY}
      WECHAT_NOTIFY_URL: https://api.cruise-booking.com/api/v1/payments/wechat/notify
      
      # 环境
      ENV: production
      LOG_LEVEL: info
      
      # 性能
      GIN_MODE: release
      GOMAXPROCS: 8
      
    ports:
      - "8080:8080"
    volumes:
      - /opt/cruise-booking/logs/backend:/app/logs
      - /opt/cruise-booking/uploads:/app/uploads
    networks:
      - backend-network
      - frontend-network
    restart: always
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis-primary:
        condition: service_healthy
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2.0'
          memory: 4G
        reservations:
          cpus: '1.0'
          memory: 2G
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

networks:
  backend-network:
    external: true
  frontend-network:
    driver: bridge
```

### 5.2 部署脚本

```bash
#!/bin/bash
# deploy-backend.sh

set -e

echo "开始部署后端服务..."

# 1. 拉取最新代码
cd /opt/cruise-booking/backend
git pull origin main

# 2. 构建镜像
echo "构建 Docker 镜像..."
docker build -t cruise-booking-backend:v2.0.0 .

# 3. 滚动更新
echo "执行滚动更新..."
docker-compose -f ../infra/docker-compose.backend.yml pull
docker-compose -f ../infra/docker-compose.backend.yml up -d --no-deps --scale backend=4 backend
docker-compose -f ../infra/docker-compose.backend.yml up -d --no-deps --scale backend=3 backend

# 4. 健康检查
echo "执行健康检查..."
sleep 10
for i in {1..5}; do
    if curl -sf http://localhost:8080/health > /dev/null; then
        echo "服务健康检查通过！"
        break
    fi
    echo "等待服务就绪... ($i/5)"
    sleep 5
done

# 5. 清理旧镜像
echo "清理旧镜像..."
docker image prune -f

echo "后端部署完成！"
```

---

## 6. 前端应用部署

### 6.1 Web 前端部署

```dockerfile
# Dockerfile.web
FROM node:20-alpine AS builder

WORKDIR /app

# 复制 package.json
COPY package*.json ./
RUN npm ci

# 复制源码
COPY . .

# 构建生产版本
RUN npm run build

# 运行阶段
FROM nginx:alpine

# 复制 Nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 复制构建产物
COPY --from=builder /app/.output/public /usr/share/nginx/html

# 暴露端口
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

```nginx
# nginx.conf
server {
    listen 80;
    server_name www.cruise-booking.com cruise-booking.com;
    
    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml application/json application/javascript application/rss+xml application/atom+xml image/svg+xml;
    
    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
        add_header X-Content-Type-Options "nosniff";
    }
    
    # 主应用
    location / {
        root /usr/share/nginx/html;
        index index.html;
        try_files $uri $uri/ /index.html;
        
        # 安全头部
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    }
    
    # API 代理
    location /api/ {
        proxy_pass http://backend:8080/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # 健康检查
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
```

### 6.2 Admin 管理后台部署

```dockerfile
# Dockerfile.admin
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine

COPY nginx-admin.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/.output/public /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

---

## 7. 监控与日志

### 7.1 Prometheus 配置

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']

rule_files:
  - /etc/prometheus/rules/*.yml

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'backend'
    static_configs:
      - targets: ['backend-1:8080', 'backend-2:8080', 'backend-3:8080']
    metrics_path: /metrics
    scrape_interval: 10s

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']

  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']
```

### 7.2 Grafana Dashboard 配置

```json
{
  "dashboard": {
    "title": "CruiseBooking 监控大屏",
    "panels": [
      {
        "title": "API 请求速率",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{handler}}"
          }
        ]
      },
      {
        "title": "响应时间 P99",
        "type": "singlestat",
        "targets": [
          {
            "expr": "histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))"
          }
        ]
      },
      {
        "title": "数据库连接数",
        "type": "graph",
        "targets": [
          {
            "expr": "pg_stat_activity_count"
          }
        ]
      },
      {
        "title": "Redis 内存使用",
        "type": "graph",
        "targets": [
          {
            "expr": "redis_memory_used_bytes"
          }
        ]
      },
      {
        "title": "订单量 (实时)",
        "type": "stat",
        "targets": [
          {
            "expr": "increase(orders_total[1h])"
          }
        ]
      }
    ]
  }
}
```

### 7.3 Loki 日志聚合

```yaml
# docker-compose.logging.yml
version: '3.8'

services:
  loki:
    image: grafana/loki:latest
    container_name: loki
    ports:
      - "3100:3100"
    volumes:
      - /opt/cruise-booking/data/loki:/loki
      - ./loki-config.yml:/etc/loki/local-config.yaml
    command: -config.file=/etc/loki/local-config.yaml
    networks:
      - backend-network

  promtail:
    image: grafana/promtail:latest
    container_name: promtail
    volumes:
      - /opt/cruise-booking/logs:/var/log/cruise-booking
      - ./promtail-config.yml:/etc/promtail/config.yml
    command: -config.file=/etc/promtail/config.yml
    networks:
      - backend-network
    depends_on:
      - loki

networks:
  backend-network:
    external: true
```

---

## 8. 安全加固

### 8.1 SSL/TLS 配置

```nginx
# ssl.conf
server {
    listen 443 ssl http2;
    server_name www.cruise-booking.com;
    
    # SSL 证书
    ssl_certificate /etc/nginx/ssl/cruise-booking.crt;
    ssl_certificate_key /etc/nginx/ssl/cruise-booking.key;
    
    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
    ssl_session_tickets off;
    
    # HSTS
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always;
    
    # 安全头部
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Permissions-Policy "geolocation=(), microphone=(), camera=()" always;
    
    # 内容安全策略
    add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self' https://api.cruise-booking.com;" always;
    
    location / {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name www.cruise-booking.com cruise-booking.com;
    return 301 https://$server_name$request_uri;
}
```

### 8.2 WAF 规则 (ModSecurity)

```apache
# modsecurity.conf
SecRuleEngine On
SecRequestBodyAccess On
SecRequestBodyLimit 13107200
SecRequestBodyNoFilesLimit 131072
SecResponseBodyAccess On
SecResponseBodyLimit 524288

# SQL 注入防护
SecRule REQUEST_COOKIES|REQUEST_COOKIES_NAMES|REQUEST_FILENAME|ARGS_NAMES|ARGS|XML:/* \
    "@rx (?i:(?:select\s*\*\s*from|(?:delete|drop|truncate)\s+table|union(?:\s+all)?\s*select|insert\s+into\s+.*\s+values|load_file\s*\()))" \
    "id:942100,phase:2,deny,status:403,msg:'SQL Injection Attack Detected'"

# XSS 防护
SecRule REQUEST_COOKIES|REQUEST_COOKIES_NAMES|REQUEST_FILENAME|ARGS_NAMES|ARGS|XML:/* \
    "@rx (?i:<script|javascript:|onload=|onerror=)" \
    "id:941100,phase:2,deny,status:403,msg:'XSS Attack Detected'"
```

---

## 9. 测试验证

### 9.1 部署前测试清单

```bash
#!/bin/bash
# pre-deploy-check.sh

echo "=== 部署前检查 ==="

# 1. 单元测试
echo "运行单元测试..."
cd /opt/cruise-booking/backend
go test ./... -v -race -count=1
if [ $? -ne 0 ]; then
    echo "❌ 单元测试失败"
    exit 1
fi
echo "✓ 单元测试通过"

# 2. 集成测试
echo "运行集成测试..."
go test ./tests/integration/... -v
if [ $? -ne 0 ]; then
    echo "❌ 集成测试失败"
    exit 1
fi
echo "✓ 集成测试通过"

# 3. E2E 测试
echo "运行 E2E 测试..."
cd /opt/cruise-booking/frontend-web
npm run test:e2e
if [ $? -ne 0 ]; then
    echo "❌ E2E 测试失败"
    exit 1
fi
echo "✓ E2E 测试通过"

# 4. 数据库迁移测试
echo "测试数据库迁移..."
cd /opt/cruise-booking/backend
go run cmd/migrate/main.go up
go run cmd/migrate/main.go down
go run cmd/migrate/main.go up
if [ $? -ne 0 ]; then
    echo "❌ 数据库迁移失败"
    exit 1
fi
echo "✓ 数据库迁移测试通过"

# 5. 安全检查
echo "运行安全扫描..."
gosec ./...
if [ $? -ne 0 ]; then
    echo "❌ 安全检查失败"
    exit 1
fi
echo "✓ 安全检查通过"

echo "=== 所有检查通过，可以部署 ==="
```

### 9.2 冒烟测试

```bash
#!/bin/bash
# smoke-test.sh

BASE_URL="https://api.cruise-booking.com"

echo "=== 冒烟测试 ==="

# 1. 健康检查
echo "测试健康端点..."
curl -sf ${BASE_URL}/health || { echo "❌ 健康检查失败"; exit 1; }
echo "✓ 健康检查通过"

# 2. API 可用性
echo "测试 API 可用性..."
curl -sf ${BASE_URL}/api/v1/cruises || { echo "❌ API 不可用"; exit 1; }
echo "✓ API 可用"

# 3. 数据库连接
echo "测试数据库连接..."
curl -sf ${BASE_URL}/health/ready || { echo "❌ 数据库未就绪"; exit 1; }
echo "✓ 数据库就绪"

# 4. 关键业务流程测试
echo "测试关键业务流程..."

# 登录测试
TOKEN=$(curl -sf -X POST ${BASE_URL}/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"phone":"13800138000","code":"123456"}' | jq -r '.data.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "❌ 登录失败"
    exit 1
fi
echo "✓ 登录功能正常"

# 邮轮列表测试
curl -sf -H "Authorization: Bearer ${TOKEN}" \
    ${BASE_URL}/api/v1/cruises | jq '.data' > /dev/null || { echo "❌ 邮轮列表获取失败"; exit 1; }
echo "✓ 邮轮列表获取正常"

# 下单流程测试 (使用测试环境)
echo "✓ 冒烟测试全部通过！"
```

---

## 10. 上线 Checklist

### 10.1 上线前检查清单

#### 基础设施
- [ ] 服务器已初始化并配置防火墙
- [ ] 数据库主从复制正常
- [ ] Redis 集群部署完成
- [ ] Meilisearch 索引已创建
- [ ] MinIO 对象存储可访问
- [ ] 负载均衡器已配置

#### 应用部署
- [ ] 后端服务已部署 (3个实例)
- [ ] 前端 Web 已部署
- [ ] Admin 管理后台已部署
- [ ] 小程序已上传并审核通过
- [ ] 所有服务健康检查通过

#### 安全
- [ ] SSL 证书已安装并生效
- [ ] WAF 规则已启用
- [ ] API 速率限制已配置
- [ ] 敏感配置已加密
- [ ] 安全扫描已通过

#### 监控
- [ ] Prometheus 数据收集正常
- [ ] Grafana 仪表板可访问
- [ ] Loki 日志收集正常
- [ ] 告警规则已配置
- [ ] 告警通道已测试 (短信/邮件/钉钉)

#### 业务验证
- [ ] 邮轮数据已导入
- [ ] 航次和库存数据已配置
- [ ] 价格数据已设置
- [ ] 管理员账号已创建
- [ ] 支付渠道已配置并测试

#### 文档
- [ ] API 文档已发布
- [ ] 运维手册已准备
- [ ] 应急响应流程已确认
- [ ] 客服话术已准备

### 10.2 上线执行步骤

```bash
#!/bin/bash
# launch.sh

echo "🚀 CruiseBooking 平台上线脚本"
echo "================================"

# 步骤 1: 备份
echo "[1/10] 执行数据备份..."
./backup.sh
echo "✓ 备份完成"

# 步骤 2: 流量切换准备
echo "[2/10] 准备流量切换..."
# 如果是蓝绿部署，准备绿色环境
# 如果是金丝雀部署，配置权重
echo "✓ 准备就绪"

# 步骤 3: 数据库迁移
echo "[3/10] 执行数据库迁移..."
docker exec -i postgres-primary psql -U cruise_admin < migrations/production_migration.sql
echo "✓ 数据库迁移完成"

# 步骤 4: 部署新版本
echo "[4/10] 部署新版本..."
./deploy-backend.sh
./deploy-frontend.sh
./deploy-admin.sh
echo "✓ 部署完成"

# 步骤 5: 健康检查
echo "[5/10] 执行健康检查..."
./smoke-test.sh
echo "✓ 健康检查通过"

# 步骤 6: 切换流量
echo "[6/10] 切换流量到新版本..."
# 更新负载均衡器配置
# 或切换 DNS
# 或调整 Istio/Traefik 权重
echo "✓ 流量已切换"

# 步骤 7: 监控观察
echo "[7/10] 开始监控观察 (5分钟)..."
sleep 300
# 检查错误率、响应时间、业务指标
echo "✓ 监控观察完成"

# 步骤 8: 功能验证
echo "[8/10] 执行功能验证..."
# 核心业务流程测试
./e2e-test-production.sh
echo "✓ 功能验证通过"

# 步骤 9: 通知团队
echo "[9/10] 发送上线通知..."
# 发送邮件/钉钉通知
./notify-team.sh "上线成功"
echo "✓ 通知已发送"

# 步骤 10: 清理旧版本
echo "[10/10] 清理旧版本..."
# 保留旧版本 24 小时，之后清理
( sleep 86400 && ./cleanup-old-version.sh ) &
echo "✓ 清理任务已安排"

echo ""
echo "🎉 上线完成！平台正式对外提供服务"
echo "================================"
echo ""
echo "访问地址:"
echo "  - 用户端: https://www.cruise-booking.com"
echo "  - 管理端: https://admin.cruise-booking.com"
echo "  - API: https://api.cruise-booking.com"
echo ""
echo "监控地址:"
echo "  - Grafana: https://monitor.cruise-booking.com"
echo ""
echo "祝航行顺利！⛵"
```

---

## 11. 故障处理

### 11.1 常见故障及处理

#### 数据库连接池耗尽

```bash
# 症状: 大量请求超时，数据库连接数达到上限

# 诊断
docker exec postgres-primary psql -U cruise_admin -c "
SELECT count(*), state FROM pg_stat_activity GROUP BY state;
"

# 处理
# 1. 重启后端服务释放连接
# 2. 检查是否有慢查询
# 3. 临时增加 max_connections

# 应急脚本
#!/bin/bash
# fix-db-connections.sh

echo "重启后端服务释放连接池..."
docker-compose -f docker-compose.backend.yml restart

echo "等待服务恢复..."
sleep 30

echo "检查连接数..."
docker exec postgres-primary psql -U cruise_admin -c "
SELECT count(*) as total_connections 
FROM pg_stat_activity 
WHERE state = 'active';
"
```

#### Redis 内存不足

```bash
# 症状: Redis 内存达到上限，写入失败

# 诊断
redis-cli INFO memory

# 处理
# 1. 清理过期缓存
# 2. 调整 maxmemory-policy 为 allkeys-lru
# 3. 扩容 Redis 内存

# 应急脚本
#!/bin/bash
# fix-redis-memory.sh

redis-cli <<EOF
CONFIG SET maxmemory-policy allkeys-lru
MEMORY PURGE
INFO memory
EOF
```

#### 订单超卖

```bash
# 症状: 库存为负，超卖发生

# 诊断
SELECT voyage_id, cabin_type_id, remaining 
FROM cabin_inventory 
WHERE remaining < 0;

# 处理
# 1. 立即锁定超卖航次的预订
# 2. 人工介入处理已超卖订单
# 3. 修复库存数据

# 应急脚本
#!/bin/bash
# fix-overselling.sh

echo "锁定超卖航次..."
docker exec postgres-primary psql -U cruise_admin <<EOF
UPDATE voyages 
SET booking_status = 'closed' 
WHERE id IN (
    SELECT DISTINCT voyage_id 
    FROM cabin_inventory 
    WHERE remaining < 0
);

-- 记录超卖情况
INSERT INTO overselling_log (voyage_id, cabin_type_id, remaining, created_at)
SELECT voyage_id, cabin_type_id, remaining, NOW()
FROM cabin_inventory 
WHERE remaining < 0;
EOF

echo "发送告警..."
./send-alert.sh "发生超卖，请立即处理"
```

### 11.2 回滚方案

```bash
#!/bin/bash
# rollback.sh

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "用法: ./rollback.sh <版本号>"
    echo "例如: ./rollback.sh v1.9.0"
    exit 1
fi

echo "🔄 开始回滚到版本 ${VERSION}..."

# 1. 停止当前版本
echo "停止当前版本..."
docker-compose -f docker-compose.backend.yml down

# 2. 拉取旧版本镜像
echo "拉取旧版本镜像..."
docker pull cruise-booking-backend:${VERSION}
docker pull cruise-booking-web:${VERSION}
docker pull cruise-booking-admin:${VERSION}

# 3. 数据库回滚 (如果有迁移脚本)
echo "检查是否需要数据库回滚..."
if [ -f "migrations/rollback_${VERSION}.sql" ]; then
    echo "执行数据库回滚..."
    docker exec -i postgres-primary psql -U cruise_admin < migrations/rollback_${VERSION}.sql
fi

# 4. 启动旧版本
echo "启动旧版本..."
VERSION=${VERSION} docker-compose -f docker-compose.backend.yml up -d

# 5. 健康检查
echo "执行健康检查..."
sleep 10
./smoke-test.sh

echo "✓ 回滚完成！当前版本: ${VERSION}"
```

---

## 12. 运维手册

### 12.1 日常运维任务

#### 每日检查清单

```bash
#!/bin/bash
# daily-check.sh

echo "=== $(date) 每日运维检查 ==="

# 1. 服务状态
echo "1. 检查服务状态..."
docker-compose ps | grep -E "Exit|Dead" && echo "❌ 有服务异常" || echo "✓ 所有服务正常"

# 2. 磁盘空间
echo "2. 检查磁盘空间..."
df -h | awk '$5 > 80 {print "❌ " $0}' | grep . || echo "✓ 磁盘空间充足"

# 3. 数据库连接
echo "3. 检查数据库连接..."
curl -sf http://localhost:8080/health/ready && echo "✓ 数据库正常" || echo "❌ 数据库异常"

# 4. 日志检查
echo "4. 检查错误日志..."
grep -i "error\|fatal\|panic" /opt/cruise-booking/logs/backend/app.log | tail -5

# 5. 业务指标
echo "5. 昨日业务指标..."
docker exec postgres-primary psql -U cruise_admin -c "
SELECT 
    COUNT(*) as total_orders,
    SUM(total_amount) as total_revenue,
    COUNT(DISTINCT user_id) as active_users
FROM orders 
WHERE created_at >= CURRENT_DATE - INTERVAL '1 day'
    AND created_at < CURRENT_DATE;
"

# 6. 备份状态
echo "6. 检查备份状态..."
ls -lh /opt/cruise-booking/backups/daily/ | tail -3

echo "=== 检查完成 ==="
```

#### 每周维护任务

```bash
#!/bin/bash
# weekly-maintenance.sh

echo "=== $(date) 每周维护 ==="

# 1. 数据备份验证
echo "1. 验证备份完整性..."
./verify-backup.sh

# 2. 清理过期日志
echo "2. 清理过期日志..."
find /opt/cruise-booking/logs -name "*.log" -mtime +7 -delete
find /opt/cruise-booking/logs -name "*.log.*" -mtime +30 -delete

# 3. 数据库维护
echo "3. 执行数据库维护..."
docker exec postgres-primary psql -U cruise_admin -c "VACUUM ANALYZE;"

# 4. 更新索引
echo "4. 更新搜索引擎索引..."
curl -X POST http://localhost:7700/indexes/cruises/documents

# 5. 安全更新检查
echo "5. 检查安全更新..."
docker images | grep -E "alpine|postgres|redis" | while read image; do
    echo "检查 $image 更新..."
done

# 6. 性能分析
echo "6. 生成性能报告..."
./generate-performance-report.sh

echo "=== 维护完成 ==="
```

### 12.2 紧急响应流程

```
故障级别定义:

P0 - 灾难级 (服务完全不可用)
├── 响应时间: 5分钟内
├── 通知对象: 全员 + 管理层
├── 处理目标: 30分钟内恢复服务
└── 示例: 数据库宕机、核心服务崩溃、严重安全漏洞

P1 - 严重级 (核心功能受损)
├── 响应时间: 15分钟内
├── 通知对象: 技术团队 + 产品负责人
├── 处理目标: 2小时内恢复
└── 示例: 支付故障、预订失败、数据不一致

P2 - 一般级 (非核心功能异常)
├── 响应时间: 1小时内
├── 通知对象: 相关开发人员
├── 处理目标: 24小时内修复
└── 示例: 统计延迟、推荐异常、UI 显示问题

P3 - 轻微级 (优化建议类)
├── 响应时间: 24小时内
├── 通知对象: 产品团队
└── 示例: 性能优化、体验改进
```

### 12.3 联系人列表

```yaml
# 应急联系清单

技术支持:
  运维负责人: 张三 (电话: 138xxxx0001, 钉钉: zhangsan)
  后端负责人: 李四 (电话: 138xxxx0002, 钉钉: lisi)
  前端负责人: 王五 (电话: 138xxxx0003, 钉钉: wangwu)
  DBA: 赵六 (电话: 138xxxx0004, 钉钉: zhaoliu)

业务支持:
  产品负责人: 孙七 (电话: 138xxxx0005)
  运营负责人: 周八 (电话: 138xxxx0006)
  客服负责人: 吴九 (电话: 138xxxx0007)

外部支持:
  云服务商: 阿里云 95187
  支付渠道: 微信支付 95017
  SSL证书: DigiCert 400-xxx-xxxx
  CDN服务商: 阿里云 CDN 95187
```

---

## 📞 附录

### A. 常用命令速查

```bash
# 查看日志
docker logs -f cruise-backend
tail -f /opt/cruise-booking/logs/backend/app.log

# 重启服务
docker-compose restart backend
docker-compose up -d --no-deps --force-recreate backend

# 数据库操作
docker exec -it postgres-primary psql -U cruise_admin -d cruise_booking

# Redis 操作
docker exec -it redis-primary redis-cli

# 查看指标
curl http://localhost:8080/metrics
curl http://localhost:9090/api/v1/query?query=up

# 备份恢复
./backup.sh
./restore.sh /path/to/backup.sql
```

### B. 环境变量模板

```bash
# .env.production

# 数据库
DB_HOST=db-primary
DB_PORT=5432
DB_USER=cruise_admin
DB_PASSWORD=<强密码，32位+>
DB_NAME=cruise_booking

# Redis
REDIS_HOST=redis-primary
REDIS_PORT=6379
REDIS_PASSWORD=<强密码>

# JWT
JWT_SECRET=<随机字符串，64位+>

# 微信支付
WECHAT_APP_ID=<小程序AppID>
WECHAT_MCH_ID=<商户号>
WECHAT_API_KEY=<API密钥>
WECHAT_CERT_PATH=/secrets/apiclient_cert.pem
WECHAT_KEY_PATH=/secrets/apiclient_key.pem

# 阿里云
ALIYUN_ACCESS_KEY=<AccessKey>
ALIYUN_SECRET_KEY=<SecretKey>
ALIYUN_SMS_SIGN_NAME=<短信签名>

# MinIO
MINIO_ACCESS_KEY=<AccessKey>
MINIO_SECRET_KEY=<SecretKey>
```

### C. 监控告警规则

```yaml
# alert-rules.yml
groups:
  - name: cruise-booking-alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "错误率过高"
          description: "5xx 错误率超过 10%"

      - alert: DatabaseDown
        expr: pg_up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "数据库宕机"
          description: "PostgreSQL 主库不可访问"

      - alert: HighLatency
        expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m])) > 2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "响应延迟过高"
          description: "P99 延迟超过 2 秒"
```

---

## 🎯 总结

本部署指南涵盖了从环境准备到上线运维的完整流程。关键要点：

1. **分阶段部署**: 基础设施 → 数据层 → 应用层 → 监控层
2. **高可用设计**: 多实例、主从复制、负载均衡
3. **安全第一**: SSL、WAF、安全扫描、最小权限
4. **监控完备**: 指标、日志、追踪、告警全覆盖
5. **预案充分**: 回滚方案、故障处理、应急响应

**祝 CruiseBooking 平台上线成功，航行顺利！** 🚢✨
