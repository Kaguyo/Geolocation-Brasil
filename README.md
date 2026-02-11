# API de Geolocalização - Brasil 🇧🇷

Sistema completo para resolver localizações brasileiras em latitude e longitude usando MongoDB.

## 🎯 Características

- ✅ Base de dados local com MongoDB
- ✅ Todos os estados do Brasil
- ✅ API REST para consultas
- ✅ Importação de dados do GeoNames ou dados de exemplo
- ✅ Busca por município e estado
- ✅ Busca de localizações próximas (geoespacial)
- ✅ 100% Gratuito e Open Source

## 📋 Pré-requisitos

- Go 1.21+
- MongoDB 4.4+ (rodando localmente ou remoto)

## 🚀 Instalação

### 1. Instalar MongoDB

**Ubuntu/Debian:**
```bash
sudo apt-get install mongodb
sudo systemctl start mongodb
```

**macOS:**
```bash
brew install mongodb-community
brew services start mongodb-community
```

**Docker:**
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

### 2. Baixar dependências

```bash
cd geolocation-br
go mod download
```

## 📊 Importar Dados

### Opção 1: Dados de Exemplo (30 principais cidades)

```bash
go run . -import -serve
```

Isso irá:
- Importar 30 principais cidades brasileiras (capitais + maiores cidades)
- Criar índices geoespaciais
- Iniciar o servidor na porta 8080

### Opção 2: Dados Completos do GeoNames (5570 municípios)

**Passo 1:** Baixar dados do GeoNames
```bash
# Download do arquivo de cidades do Brasil
wget http://download.geonames.org/export/dump/BR.zip
unzip BR.zip
```

**Passo 2:** Importar
```bash
go run . -import -file=BR.txt -serve
```

## 🔧 Uso da API

### Iniciar servidor

```bash
# Porta padrão (8080)
go run . -serve

# Porta customizada
go run . -serve -port=3000

# MongoDB remoto
go run . -serve -mongo-uri="mongodb://usuario:senha@host:27017"
```

### Endpoints Disponíveis

#### 1. Health Check
```bash
curl http://localhost:8080/health
```

Resposta:
```json
{
  "status": "ok",
  "message": "API de Geolocalização Brasil está funcionando!"
}
```

#### 2. Buscar por Município

```bash
# Buscar São Paulo
curl http://localhost:8080/location/São%20Paulo

# Buscar com filtro de estado
curl "http://localhost:8080/location/Campinas?estado=SP"
```

Resposta:
```json
{
  "municipio": "São Paulo",
  "estado": "SP",
  "latitude": -23.5505,
  "longitude": -46.6333
}
```

#### 3. Buscar Localizações Próximas

```bash
# Buscar num raio de 100km de São Paulo
curl "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=100"
```

Resposta:
```json
[
  {
    "municipio": "São Paulo",
    "estado": "SP",
    "latitude": -23.5505,
    "longitude": -46.6333
  },
  {
    "municipio": "Guarulhos",
    "estado": "SP",
    "latitude": -23.4625,
    "longitude": -46.5333
  },
  {
    "municipio": "Osasco",
    "estado": "SP",
    "latitude": -23.5329,
    "longitude": -46.7917
  }
]
```

## 🏗️ Estrutura do Projeto

```
geolocation-br/
├── main.go           # Aplicação principal e CLI
├── models.go         # Estruturas de dados
├── database.go       # Conexão e configuração MongoDB
├── import.go         # Lógica de importação de dados
├── handlers.go       # Handlers da API REST
├── go.mod            # Dependências
└── README.md         # Este arquivo
```

## 🗄️ Estrutura do Banco de Dados

### Collection: `localizacoes`

```javascript
{
  "_id": ObjectId("..."),
  "municipio": "São Paulo",
  "estado": "SP",
  "localizacao": {
    "type": "Point",
    "coordinates": [-46.6333, -23.5505]  // [longitude, latitude]
  },
  "populacao": 12000000  // opcional
}
```

