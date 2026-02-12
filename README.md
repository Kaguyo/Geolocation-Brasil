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

### Opção 1: Dados Completos do GeoNames (234.691 registros) - RECOMENDADO

O novo comando `-importall` faz tudo automaticamente: baixa o arquivo do GeoNames, descompacta e importa para o banco de dados.

**Sem servir a API:**
```bash
go run ./cmd -importall
```

**Importar dados e iniciar servidor:**
```bash
go run ./cmd -importall -serve
```

Isso irá:
- ✅ Baixar BR.zip automaticamente do servidor GeoNames
- ✅ Descompactar o arquivo
- ✅ Importar **234.691 registros brasileiros** (municípios, bairros e localidades)
- ✅ Filtrar automaticamente apenas dados do Brasil com estados válidos
- ✅ Criar índices geoespaciais e de texto
- ✅ Iniciar o servidor na porta 8080 (se `-serve` foi usado)

**Tempo de importação:** ~5-10 minutos (depende da conexão)

### Opção 2: Dados de Exemplo (30 principais cidades)

Para testes rápidos, use:
```bash
go run ./cmd -import -serve
```

Isso irá:
- Importar 30 principais cidades brasileiras (capitais + maiores cidades)
- Criar índices geoespaciais
- Iniciar o servidor na porta 8080

### Opção 3: Importar de Arquivo CSV Customizado

Se você tem um arquivo CSV no formato GeoNames:
```bash
go run ./cmd -import -file=seu_arquivo.txt -serve
```

**Nota:** O arquivo deve estar no formato GeoNames com campos separados por tab.

## 🔧 Uso da API

### Primeira Vez: Importar Dados + Iniciar Servidor

```bash
# Opção 1: Usar dados completos do GeoNames (recomendado)
go run ./cmd -importall -serve

# Opção 2: Usar dados de exemplo (mais rápido para testes)
go run ./cmd -import -serve
```

### Após Importação: Apenas Iniciar Servidor

Depois que os dados foram importados uma vez, você pode apenas iniciar o servidor:

```bash
# Porta padrão (8080)
go run ./cmd -serve

# Porta customizada
go run ./cmd -serve -port=3000

# MongoDB remoto
go run ./cmd -serve -mongo-uri="mongodb://usuario:senha@host:27017"
```

**Nota:** O servidor cria os índices automaticamente na primeira execução, então você não precisa fazer nada.

### Flags Disponíveis

```
-import              Importar dados de exemplo (30 principais cidades)
-importall          Baixar BR.zip do GeoNames, descompactar e importar todos os dados (~5570 municípios)
-file string        Arquivo CSV para importar (formato GeoNames) - usado com -import
-serve              Iniciar servidor API
-port string        Porta do servidor (padrão: 8080)
-mongo-uri string   URI de conexão do MongoDB (padrão: mongodb://localhost:27017)
```

**Exemplos de uso:**
```bash
# Apenas importar dados de exemplo
go run ./cmd -import

# Importar dados de exemplo e iniciar servidor
go run ./cmd -import -serve

# Importar dados completos do GeoNames e iniciar servidor
go run ./cmd -importall -serve

# Apenas iniciar servidor (dados já importados)
go run ./cmd -serve

# Importar arquivo customizado en iniciar servidor
go run ./cmd -import -file=dados.txt -serve

# Servidor em porta customizada com MongoDB remoto
go run ./cmd -serve -port=3000 -mongo-uri="mongodb://user:pass@host:27017"
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

**Sem especificar estado (retorna mais populoso):**
```bash
# Buscar São Paulo (retorna a mais populosa)
curl "http://localhost:8080/location/S%C3%A3o%20Paulo"

# Buscar Porto Alegre (retorna a mais populosa)
curl "http://localhost:8080/location/Porto%20Alegre"
```

**Com filtro de estado (busca específica):**
```bash
# Buscar São Paulo em SP
curl "http://localhost:8080/location/S%C3%A3o%20Paulo?estado=SP"

# Buscar Campinas em SP
curl "http://localhost:8080/location/Campinas?estado=SP"

