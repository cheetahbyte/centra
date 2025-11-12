FROM golang:1.21-alpine AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o /app/main .

FROM scratch

COPY --from=builder /app/main /main

EXPOSE 8080

ENTRYPOINT ["/main"]
