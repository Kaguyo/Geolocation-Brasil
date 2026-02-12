# Exemplos de Uso

## Setup Inicial - Importar Dados

### Opção 1: Importar Dados Completos do GeoNames (RECOMENDADO)

```bash
# Baixa BR.zip, descompacta e importa 234.691 registros brasileiros, depois inicia servidor
go run ./cmd -importall -serve

# Apenas importar dados sem iniciar servidor
go run ./cmd -importall
```

**Inclui:**
- ✅ 234.691 registros (municípios, bairros, localidades)
- ✅ Filtro automático de dados não-brasileiros
- ✅ Validação de estados
- ✅ Índices geoespaciais e de texto
- ✅ Tempo de importação: ~5-10 minutos

### Opção 2: Importar Dados de Exemplo (30 cidades principais)

```bash
# Importar dados de exemplo e iniciar servidor
go run ./cmd -import -serve

# Apenas importar dados de exemplo
go run ./cmd -import
```

### Opção 3: Importar de Arquivo Customizado

```bash
# Importar dados de arquivo no formato GeoNames
go run ./cmd -import -file=meu_arquivo.txt -serve
```

## Executar Servidor (Após Importação)

```bash
# Servidor padrão (porta 8080)
go run ./cmd -serve

# Servidor em porta customizada
go run ./cmd -serve -port=3000

# Servidor com MongoDB remoto
go run ./cmd -serve -mongo-uri="mongodb://user:pass@host:27017"
```

---

# Exemplos de Requisições HTTP
# Use com ferramentas como curl, httpie, ou Postman

## 1. Health Check
### curl
curl http://localhost:8080/health

### httpie
http GET http://localhost:8080/health


## 2. Buscar por Município (sem estado - retorna mais populoso)

### São Paulo
```bash
curl "http://localhost:8080/location/S%C3%A3o%20Paulo"
```
Resposta:
```json
{"municipio":"São Paulo","estado":"SP","latitude":-22,"longitude":-49}
```

### Rio de Janeiro
```bash
curl "http://localhost:8080/location/Rio%20de%20Janeiro"
```

### Porto Alegre
```bash
curl "http://localhost:8080/location/Porto%20Alegre"
```

### Brasília
```bash
curl "http://localhost:8080/location/Brasília"
```


## 3. Buscar por Município e Estado Específico (busca exata)

### São Paulo, SP
```bash
curl "http://localhost:8080/location/S%C3%A3o%20Paulo?estado=SP"
```

### Campinas, SP
```bash
curl "http://localhost:8080/location/Campinas?estado=SP"
```

### Santos, SP
```bash
curl "http://localhost:8080/location/Santos?estado=SP"
```

### Salvador, BA
```bash
curl "http://localhost:8080/location/Salvador?estado=BA"
```

### Porto Alegre, RS (específico)
```bash
curl "http://localhost:8080/location/Porto%20Alegre?estado=RS"
```


## 4. Buscar Localizações Próximas (Geoespacial)

### Próximos a São Paulo (raio de 50km)
curl "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=50"

### Próximos a São Paulo (raio de 100km)
curl "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=100"

### Próximos ao Rio de Janeiro (raio de 80km)
curl "http://localhost:8080/nearby?lat=-22.9068&lon=-43.1729&distance=80"

### Próximos a Brasília (raio de 150km)
curl "http://localhost:8080/nearby?lat=-15.7801&lon=-47.9292&distance=150"

### Com formatação
curl -s "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=100" | jq .


## 5. URL Encoding - Importante!

Nomes de cidades com espaços e caracteres especiais precisam ser codificados na URL:

| Texto Original | URL Encoded |
|---|---|
| `São Paulo` | `S%C3%A3o%20Paulo` |
| `Rio de Janeiro` | `Rio%20de%20Janeiro` |
| `Brasília` | `Brasília` ou `Bras%C3%ADlia` |
| `Santa Catarina` | `Santa%20Catarina` |
| `Belo Horizonte` | `Belo%20Horizonte` |
| `Porto Alegre` | `Porto%20Alegre` |

**Exemplos:**
```bash
# ✅ Correto
curl "http://localhost:8080/location/S%C3%A3o%20Paulo"

# ❌ Incorreto (pode não encontrar)
curl "http://localhost:8080/location/São Paulo"
```