# Buscar Porto Alegre em RS
curl "http://localhost:8080/location/Porto%20Alegre?estado=RS"
```

**Resposta (sucesso):**
```json
{
  "municipio": "São Paulo",
  "estado": "SP",
  "latitude": -23.5505,
  "longitude": -46.6333
}
```

**Resposta (não encontrado):**
```json
{
  "error": "Not Found",
  "message": "Localização não encontrada"
}
```

**⚠️ Importante: URL Encoding**
- Espaços devem ser codificados como `%20`
- Caracteres especiais (ç, ã, é, ó) são codificados como UTF-8
- Exemplos:
  - `São Paulo` → `S%C3%A3o%20Paulo`
  - `Rio de Janeiro` → `Rio%20de%20Janeiro`
  - `Brasília` → `Brasília` (ou `Bras%C3%ADlia`)

#### 3. Buscar Localizações Próximas (Geoespacial)

```bash
# Buscar num raio de 50km de São Paulo (padrão)
curl "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333"

# Buscar num raio de 100km de São Paulo
curl "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=100"

# Buscar num raio de 200km de São Paulo
curl "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=200"
```

**Parâmetros:**
- `lat` (obrigatório): Latitude em graus decimais
- `lon` (obrigatório): Longitude em graus decimais
- `distance` (opcional): Distância em quilômetros (padrão: 50km)

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
├── cmd/
│   └── main.go                    # Aplicação principal e CLI
├── internal/
│   ├── api/
│   │   ├── handlers.go            # Handlers da API REST
│   │   └── response.go            # Estruturas de resposta
│   ├── application/
│   │   └── services/
│   │       ├── import.go          # Lógica de importação de dados
│   │       └── interfaces/
│   │           └── import.go      # Interfaces dos serviços
│   ├── bootstrap/
│   │   └── container.go           # Configuração da aplicação
│   ├── domain/
│   │   ├── entities/
│   │   │   └── models.go          # Modelos de dados (Location, GeoJSON)
│   │   └── interfaces/
│   │       └── geo_repository.go  # Interface do repositório
│   ├── infrastructure/
│   │   └── mongodb/
│   │       ├── connection.go      # Conexão com MongoDB
│   │       └── geo_repository.go  # Implementação do repositório
│   └── utils/
│       └── zip.go                 # Utilitários (download, unzip)
├── go.mod                          # Dependências
├── Dockerfile                       # Container Docker
├── docker-compose.yml              # Orquestração Docker
└── README.md                        # Este arquivo
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

### Dados de Exemplo Inclusos
- **O que é**: 30 principais cidades brasileiras (todas as capitais + maiores cidades)
- **Como usar**: `go run ./cmd -import`
- **Tempo**: Instantâneo
- **Uso**: Testes e prototipagem rápida

### GeoNames - Dados Completos (RECOMENDADO)
- **URL**: http://download.geonames.org/export/dump/
- **Arquivo**: BR.zip
- **Contém**: ~5.570 municípios brasileiros
- **Campos**: nome, coordenadas, população, código IBGE, etc.
- **Licença**: Creative Commons Attribution 4.0
- **Como usar**: `go run ./cmd -importall -serve`
- **Tempo**: ~5-10 minutos na primeira vez (depende da conexão)
- **Nota**: Após a primeira importação, você só precisa usar `-serve`

### Alternativas de Dados
- **IBGE**: https://www.ibge.gov.br/geociencias/organizacao-do-territorio/malhas-territoriais.html
- **Brasil API**: https://brasilapi.com.br/docs (para integração híbrida)
- **Sua própria fonte**: Use o flag `-import -file=seu_arquivo.txt` com dados no formato GeoNames

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

**Com Docker Compose (recomendado):**
```bash
# Iniciar MongoDB + API (com importação automática de dados)
docker-compose up --build

# Em modo detached (background)
docker-compose up -d --build

# Parar serviços
docker-compose down
```

Este comando:
- ✅ Cria container MongoDB na porta 27018
- ✅ Cria container API na porta 8080
- ✅ Baixa e importa automaticamente ~5.570 municípios do GeoNames
- ✅ Cria índices geoespaciais
- ✅ Inicia o servidor API

**Com Docker diretamente:**
```bash
# Build
docker build -t geolocation-api .

# Rodar com MongoDB local
docker run -d \
  --name geolocation-api \
  -p 8080:8080 \
  -e MONGO_URI="mongodb://host.docker.internal:27017" \
  geolocation-api
```

**Customizar comportamento do Docker:**

Para usar apenas dados de exemplo ao invés de baixar todos os dados:
```dockerfile
# Editar Dockerfile e alterar CMD para:
CMD ["./geolocation-api", "-import", "-serve", "-port=8080"]
```

Ou via docker-compose:
```yaml
command:
  [
    "./geolocation-api",
    "-import",           # Dados de exemplo apenas
    "-serve",
    "-mongo-uri=mongodb://mongodb:27017"
  ]
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
