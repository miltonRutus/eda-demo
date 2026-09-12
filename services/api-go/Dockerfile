# ==============================================================================
# Development Stage (Target: dev)
# - Incluye compilador Go, herramientas para pruebas unitarias (go test)
# - Ejecución interactiva / recarga
# ==============================================================================
FROM golang:1.23-alpine AS dev

RUN apk add --no-cache git bash curl build-base

WORKDIR /app

# Descargar dependencias para aprovechar la capa de caché
COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .

# Descarga y verificación final de módulos
RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "cmd/server/main.go"]

# ==============================================================================
# Builder Stage: Compilación binaria optimizada
# ==============================================================================
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server cmd/server/main.go

# ==============================================================================
# Production Stage (Target: prod)
# - Inmutable, mínimo y seguro
# - Usuario sin privilegios root
# ==============================================================================
FROM alpine:3.20 AS prod

RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

COPY --from=builder /app/server /app/server

USER appuser
EXPOSE 8080

CMD ["/app/server"]