## 6. Usando JavaScript (Fetch API)

```javascript
// Health Check
fetch('http://localhost:8080/health')
  .then(res => res.json())
  .then(data => console.log(data));

// Buscar coordenadas (sem estado - retorna mais populosa)
const cidade = 'São Paulo';
fetch(`http://localhost:8080/location/${encodeURIComponent(cidade)}`)
  .then(res => res.json())
  .then(data => console.log(data));

// Buscar com estado específico
fetch('http://localhost:8080/location/Campinas?estado=SP')
  .then(res => res.json())
  .then(data => console.log(data));

// Buscar próximos
fetch('http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=100')
  .then(res => res.json())
  .then(data => console.log(data));
```


## 7. Usando Python (requests)

```python
import requests
from urllib.parse import quote

# Health Check
response = requests.get('http://localhost:8080/health')
print(response.json())

# Buscar coordenadas (sem estado - retorna mais populosa)
cidade = 'São Paulo'
response = requests.get(f'http://localhost:8080/location/{quote(cidade)}')
print(response.json())

# Buscar com estado específico
response = requests.get('http://localhost:8080/location/Campinas', params={'estado': 'SP'})
print(response.json())

# Buscar próximos (100km de São Paulo)
response = requests.get(
    'http://localhost:8080/nearby',
    params={'lat': -23.5505, 'lon': -46.6333, 'distance': 100}
)
print(response.json())
```


## 8. Usando Go (net/http)

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
)

type Location struct {
    Municipio string  `json:"municipio"`
    Estado    string  `json:"estado"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

func main() {
    // Buscar coordenadas
    cidade := "São Paulo"
    resp, err := http.Get("http://localhost:8080/location/" + url.QueryEscape(cidade))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var location Location
    json.Unmarshal(body, &location)

    fmt.Printf("Município: %s\n", location.Municipio)
    fmt.Printf("Estado: %s\n", location.Estado)
    fmt.Printf("Latitude: %f\n", location.Latitude)
    fmt.Printf("Longitude: %f\n", location.Longitude)
}
```


## 9. Teste de Performance

### Apache Bench (100 requisições, 10 concorrentes)
```bash
ab -n 100 -c 10 http://localhost:8080/health
```

### wrk (teste de 30 segundos, 10 threads, 100 conexões)
```bash
wrk -t10 -c100 -d30s http://localhost:8080/health
```


## 10. Status e Logs da Importação

Durante a importação, você verá logs como:

```
2026/02/12 12:23:46 ✅ Conectado ao MongoDB com sucesso!
2026/02/12 12:23:46 🔄 Iniciando importação completa do GeoNames...
2026/02/12 12:23:46 📥 Baixando BR.zip do GeoNames...
2026/02/12 12:23:48 ✅ Download concluído!
2026/02/12 12:23:48 📦 Descompactando arquivo...
2026/02/12 12:23:48 ✅ Descompactação concluída!
2026/02/12 12:23:48 📂 Importando dados de BR.txt (~5570 municípios) aguarde...
2026/02/12 12:25:26 ✅ Importação concluída! Total: 234691 registros
2026/02/12 12:25:26 📊 Estatísticas de rejeição:
2026/02/12 12:25:26    - Total de linhas processadas: 235522
2026/02/12 12:25:26    - Rejeitadas por país diferente de BR: 0
2026/02/12 12:25:26    - Rejeitadas por estado inválido: 828
2026/02/12 12:25:26    - Rejeitadas por erro ao ler coordenadas: 0
2026/02/12 12:25:26    - Rejeitadas por coordenadas fora dos bounds: 3
2026/02/12 12:25:26    - ✓ Aceitas e importadas: 234691
```


## 9. Postman Collection

Importe esta collection no Postman:

```json
{
  "info": {
    "name": "Geolocalização Brasil API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "http://localhost:8080/health",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["health"]
        }
      }
    },
    {
      "name": "Buscar São Paulo",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "http://localhost:8080/location/São Paulo",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["location", "São Paulo"]
        }
      }
    },
    {
      "name": "Buscar Próximos",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=100",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["nearby"],
          "query": [
            {"key": "lat", "value": "-23.5505"},
            {"key": "lon", "value": "-46.6333"},
            {"key": "distance", "value": "100"}
          ]
        }
      }
    }
  ]
}
```
