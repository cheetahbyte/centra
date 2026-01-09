# Configuration

All configuration values are provided via environment variables.

| Category | Env Variable | Type | Default | Description |
|--------|--------------|------|---------|-------------|
| Core | `PORT` | string | `3000` | Port the HTTP server listens on |
| Core | `CONTENT_ROOT` | string | `/content` | Root directory for all content files |
| Core | `CENTRA_API_KEY` | string | – | Optional API key for protected/internal endpoints |
| Git & Keys | `GITHUB_REPO_URL` | string | – | Git repository URL used for content sync |
| Git & Keys | `KEYS_DIR` | string | `/keys` | Directory where SSH keys are stored |
| Git & Keys | `SSH_PRIVATE_KEY` | string | – | Private SSH key for Git access |
| Git & Keys | `SSH_PUBLIC_KEY` | string | – | Public SSH key corresponding to the private key |
| Git & Keys | `WEBHOOK_SECRET` | string | – | Secret for validating incoming webhooks |
| CORS | `CORS_ALLOWED_ORIGINS` | string[] | `*` | Allowed CORS origins (`*` allows all) |
| CORS | `CORS_ALLOWED_METHODS` | string[] | `GET,HEAD,OPTIONS` | Allowed HTTP methods for CORS |
| CORS | `CORS_ALLOWED_HEADERS` | string[] | `*` | Allowed request headers for CORS |
| CORS | `CORS_EXPOSED_HEADERS` | string[] | `Cache-Control,Content-Language,Content-Type,Expires,Last-Modified` | Response headers exposed to the browser |
| CORS | `CORS_MAX_AGE` | int | `360` | Time (seconds) a CORS preflight response is cached |
| CORS | `CORS_ALLOW_CREDENTIALS` | bool | `false` | Allow credentials in CORS requests |
| Logging & Limits | `LOG_LEVEL` | string | `INFO` | Log level (`DEBUG`, `INFO`, `WARN`, `ERROR`) |
| Logging & Limits | `LOG_STRUC` | bool | `false` | Enable structured (e.g. JSON) logging |
| Logging & Limits | `RATELIMIT_QUOTA` | int | `100` | Maximum requests per time window |
| Features | `CACHE_BINARIES` | bool | `false` | Enable caching of binary files |
| Features | `ALLOWED_BINARIES` | string[] | `*` | Allowed binary file extensions (`*` allows all) |
| Features | `AnyBinaries` | internal | derived | Automatically set when `ALLOWED_BINARIES=*`; not configurable |

## Example

```bash
PORT=3000
CONTENT_ROOT=/content
GITHUB_REPO_URL=git@github.com:org/repo.git
CACHE_BINARIES=true
ALLOWED_BINARIES=.jpg,.png,.pdf
LOG_LEVEL=INFO
```
