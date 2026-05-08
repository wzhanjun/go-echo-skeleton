# go-echo-skeleton

## Commands

| Command | What it does |
|---------|-------------|
| `make generate` | `go generate ./internal/enum` — regenerates `internal/enum/errcode_string.go` from stringer directive |
| `make lint` | `golangci-lint run` |
| `make all` | `make generate` **then** `make lint` — run this before committing |
| `make docs` | `swag init` (swag is configured but not wired to code yet) |
| `make win64` | Cross-build Windows binary to `bin/go-echo-skeleton.exe` |
| `make linux-amd64` | Cross-build Linux binary with UPX compression to `bin/go-echo-skeleton` |
| `go run ./cmd/server -c ../../config` | Run locally from repo root |

**Order matters**: `make generate` must run before `make lint` (enforced by `make all`).

## Architecture

- **Entrypoint**: `cmd/server/main.go` — Echo v4 HTTP server
- **Config**: Viper reads YAML from `-c <dir>` flag (default `../../config`). Config hot-reloads via fsnotify. Copy `config/config.yaml.example` to `config/config.yaml`. Environment variables (e.g. `SYSTEM_ADDR`) override config keys via `AutomaticEnv()`. `.env` file at project root is loaded automatically.
- **Router**: `internal/routers/router.go` — sets up Echo with recover, CORS, and request-time middleware
- **Handlers**: builder pattern (`NewXxxController() *xxxController`) embedding `handler.BaseController`. Response envelope: `{code, msg, data}`.
- **Errors**: strongly typed via `enum.ErrCode` (int32) and `enum.ApiError{Code, Msg}`. Handled in `handler/base.go` with type switch.
- **DB**: XORM + MySQL, lazy singleton via `db.GetEngine()`. All repo methods take `*xorm.Session`.
- **Cache**: Redis, lazy singleton with sentinel failover support. `cache.GetRedisLock` / `cache.RedisUnLock` for distributed locks.
- **Cron**: `robfig/cron` v3 with seconds support. Jobs registered in `internal/jobs/start_job.go`. Enabled via config key `system.start-cron`.
- **Repo layer**: generic `repo.BaseRepo[T any]` interface with `repo/impl/base.go` — uses XORM sessions + optional Redis cache.
- **`internal/models/`** and `internal/services/` are intentionally empty stubs (only `.gitignore`).
- **Logging**: custom `log-service` client (slog wrapper), not standard library `log`.

## Conventions

- Error codes use `go:generate stringer -type ErrCode -linecomment`. Add new codes in `internal/enum/errcode.go`, then run `make generate`. The generated `errcode_string.go` **must** be committed.
- API handlers: return `s.Success(c, data)` or `s.Fail(c, err)`. Business errors return `enum.ApiError{Code: enum.Xxx, Msg: "..."}`.
- Handlers accept `echo.Context` — parameter binding errors are caught and converted to `ParamsError` automatically.
- Date/time JSON binding: use `dto.Date` and `dto.DateTime` types (custom JSON marshal/unmarshal with Chinese date formats).
- Config struct tags use `mapstructure` for Viper compatibility.

## Config

- Requires `config/config.yaml` (gitignored, use `config/config.yaml.example` as template)
- `.env` file at project root (with `APP_ID`, `APP_DEBUG`) is loaded via godotenv before Viper reads YAML. Env vars can override any config key (e.g. `SYSTEM_ADDR=0.0.0.0:9999` overrides `system.addr`).
- Key config sections: `system` (env, addr, location, start-cron, show-sql), `mysql`, `redis` (incl. sentinel), `jwt` (secret, expires, issuer), `email`
