# builder
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o centra-server ./src/cmd/centra-server

# runtime
FROM scratch

WORKDIR /app

VOLUME ["/content"]

COPY --from=builder /app/centra-server /app/centra-server

EXPOSE 3000

ENTRYPOINT ["/app/centra-server"]
