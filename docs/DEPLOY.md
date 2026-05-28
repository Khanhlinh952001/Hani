# Hướng dẫn deploy Hani

Tài liệu này mô tả cách deploy **backend Go** (API + WebSocket) và **frontend Next.js** lên VPS/Linux production, dùng **PostgreSQL (pgvector)**, **PM2** và **Nginx** reverse proxy.

## Kiến trúc

```
Internet
   │
   ▼
 Nginx (443) ──► hani-fe (Next.js, PM2, port 3005)
   │                  │
   │                  └── NEXT_PUBLIC_API_URL / WS_URL
   │
   └──► hani-be (Go API, PM2, port 8080)
              │
              └── PostgreSQL (pgvector, port 5432)
```

| Thành phần | Thư mục | Port mặc định |
|------------|---------|---------------|
| Frontend | `fe/` | `3005` (PM2) |
| Backend API | `be/` | `8080` |
| PostgreSQL | Docker | `5432` |

---

## 1. Yêu cầu server

- Ubuntu 22.04+ (hoặc tương đương)
- **Go** ≥ 1.26
- **Node.js** ≥ 20 + **pnpm**
- **PM2**: `npm install -g pm2`
- **Docker** + Docker Compose (cho Postgres)
- **Nginx** + **Certbot** (HTTPS, tùy chọn nhưng khuyến nghị)

Clone repo:

```bash
git clone <repo-url> /opt/hani
cd /opt/hani
```

---

## 2. Database (PostgreSQL + pgvector)

### 2.1 Tạo file env cho Docker

Trong thư mục `be/`, tạo `.env`:

```env
POSTGRES_USER=hani
POSTGRES_PASSWORD=<mật-khẩu-mạnh>
POSTGRES_DB=hani_db
POSTGRES_PORT=5432
```

### 2.2 Chạy Postgres

```bash
cd /opt/hani/be
docker compose up -d
docker compose ps
```

Lần đầu backend chạy sẽ **AutoMigrate** bảng và **seed** plan billing + admin.

### 2.3 Backup (khuyến nghị)

```bash
docker exec postgres_db pg_dump -U hani hani_db > backup.sql
```

---

## 3. Backend (Go API)

### 3.1 Biến môi trường

Tạo `be/.env` (cùng file có thể dùng cho Docker nếu đặt chung):

```env
# Database
DB_HOST=localhost
POSTGRES_USER=hani
POSTGRES_PASSWORD=<mật-khẩu-mạnh>
POSTGRES_DB=hani_db
POSTGRES_PORT=5432

# Auth (đổi trên production!)
JWT_SECRET=<chuỗi-ngẫu-nhiên-dài>
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h

# Admin mặc định (tạo lần đầu nếu chưa có)
ADMIN_EMAIL=admin@hani.app
ADMIN_PASSWORD=<mật-khẩu-admin-mạnh>
ADMIN_BYPASS_QUOTA=false

# AI / Voice (bắt buộc cho chat + TTS/STT)
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o-mini
SONIOX_API_KEY=<soniox-key>
TTS_PROVIDER=soniox

# Tùy chọn
TRANSLATE_PROVIDER=openai
```

### 3.2 Build binary

```bash
cd /opt/hani/be
go build -o bin/api ./cmd/api
```

### 3.3 Chạy với PM2

```bash
cd /opt/hani
bash scripts/deploy-be.sh
```

Hoặc thủ công:

```bash
pm2 start ecosystem.config.cjs --only hani-be --env production
pm2 logs hani-be
```

### 3.4 Kiểm tra

```bash
curl http://127.0.0.1:8080/api/billing/plans
```

Backend phục vụ file tĩnh tại `/uploads` — giữ thư mục `be/uploads/` khi deploy (avatar, voice cache).

---

## 4. Frontend (Next.js)

### 4.1 Biến môi trường build-time

Tạo `fe/.env.production` **trước khi build** (Next.js nhúng biến `NEXT_PUBLIC_*` lúc build):

```env
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
NEXT_PUBLIC_WS_URL=wss://api.yourdomain.com
```

Nếu chạy thử trên cùng máy không HTTPS:

```env
NEXT_PUBLIC_API_URL=http://YOUR_SERVER_IP:8080
NEXT_PUBLIC_WS_URL=ws://YOUR_SERVER_IP:8080
```

### 4.2 Build + PM2

```bash
cd /opt/hani
bash scripts/deploy-fe.sh
```

Script sẽ: `pnpm install` → `pnpm build` → `pm2 start/reload hani-fe`.

Port mặc định: **3005** (xem `ecosystem.config.cjs`).

```bash
pm2 status
pm2 logs hani-fe
curl -I http://127.0.0.1:3005
```

### 4.3 Cập nhật phiên bản mới

