# Exemplos de Requisições HTTP
# Use com ferramentas como curl, httpie, ou Postman

## 1. Health Check
### curl
curl http://localhost:8080/health

### httpie
http GET http://localhost:8080/health


## 2. Buscar por Município

### São Paulo
curl http://localhost:8080/location/São%20Paulo

### Rio de Janeiro
curl http://localhost:8080/location/Rio%20de%20Janeiro

### Brasília
curl http://localhost:8080/location/Brasília

### Com formatação JSON (usando jq)
curl -s http://localhost:8080/location/São%20Paulo | jq .


## 3. Buscar por Município com Estado

### Campinas, SP
curl "http://localhost:8080/location/Campinas?estado=SP"

### Santos, SP
curl "http://localhost:8080/location/Santos?estado=SP"

### Salvador, BA
curl "http://localhost:8080/location/Salvador?estado=BA"


## 4. Buscar Localizações Próximas

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


## 5. Usando JavaScript (Fetch API)

```javascript
// Health Check
fetch('http://localhost:8080/health')
  .then(res => res.json())
  .then(data => console.log(data));

// Buscar coordenadas
fetch('http://localhost:8080/location/São Paulo')
  .then(res => res.json())
  .then(data => console.log(data));

// Buscar próximos
fetch('http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=100')
  .then(res => res.json())
  .then(data => console.log(data));
```


## 6. Usando Python (requests)

```python
import requests

# Health Check
response = requests.get('http://localhost:8080/health')
print(response.json())

# Buscar coordenadas
response = requests.get('http://localhost:8080/location/São Paulo')
print(response.json())

# Buscar próximos
response = requests.get(
    'http://localhost:8080/nearby',
    params={'lat': -23.5505, 'lon': -46.6333, 'distance': 100}
)
print(response.json())
```


## 7. Usando Go (net/http)

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

func main() {
    // Buscar coordenadas
    resp, err := http.Get("http://localhost:8080/location/São Paulo")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    
    var result map[string]interface{}
    json.Unmarshal(body, &result)
    
    fmt.Printf("Município: %s\n", result["municipio"])
    fmt.Printf("Latitude: %v\n", result["latitude"])
    fmt.Printf("Longitude: %v\n", result["longitude"])
}
```


## 8. Teste de Performance

### Apache Bench (100 requisições, 10 concorrentes)
ab -n 100 -c 10 http://localhost:8080/health

### wrk (teste de 30 segundos, 10 threads, 100 conexões)
wrk -t10 -c100 -d30s http://localhost:8080/health


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
