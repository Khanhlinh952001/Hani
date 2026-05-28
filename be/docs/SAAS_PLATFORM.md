# Hani SaaS Platform — Production Architecture

> **Status:** Architecture spec (v1) — aligned with existing Go/Gin + Next.js codebase.  
> **Principle:** No unlimited free AI. Every chat/voice/STT/TTS path must pass quota checks.

---

## 0. Current state (baseline)

| Area | Exists today | Gap |
|------|----------------|-----|
| Auth | Email/password, JWT (7d, no refresh), bcrypt | OAuth, refresh rotation, sessions/devices |
| Roles | `users.role` (0=user, 1=admin) | `guest`, plan-based entitlements |
| Admin | `/api/admin/*` stats, users, sessions, memories | Analytics, billing, moderation, usage reset |
| WS | JWT on upgrade, per-user session | Quota gate, connection limits, guest deny |
| Billing | — | subscriptions, payments, usage counters |
| Redis | — | rate limits, quota cache, session store |
| FE admin | Next.js `AdminView` | Full SaaS dashboard (charts, finance, moderation) |

**Existing packages to extend (do not rewrite):**

- `be/internal/auth` — JWT, middleware, register/login
- `be/internal/modules/users` — user entity, CRUD
- `be/internal/admin` — admin API + `RequireAdmin`
- `be/internal/websocket` — realtime pipeline (quota hooks go here)
- `fe/app/admin`, `fe/lib/admin/api.ts` — admin UI shell

---

## 1. Auth architecture

### 1.1 Token model (dual-token)

```
┌─────────────┐     login/register/OAuth     ┌──────────────────┐
│   Client    │ ───────────────────────────► │  Auth Service    │
└─────────────┘                              └────────┬─────────┘
       │                                              │
       │  access_token (JWT, 15m)                      │ persist
       │  refresh_token (opaque, 30d, rotatable)        ▼
       │◄────────────────────────────────────  refresh_tokens
       │                                              sessions
       └─ Authorization: Bearer <access>              devices
```

**Access JWT claims (extend current `auth.Claims`):**

```go
type Claims struct {
    UserID   int    `json:"uid"`
    Email    string `json:"email"`
    Name     string `json:"name"`
    Role     int    `json:"role"`      // 0 user, 1 admin
    Plan     string `json:"plan"`      // free | plus | premium
    Guest    bool   `json:"guest"`     // true = guest session only
    SessionID string `json:"sid"`     // server session row id
    jwt.RegisteredClaims // exp ~15m
}
```

**Refresh token:** random 32-byte hex, stored hashed (SHA-256) in `refresh_tokens`, bound to `session_id` + `device_id`. Rotation on every refresh; reuse detection → revoke family.

### 1.2 Endpoints (`/api/auth`)

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/register` | — | Email/password (existing) |
| POST | `/login` | — | Returns access + refresh |
| POST | `/refresh` | refresh cookie/body | Rotate refresh, new access |
| POST | `/logout` | Bearer | Revoke session + refresh |
| POST | `/logout-all` | Bearer | Revoke all devices |
| GET | `/me` | Bearer | Profile + plan + usage summary |
| POST | `/oauth/google` | — | ID token → user + tokens |
| POST | `/oauth/apple` | — | identity token → user + tokens |
| POST | `/guest` | — | Issue guest access (limited claims) |

### 1.3 OAuth flow

1. Client obtains provider ID token (Google Sign-In / Sign in with Apple).
2. `POST /api/auth/oauth/{provider}` validates token with provider JWKS.
3. `users.FindOrCreateOAuth(provider, sub, email, name)`.
4. Create `sessions` row + `devices` row + tokens.

### 1.4 Session & device tracking

```sql
-- sessions: logical login session
CREATE TABLE sessions (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       INT REFERENCES users(id),
  guest_id      UUID,                    -- nullable; set for guest
  ip_hash       TEXT,
  user_agent    TEXT,
  revoked_at    TIMESTAMPTZ,
  last_seen_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE devices (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       INT REFERENCES users(id),
  fingerprint   TEXT NOT NULL,         -- client-generated stable id
  platform      TEXT,                    -- ios | android | web
  push_token    TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, fingerprint)
);