```bash
cd /opt/hani
git pull
bash scripts/deploy-fe.sh   # FE
bash scripts/deploy-be.sh   # BE (nếu có thay đổi API)
```

---

## 5. PM2 — lệnh thường dùng

```bash
pm2 status
pm2 logs hani-fe
pm2 logs hani-be
pm2 restart hani-fe
pm2 restart hani-be
pm2 save                  # lưu danh sách process
pm2 startup               # tự chạy lại khi reboot (chạy lệnh PM2 in ra)
```

Deploy cả hai:

```bash
bash scripts/deploy-be.sh && bash scripts/deploy-fe.sh
```

---

## 6. Nginx reverse proxy

Ví dụ domain:

- `hani.yourdomain.com` → frontend (3005)
- `api.yourdomain.com` → backend (8080), kèm WebSocket

```nginx
# /etc/nginx/sites-available/hani

upstream hani_fe {
  server 127.0.0.1:3005;
}

upstream hani_be {
  server 127.0.0.1:8080;
}

server {
  listen 80;
  server_name hani.yourdomain.com;
  return 301 https://$host$request_uri;
}

server {
  listen 443 ssl http2;
  server_name hani.yourdomain.com;

  ssl_certificate     /etc/letsencrypt/live/hani.yourdomain.com/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/hani.yourdomain.com/privkey.pem;

  location / {
    proxy_pass http://hani_fe;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }
}

server {
  listen 80;
  server_name api.yourdomain.com;
  return 301 https://$host$request_uri;
}

server {
  listen 443 ssl http2;
  server_name api.yourdomain.com;

  ssl_certificate     /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;

  client_max_body_size 20M;

  location / {
    proxy_pass http://hani_be;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }

  # WebSocket chat (/api/ws/chat)
  location /api/ws/ {
    proxy_pass http://hani_be;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host $host;
    proxy_read_timeout 3600s;
    proxy_send_timeout 3600s;
  }
}
```

Kích hoạt:

```bash
sudo ln -s /etc/nginx/sites-available/hani /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```

HTTPS (Let's Encrypt):

```bash
sudo certbot --nginx -d hani.yourdomain.com -d api.yourdomain.com
```

Sau khi có HTTPS, **build lại frontend** với `NEXT_PUBLIC_*` trỏ `https://` / `wss://`.

---

## 7. Admin & gói subscription

| Mục | Giá trị |
|-----|---------|
| URL admin | `https://hani.yourdomain.com/admin` |
| Tài khoản mặc định | `ADMIN_EMAIL` / `ADMIN_PASSWORD` trong `be/.env` |
| Đổi gói user | Admin → chọn user → Free / Plus / Premium |
| Reset hạn mức | Nút **Reset hạn mức** (tin nhắn & voice trong ngày) |

Gói áp dụng ngay từ database; user xem hạn mức tại **Cài đặt → Gói & hạn mức**.

Chi tiết kiến trúc SaaS: [`be/docs/SAAS_PLATFORM.md`](../be/docs/SAAS_PLATFORM.md).

---

## 8. Checklist production

- [ ] Đổi `JWT_SECRET`, `ADMIN_PASSWORD`, mật khẩu Postgres
- [ ] `ADMIN_BYPASS_QUOTA=false`
- [ ] `fe/.env.production` dùng URL HTTPS/WSS đúng domain API
- [ ] Build lại FE sau khi đổi `NEXT_PUBLIC_*`
- [ ] Firewall: mở `80`, `443`; **không** expose `5432`, `8080`, `3005` ra internet (chỉ localhost + Nginx)
- [ ] `pm2 save` + `pm2 startup`
- [ ] Backup database định kỳ

---

## 9. Xử lý lỗi thường gặp

| Triệu chứng | Cách xử lý |
|-------------|------------|
| FE gọi API lỗi CORS / 404 | Kiểm tra `NEXT_PUBLIC_API_URL`, build lại FE |
| WebSocket không kết nối | Nginx cần block `location` upgrade; `NEXT_PUBLIC_WS_URL` phải `wss://` khi dùng HTTPS |
| `profile already exists` | User đã có profile; vào app hoặc admin xóa/sync — xem commit fix lover upsert |
| Hết quota dù đã lên Plus | Admin → **Reset hạn mức** |
| Go IDE báo `missing metadata` | Mở workspace root `Hani/` (có `go.work` trỏ `./be`) |
| PM2 FE crash ngay | Chạy `pnpm build` trong `fe/` trước; xem `pm2 logs hani-fe` |

---

## 10. Dev local (tham khảo)

```bash
# DB
cd be && docker compose up -d

# API
cd be && go run ./cmd/api

# FE
cd fe && pnpm dev
```

Mặc định: FE `http://localhost:3000`, API `http://localhost:8080`.
