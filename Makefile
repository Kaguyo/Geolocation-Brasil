.PHONY: help install run import serve build clean test docker

help: ## Mostrar ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Instalar dependências
	go mod download
	go mod tidy

run: ## Importar dados de exemplo e iniciar servidor
	go run . -import -serve

import: ## Importar apenas dados de exemplo
	go run . -import

import-geonames: ## Importar dados completos do GeoNames (requer BR.txt)
	@if [ ! -f BR.txt ]; then \
		echo "Baixando dados do GeoNames..."; \
		wget http://download.geonames.org/export/dump/BR.zip; \
		unzip BR.zip; \
	fi
	go run . -import -file=BR.txt

serve: ## Iniciar servidor (porta 8080)
	go run . -serve

serve-3000: ## Iniciar servidor na porta 3000
	go run . -serve -port=3000

build: ## Compilar binário otimizado
	go build -ldflags="-s -w" -o bin/geolocation-api .

build-linux: ## Compilar para Linux
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/geolocation-api-linux .

build-mac: ## Compilar para macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/geolocation-api-mac .

build-windows: ## Compilar para Windows
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/geolocation-api.exe .

clean: ## Limpar arquivos compilados
	rm -rf bin/
	rm -f BR.txt BR.zip

test: ## Executar testes
	go test -v ./...

docker-build: ## Construir imagem Docker
	docker build -t geolocation-api .

docker-run: ## Executar container Docker
	docker-compose up -d

docker-stop: ## Parar container Docker
	docker-compose down

test-api: ## Testar endpoints da API (requer servidor rodando)
	@echo "Testando /health..."
	@curl -s http://localhost:8080/health | jq .
	@echo "\nTestando /location/São Paulo..."
	@curl -s http://localhost:8080/location/São%20Paulo | jq .
	@echo "\nTestando /nearby..."
	@curl -s "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=100" | jq .

mongodb-start: ## Iniciar MongoDB com Docker
	docker run -d -p 27017:27017 --name mongodb-geo mongo:latest

mongodb-stop: ## Parar MongoDB Docker
	docker stop mongodb-geo
	docker rm mongodb-geo

dev: ## Modo desenvolvimento com hot reload (requer air)
	air
