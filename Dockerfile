# --- Builder Stage -----------------------------------------------------------
FROM golang:1.25-bookworm AS builder

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates \
      tzdata && \
    rm -rf /var/lib/apt/lists/*

RUN useradd \
    --uid 10001 \
    --no-create-home \
    --home-dir /nonexistent \
    --shell /usr/sbin/nologin \
    appuser

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o centra-server ./cmd/centra-server


# --- Runtime Stage -----------------------------------------------------------
FROM debian:12-slim AS runtime

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      git \
      openssh-client \
      ca-certificates \
      tzdata && \
    rm -rf /var/lib/apt/lists/*


# make friends with github
RUN mkdir -p /etc/ssh && \
    ssh-keyscan github.com >> /etc/ssh/ssh_known_hosts

# Prepare a writable keys dir for appuser
RUN mkdir -p /keys && chown 10001:10001 /keys
RUN mkdir -p /content && chown 10001:10001 /content

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/centra-server /app/centra-server

ENV KEYS_DIR=/keys

USER appuser:appuser

VOLUME ["/keys", "/content"]

EXPOSE 3000

ENTRYPOINT ["/app/centra-server"]
