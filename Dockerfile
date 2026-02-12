# ==============================
# Etapa 1 - Build
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copiar apenas dependências primeiro (melhora cache)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código
COPY . .

# Compilar apontando para ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o geolocation-api ./cmd

# ==============================
# Etapa 2 - Runtime mínima
# ==============================
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/geolocation-api .

EXPOSE 8080

# Default: baixar, descompactar e importar dados, depois iniciar servidor
CMD ["./geolocation-api", "-importall", "-serve", "-port=8080"]