CREATE TABLE refresh_tokens (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id    UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  token_hash    TEXT NOT NULL UNIQUE,
  expires_at    TIMESTAMPTZ NOT NULL,
  replaced_by   UUID REFERENCES refresh_tokens(id),
  revoked_at    TIMESTAMPTZ,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 1.5 Security controls

| Control | Implementation |
|---------|----------------|
| Password | bcrypt (existing), min 6 → raise to 8 in prod |
| Rate limit | Redis: `rl:login:{ip}` 10/min, `rl:register:{ip}` 5/hour |
| Anti-spam | honeypot field, email domain blocklist |
| API protection | `RequireAuth` + `RequirePlan` + `RequireQuota` middleware chain |
| Audit | `admin_logs` + `security_events` |

**Package layout:**

```
be/internal/auth/
  jwt.go          # access token (extend claims)
  refresh.go      # refresh CRUD + rotation
  oauth_google.go
  oauth_apple.go
  guest.go
  middleware.go   # RequireAuth, OptionalAuth
  session.go

be/internal/modules/identity/   # optional split later
  repository.go
  service.go
```

---

## 2. User types & entitlements

### 2.1 Identity matrix

| Type | `users` row | `subscription_plan` | Chat | Voice | Memory | Save companion |
|------|-------------|---------------------|------|-------|--------|----------------|
| Guest | no / `guests` table | — | 5/day | 0 | no | no |
| Free | yes | `free` | 30/day | 5 min/day | basic | yes |
| Plus | yes | `plus` | 1000/day | 60 min/day | full | yes |
| Premium | yes | `premium` | unlimited | unlimited | advanced | yes |
| Admin | yes, `role=1` | any | bypass* | bypass* | bypass* | yes |

\*Admin bypass only when `ADMIN_BYPASS_QUOTA=true` in env (default false in prod).

### 2.2 Plan limits (source of truth)

```sql
CREATE TABLE plan_limits (
  plan              TEXT PRIMARY KEY,  -- guest | free | plus | premium
  daily_messages    INT,               -- NULL = unlimited
  daily_voice_sec   INT,
  monthly_messages  INT,
  max_memories      INT,
  max_companions    INT DEFAULT 1,
  allow_voice       BOOLEAN NOT NULL,
  allow_memory      BOOLEAN NOT NULL,
  allow_premium_voices BOOLEAN NOT NULL,
  ai_model_tier     TEXT NOT NULL,     -- basic | standard | fast
  tts_tier          TEXT NOT NULL,
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO plan_limits VALUES
  ('guest',   5,    0,    5,    0, 0, false, false, false, 'basic',    'none'),
  ('free',   30,  300,  900,   50, 1, false, true,  false, 'basic',    'basic'),
  ('plus', 1000, 3600, NULL, 500, 1, true,  true,  true,  'standard', 'premium'),
  ('premium', NULL, NULL, NULL, NULL, 3, true, true, true, 'fast', 'premium');
```

Entitlement resolver:

```go
// be/internal/billing/entitlements.go
func PlanForUser(u *users.User, sub *Subscription) string
func LimitsForPlan(plan string) PlanLimits
func CanUseFeature(plan, feature string) bool
```

---

## 3. Usage limit system

### 3.1 Tables

```sql
CREATE TABLE user_usage (
  id                  BIGSERIAL PRIMARY KEY,
  user_id             INT REFERENCES users(id),
  guest_id            UUID,
  period_date         DATE NOT NULL,           -- UTC day for daily counters
  daily_messages      INT NOT NULL DEFAULT 0,
  daily_voice_seconds INT NOT NULL DEFAULT 0,
  monthly_messages    INT NOT NULL DEFAULT 0,
  tokens_in           BIGINT NOT NULL DEFAULT 0,
  tokens_out          BIGINT NOT NULL DEFAULT 0,
  tts_characters      BIGINT NOT NULL DEFAULT 0,
  stt_seconds         INT NOT NULL DEFAULT 0,
  embedding_calls     INT NOT NULL DEFAULT 0,
  estimated_cost_usd  NUMERIC(12,6) NOT NULL DEFAULT 0,
  reset_at            TIMESTAMPTZ NOT NULL,
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, period_date),
  UNIQUE (guest_id, period_date)
);

CREATE TABLE usage_logs (
  id          BIGSERIAL PRIMARY KEY,
  user_id     INT,
  guest_id    UUID,
  event_type  TEXT NOT NULL,  -- chat_message | voice_start | voice_end | tts | stt | embed
  units       INT NOT NULL,
  metadata    JSONB,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_usage_logs_user_created ON usage_logs(user_id, created_at DESC);
```

### 3.2 Quota check flow (every AI action)

```
Request → ResolveIdentity → ResolvePlan → LoadUsage(day) → CompareLimits
   │                                              │
   ├─ HARD_BLOCK (402/429) ◄── exceeded           │
   ├─ SOFT_WARN (header X-Quota-Warning) ◄── 80% │
   └─ ALLOW → handler → RecordUsage (async)       │
```

**Redis cache (hot path):**

```
quota:{userId}:{YYYY-MM-DD}:messages  → INT
quota:{userId}:{YYYY-MM-DD}:voice_sec → INT
```

Increment with `INCR` + `EXPIRE` at end of day; flush to Postgres every N minutes or on threshold.

### 3.3 Reset policy

| Counter | Reset |
|---------|--------|
| `daily_messages`, `daily_voice_seconds` | UTC midnight cron / lazy on first request of new day |
| `monthly_messages` | 1st of month UTC |
| Subscription validity | `subscriptions.expires_at` checked on each request |

### 3.4 Middleware

```go
// be/internal/billing/middleware.go
func RequireQuota(event UsageEvent) gin.HandlerFunc
func QuotaHeaders(c *gin.Context, usage UsageSnapshot)
```

Apply to:

- `POST` chat REST (if any)
- WebSocket: before `stream_turn`, on `final_transcript` (voice), on TTS chunk emit
- `POST /api/soniox/temporary-key` (STT quota)
- Lover profile create (guest/free companion limits)

### 3.5 Soft warnings & upsell

Response headers:

```
X-Quota-Plan: free
X-Quota-Messages-Used: 24
X-Quota-Messages-Limit: 30
X-Quota-Warning: 80
```

WS event:

```json
{ "type": "quota_warning", "feature": "messages", "used": 24, "limit": 30, "upgrade_url": "/premium" }
```

---

## 4. Subscription & payment-ready architecture

### 4.1 Schema

```sql
CREATE TYPE subscription_status AS ENUM (
  'trialing', 'active', 'past_due', 'canceled', 'expired'
);

CREATE TABLE subscriptions (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id           INT NOT NULL REFERENCES users(id),
  provider          TEXT NOT NULL,  -- stripe | app_store | play_store | manual
  provider_sub_id   TEXT,
  plan              TEXT NOT NULL REFERENCES plan_limits(plan),
  status            subscription_status NOT NULL,
  started_at        TIMESTAMPTZ NOT NULL,
  expires_at        TIMESTAMPTZ,
  canceled_at       TIMESTAMPTZ,
  metadata          JSONB,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payments (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id           INT NOT NULL REFERENCES users(id),
  subscription_id   UUID REFERENCES subscriptions(id),
  amount_cents      INT NOT NULL,
  currency          TEXT NOT NULL DEFAULT 'usd',
  provider          TEXT NOT NULL,
  provider_payment_id TEXT,
  status            TEXT NOT NULL,  -- pending | succeeded | failed | refunded
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE invoices (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id           INT NOT NULL,
  subscription_id   UUID,
  provider_invoice_id TEXT,
  amount_cents      INT NOT NULL,
  pdf_url           TEXT,
  status            TEXT NOT NULL,
  period_start      TIMESTAMPTZ,
  period_end        TIMESTAMPTZ,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 4.2 `users` migration (extend existing)

```sql
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS subscription_plan TEXT NOT NULL DEFAULT 'free',
  ADD COLUMN IF NOT EXISTS subscription_status TEXT NOT NULL DEFAULT 'active',
  ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true,
  ADD COLUMN IF NOT EXISTS banned_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS stripe_customer_id TEXT;
```

Keep `role` for admin; plan is separate from role.

### 4.3 Webhook handlers (stubs, payment-ready)

```
POST /api/webhooks/stripe
POST /api/webhooks/app-store
POST /api/webhooks/play-store
```

Each validates signature → updates `subscriptions` + `users.subscription_plan` → emits `billing.subscription_changed`.

### 4.4 Client APIs

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/billing/plans` | Public plan cards |
| GET | `/api/billing/subscription` | Current sub |
| POST | `/api/billing/checkout` | Stripe Checkout session URL |
| POST | `/api/billing/portal` | Stripe customer portal |
| GET | `/api/billing/usage` | Today + month usage |

---

## 5. WebSocket protection

### 5.1 Upgrade gate (`websocket.HandleChat`)

1. Parse JWT (existing).
2. Reject `guest` if `practice_mode` voice (query `practice_mode != chat`).
3. `billing.CheckQuota(userID, UsageEventChat)` — reject upgrade with 402 JSON before upgrade.
4. Limit concurrent connections: Redis `SET ws:conn:{userId}` with TTL; max 2 (free), 5 (premium).
5. Attach `QuotaContext` to `RealtimeSession`.

### 5.2 Per-turn metering (`websocket/stream_turn.go`)

After AI completes:

```go
usage.Record(users.ID, UsageEvent{
  Type: "chat_message",
  TokensIn: resp.Usage.PromptTokens,
  TokensOut: resp.Usage.CompletionTokens,
})
```

Voice session: accumulate seconds in session struct; on disconnect flush `voice_seconds`.

### 5.3 Guest WebSocket (optional, strict)

- Separate path `/api/ws/guest` or claim `guest: true` in JWT from `POST /api/auth/guest`.
- No memory retrieval, no lover profile, model tier `basic`, hard 5 messages.

---

## 6. API structure (full map)

```
/api
├── auth/           # login, register, refresh, oauth, guest, logout
├── billing/        # plans, subscription, checkout, usage
├── webhooks/       # stripe, app_store, play_store
├── users/          # (internal) — keep minimal; prefer /auth/me
├── characters/     # existing — gate premium presets
├── lover/          # existing — gate create by plan
├── sessions/       # existing
├── messages/       # existing
├── memories/       # existing — gate by plan_limits.allow_memory
├── soniox/         # STT key — quota middleware
├── ws/chat         # existing — quota + conn limit
├── admin/
│   ├── stats       # extend: DAU, MAU, cost
│   ├── users       # ban, plan override, reset usage
│   ├── analytics/  # retention, voice, tokens
│   ├── moderation/ # flagged content
│   ├── billing/    # MRR, churn
│   └── system/     # WS health, latency
└── notifications/  # quota warnings, upsell (phase 2)
```

---

## 7. Admin dashboard architecture

### 7.1 Backend (`be/internal/admin` → split modules)

```
admin/
  users/        # search, ban, plan patch, reset quota
  analytics/    # DAU/MAU, retention SQL + Redis
  moderation/   # flags, review queue
  billing/      # revenue aggregates (from payments)
  system/       # health, WS connections, AI latency
```

**New endpoints (examples):**

```
GET  /api/admin/analytics/overview?from=&to=
GET  /api/admin/analytics/retention
GET  /api/admin/analytics/ai-cost
GET  /api/admin/users?q=&plan=&status=
PATCH /api/admin/users/:id  { plan, is_active, role }
POST  /api/admin/users/:id/reset-usage
GET  /api/admin/moderation/flags
POST /api/admin/moderation/flags/:id/resolve
GET  /api/admin/billing/revenue
```

### 7.2 Frontend

| App | Stack | Notes |
|-----|-------|-------|
| Mobile | Flutter (future) | Consumes same APIs |
| Admin | **React/Next.js** (`fe/app/admin`) | Extend `AdminView` → multi-page dashboard |

**Admin pages:**

- `/admin` — overview cards (DAU, MAU, messages today, AI cost estimate)
- `/admin/users` — table + actions
- `/admin/analytics` — charts (recharts)
- `/admin/moderation` — flagged conversations
- `/admin/billing` — MRR, subscriptions
- `/admin/system` — WS stats, latency percentiles

Dark/light via existing Tailwind + `next-themes`.

### 7.3 Moderation pipeline

```
User message → AI response
     │
     ├─ Rule engine (regex, rate)
     ├─ Optional: OpenAI moderation API
     └─ Flag → moderation_flags table
              → admin queue
```

```sql
CREATE TABLE moderation_flags (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       INT,
  session_id    UUID,
  message_id    UUID,
  reason        TEXT NOT NULL,
  severity      TEXT NOT NULL,  -- low | medium | high
  status        TEXT NOT NULL DEFAULT 'open',
  reviewed_by   INT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## 8. AI cost control

### 8.1 Token tracking

Wrap `ai/chat.go` stream:

```go
type UsageRecorder interface {
  RecordTokens(userID int, in, out int, model string, costUSD float64)
}
```

Store in `user_usage` + `usage_logs`; aggregate daily in admin.

### 8.2 Cost optimization

| Technique | Where |
|-----------|--------|
| Memory summarization | `memory/context.go` — compress old turns |
| Embedding cache | Redis `emb:{hash}` TTL 24h |
| Context window cap | `maxRecentTurns` by plan tier |
| Response compression | truncate TTS input (`tts/sanitize.go`) |
| Vector top-K by plan | free: 3, premium: 8 |

### 8.3 Abuse prevention

- Max message length by plan
- Min interval between messages (Redis `last_msg:{userId}`)
- Duplicate prompt detection (hash in Redis, 1 min window)
- Blocked IP list in Redis

---

## 9. Infrastructure

```
┌──────────────┐     ┌─────────────┐     ┌──────────────┐
│  Next.js FE  │────►│  Gin API    │────►│  PostgreSQL  │
│  Flutter     │     │  + WS Hub   │     │  (source)    │
└──────────────┘     └──────┬──────┘     └──────────────┘
                            │
                     ┌──────▼──────┐
                     │    Redis    │
                     │ quota, rl,  │
                     │ sessions    │
                     └─────────────┘
```

**Docker Compose services:** `api`, `postgres`, `redis`, optional `worker` (cron resets, usage flush).

**Env vars:**

```
JWT_SECRET, JWT_ACCESS_TTL=15m, JWT_REFRESH_TTL=720h
REDIS_URL
STRIPE_SECRET_KEY, STRIPE_WEBHOOK_SECRET
GOOGLE_OAUTH_CLIENT_ID
APPLE_CLIENT_ID
ADMIN_BYPASS_QUOTA=false
```

---

## 10. Premium UX (mobile + web)

| Surface | Behavior |
|---------|----------|
| Locked voice | Gray card + “Plus” badge |
| Locked personality | `plan_required: plus` from API |
| Soft paywall | At 80% quota — banner + CTA |
| Hard paywall | 100% — modal, WS `quota_exceeded` |
| Copy example | “Unlock deeper emotional memory with Premium 💕” |

**API field on gated resources:**

```json
{ "id": "nina", "locked": true, "required_plan": "plus" }
```

---

## 11. Notification system (phase 2)

```
notifications/
  service.go     # enqueue
  templates.go   # quota_warning, premium_upsell, streak, expiry
```

Channels: push (FCM/APNs), email (SendGrid), in-app inbox table.

---

## 12. Implementation phases

### Phase 1 — Foundation (1–2 weeks)

- [ ] DB migrations: `plan_limits`, `user_usage`, `usage_logs`, extend `users`
- [ ] Redis client + `billing` package
- [ ] `RequireQuota` on WS + record message usage
- [ ] Refresh tokens + sessions tables
- [ ] Guest auth endpoint + 5 msg/day

### Phase 2 — Subscriptions (1–2 weeks)

- [ ] `subscriptions`, `payments` tables
- [ ] Stripe checkout + webhook stub
- [ ] `GET /api/billing/usage`, plan resolver
- [ ] Voice seconds tracking

### Phase 3 — Admin & moderation (1–2 weeks)

- [ ] Admin analytics endpoints
- [ ] User ban / plan override / reset usage
- [ ] `moderation_flags` + basic rules
- [ ] Admin UI pages (charts)

### Phase 4 — OAuth & hardening (1 week)

- [ ] Google + Apple login
- [ ] Rate limits, IP throttle
- [ ] Token rotation abuse detection
- [ ] AI cost dashboard

### Phase 5 — Mobile & payments (ongoing)

- [ ] Flutter app
- [ ] App Store / Play billing webhooks
- [ ] Push notifications

---

## 13. Error codes (client contract)

| HTTP | Code | Meaning |
|------|------|---------|
| 401 | `unauthorized` | Missing/invalid token |
| 402 | `quota_exceeded` | Hard limit hit |
| 403 | `plan_required` | Feature needs upgrade |
| 403 | `account_banned` | `is_active=false` |
| 429 | `rate_limited` | Too many requests |
| 503 | `ai_unavailable` | Upstream AI down |

WS close code `4402` + JSON body for quota exceeded.

---

## 14. Files to add (suggested tree)

```
be/internal/
  billing/
    entitlements.go
    usage.go
    quota.go
    middleware.go
    repository.go
    plans.go
  payments/
    stripe.go
    webhook.go
  oauth/
    google.go
    apple.go
  moderation/
    flags.go
    scanner.go
  analytics/
    dau.go
    retention.go
  platform/redis.go

be/migrations/
  001_plan_limits.sql
  002_usage.sql
  003_subscriptions.sql
  004_auth_sessions.sql
```

---

## 15. Definition of done (production SaaS)

- [ ] No API path generates AI without quota check
- [ ] Guest cannot exceed 5 messages/day
- [ ] Free cannot exceed 30 messages / 5 voice min
- [ ] Premium features return `403 plan_required` for free users
- [ ] Refresh token rotation enabled; stolen token reuse revokes family
- [ ] Admin can ban user and reset usage
- [ ] Daily AI cost visible in admin
- [ ] Stripe webhook updates plan within 60s
- [ ] WS concurrent connection limit enforced

---

*Document version: 2026-05-28. Update as migrations land.*
