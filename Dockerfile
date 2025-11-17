FROM oven/bun:alpine AS builder
WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --production

COPY . .

# Stage 2: Runtime
FROM oven/bun:alpine AS runner
WORKDIR /app

COPY --from=builder /app /app

VOLUME ["/content"]

ENV NODE_ENV=production
ENV PORT=3000

# Expose port
EXPOSE 3000

CMD ["bun", "run", "src/app.ts"]
