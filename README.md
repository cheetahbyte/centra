# What is this?

This is a dead-simple headless CMS. No fancy dashboards, no bloated page builders, no “enterprise-ready omnichannel synergy”.
Just a clean backend that stores your content and serves it fast. That’s it.

# Why?
Most headless CMS feel like they’re built for marketing teams with five layers of approval and a mandatory “Content Strategist” role. I don’t need that.
I just want something lightweight that stays out of my way, doesn’t force workflows on me, and simply delivers structured data to a couple of websites and apps.
This project is exactly that:
> A minimal, predictable CMS that does its job without drama.

# Setup
Read [here](https://github.com/cheetahbyte/centra/wiki/Setup)

# Environment Variables
There a few env vars that can be configured

| Variable | Description | Default |
|----------|-------------|----------|
| `PORT` | sets the port on which the webserver starts | `3000` |
| `CENTRA_API_KEY` | provided a string here will set the key used to protect the content api  | _none_ |
| `CONTENT_ROOT` | sets the directory where the server will look for content | `/content` |
| `KEYS_DIR` | sets the directory where the server will place the ssh keys | `/keys` |
| `GITHUB_REPO_URL` | sets the GitHub repo from which the content is served | _none_ |
| `CORS_ALLOWED_ORIGINS` | list of allowed origins | `*` |
| `CORS_ALLOWED_METHODS` | list of allowed methods | `["GET","HEAD","OPTIONS"]` |
| `CORS_ALLOWED_HEADERS` | list of allowed request headers | `*` |
| `CORS_EXPOSED_HEADERS` | list of headers exposed to the browser | `["Cache-Control","Content-Language","Content-Length","Content-Type","Expires","Last-Modified"]` |
| `CORS_ALLOW_CREDENTIALS` | whether credentials (cookies/auth headers) are allowed | `false` |
| `CORS_MAX_AGE` | max age of preflight cache (in seconds) | `360` |
| `SSH_PRIVATE_KEY` | ssh private key used for communication with git | _none_ |
| `SSH_PUBLIC_KEY` | ssh public key used for communication with git | _none_ |
| `LOG_LEVEL` | sets the log level | `INFO` |
| `LOG_STRUC` | turns off pretty printing of the logs and logs in plain json  | `false` |
| `WEBHOOK_SECRET` | secret used for the webhook secret validation (providing none will disable it) | _none_ |
| `RATELIMIT_QUOTA` | decides how many requests a client can send during a minute (`time.Minute`) | `100` |
