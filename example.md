# Quick Start Guide

## 🚀 Começar em 3 passos

### 1. Importar Dados (primeira vez só)

```bash
# Opção A: Dados completos do GeoNames (~234.691 registros)
go run ./cmd -importall

# Opção B: Dados de exemplo (30 principais cidades)
go run ./cmd -import
```

### 2. Iniciar Servidor

```bash
go run ./cmd -serve
```

API estará disponível em: `http://localhost:8080`

### 3. Testar

```bash
# Health check
curl http://localhost:8080/health

# Buscar São Paulo (sem especificar estado)
curl "http://localhost:8080/location/S%C3%A3o%20Paulo"

# Buscar São Paulo em SP especificamente
curl "http://localhost:8080/location/S%C3%A3o%20Paulo?estado=SP"

# Buscar cidades próximas (50km de São Paulo)
curl "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333"

# Buscar cidades próximas (100km de São Paulo)
curl "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=100"
```

---

## 🐳 Com Docker

```bash
# Tudo em um comando
docker-compose up --build

# API estará em http://localhost:8080
```

---

## 📚 Documentação Completa

- **[README.md](README.md)** - Documentação completa do projeto
- **[EXAMPLES.md](EXAMPLES.md)** - Exemplos de requisições HTTP em várias linguagens

---

## 🪪 Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/health` | Verificar status da API |
| GET | `/location/{municipio}` | Buscar coordenadas de um município (retorna mais populoso) |
| GET | `/location/{municipio}?estado=XX` | Buscar coordenadas de um município em um estado específico |
| GET | `/nearby?lat=X&lon=Y&distance=Z` | Buscar cidades próximas (Z em km, padrão 50km) |

---

## 💡 Dicas

- **Primeira execução com `-importall`**: Pode levar 5-10 minutos (depende da conexão)
- **Próximas execuções**: Use apenas `-serve` (dados já estão importados)
- **MongoDB local**: Certifique-se que MongoDB está rodando na porta 27017
- **MongoDB remoto**: Use `-mongo-uri="mongodb://host:porta"`
- **URL Encoding**: Nomes com espaços precisam de `%20` (ex: `S%C3%A3o%20Paulo`)
- **Sem estado**: Retorna o resultado mais populoso quando há nomes duplicados

---

## ❓ Problemas?

1. MongoDB não está rodando?
   ```bash
   # Linux/WSL
   sudo systemctl start mongodb

   # macOS
   brew services start mongodb-community
   ```

2. Porta 8080 já está em uso?
   ```bash
   go run ./cmd -serve -port=3000
   ```

3. Erros de download (sem internet)?
   ```bash
   # Use dados de exemplo
   go run ./cmd -import -serve
   ```

4. Nenhum resultado encontrado?
   - Verifique se o nome está escrito corretamente
   - Use URL encoding para espaços: `%20`
   - Tente com estado específico: `?estado=SP`
   - Certifique-se de que os dados foram importados
