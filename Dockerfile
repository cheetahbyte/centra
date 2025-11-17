# Stage 1: Build (if you transpile/TS, etc.)
FROM oven/bun:1 AS builder
WORKDIR /app

COPY package.json bun.lockb ./
RUN bun install --production

COPY . .
# If you have a build step (TS → JS, etc.)
# RUN bun run build

# Stage 2: Runtime
FROM oven/bun:1 AS runner
WORKDIR /app

# Copy node_modules (and build output if needed)
COPY --from=builder /app /app

# Folder where Centra expects YAML content
# This will be a volume mount in K8s
VOLUME ["/content"]

ENV NODE_ENV=production
ENV PORT=3000

# Expose port
EXPOSE 3000

CMD ["bun", "run", "src/app.ts"]
