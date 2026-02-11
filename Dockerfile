FROM golang:1.21-alpine AS builder

# Instalar dependências
RUN apk add --no-cache git

# Configurar diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependência
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o geolocation-api .

# Etapa final - imagem mínima
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binário compilado
COPY --from=builder /app/geolocation-api .

# Expor porta
EXPOSE 8080

# Comando padrão
CMD ["./geolocation-api", "-serve"]