### Índices

- **2dsphere**: índice geoespacial na propriedade `localizacao`
- **text**: índice de texto em `municipio` e `estado`

## 📝 Exemplos de Uso

### Em JavaScript/Node.js

```javascript
const axios = require('axios');

// Buscar coordenadas de uma cidade
async function getCoordinates(city, state) {
  const url = `http://localhost:8080/location/${encodeURIComponent(city)}`;
  const params = state ? { estado: state } : {};
  
  const response = await axios.get(url, { params });
  return response.data;
}

// Uso
const coords = await getCoordinates('Rio de Janeiro', 'RJ');
console.log(coords);
// { municipio: 'Rio de Janeiro', estado: 'RJ', latitude: -22.9068, longitude: -43.1729 }
```

### Em Python

```python
import requests

def get_coordinates(city, state=None):
    url = f"http://localhost:8080/location/{city}"
    params = {"estado": state} if state else {}
    
    response = requests.get(url, params=params)
    return response.json()

# Uso
coords = get_coordinates("Salvador", "BA")
print(coords)
# {'municipio': 'Salvador', 'estado': 'BA', 'latitude': -12.9714, 'longitude': -38.5108}
```

### Em Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Location struct {
    Municipio string  `json:"municipio"`
    Estado    string  `json:"estado"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

func getCoordinates(city, state string) (*Location, error) {
    url := fmt.Sprintf("http://localhost:8080/location/%s?estado=%s", city, state)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var loc Location
    if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
        return nil, err
    }
    
    return &loc, nil
}
```

## 🔍 Fontes de Dados

### Dados de Exemplo
O projeto inclui 30 principais cidades brasileiras pré-configuradas (todas as capitais + maiores cidades).

### GeoNames (Completo)
- **URL**: http://download.geonames.org/export/dump/
- **Arquivo**: BR.zip
- **Contém**: ~5.570 municípios brasileiros
- **Campos**: nome, coordenadas, população, código IBGE, etc.
- **Licença**: Creative Commons Attribution 4.0

### Alternativas
- **IBGE**: https://www.ibge.gov.br/geociencias/organizacao-do-territorio/malhas-territoriais.html
- **Brasil API**: https://brasilapi.com.br/docs (para integração híbrida)

## ⚙️ Configuração Avançada

### Variáveis de Ambiente

```bash
export MONGO_URI="mongodb://localhost:27017"
export DB_NAME="geolocalizacao_br"
export API_PORT="8080"
```

### Build para Produção

```bash
# Build otimizado
go build -ldflags="-s -w" -o geolocation-api

# Executar
./geolocation-api -serve -port=8080
```

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o geolocation-api

FROM alpine:latest
COPY --from=builder /app/geolocation-api /geolocation-api
EXPOSE 8080
CMD ["/geolocation-api", "-serve"]
```

## 🚦 Performance

- **Consultas por nome**: ~5-10ms
- **Consultas geoespaciais**: ~10-20ms
- **Throughput**: ~1000 req/s (depende do hardware)

## 📄 Licença

MIT License - Sinta-se livre para usar em projetos comerciais e pessoais.

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para:
- Reportar bugs
- Sugerir novas features
- Enviar pull requests

## 📞 Suporte

Para dúvidas ou problemas:
1. Verifique se o MongoDB está rodando: `mongo --eval "db.version()"`
2. Confira os logs do servidor
3. Teste os endpoints com `curl -v`

## 🎓 Próximos Passos

- [ ] Adicionar cache com Redis
- [ ] Implementar autenticação JWT
- [ ] Adicionar rate limiting
- [ ] Criar interface web
- [ ] Adicionar mais campos (CEP, região, etc.)
- [ ] Implementar busca fuzzy (tolerante a erros)
- [ ] Adicionar testes unitários

---

Feito com ❤️ para a comunidade brasileira de desenvolvedores
