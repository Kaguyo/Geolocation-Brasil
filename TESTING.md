# Guia de Verificação e Testes

Este arquivo fornece checklist e testes para validar a aplicação está funcionando corretamente.

## ✅ Checklist de Instalação

- [ ] MongoDB está instalado e rodando
- [ ] Go 1.21+ está instalado
- [ ] Repositório foi clonado
- [ ] Dependências foram baixadas (`go mod download`)
- [ ] Dados foram importados (`-importall` ou `-import`)
- [ ] Servidor inicia sem erros

## 🧪 Testes Rápidos

### 1. Health Check
```bash
curl http://localhost:8080/health
```
✅ Espera: `{"status":"ok",...}`

### 2. Busca Simples (sem estado)
```bash
curl "http://localhost:8080/location/S%C3%A3o%20Paulo"
```
✅ Espera: `São Paulo, SP, latitude, longitude`

### 3. Busca com Estado
```bash
curl "http://localhost:8080/location/Campinas?estado=SP"
```
✅ Espera: `Campinas, SP, latitude, longitude`

### 4. Busca Geoespacial
```bash
curl "http://localhost:8080/nearby?lat=-23.5505&lon=-46.6333&distance=50"
```
✅ Espera: Lista de cidades próximas

### 5. Não Encontrado
```bash
curl "http://localhost:8080/location/XYZABC"
```
✅ Espera: `{"error":"Not Found",...}`

## 📊 Validação de Dados

### Quantidade de Registros
```bash
# Conectar ao MongoDB para validar
mongosh
use geolocalizacao_br
db.geolocations.countDocuments()
```
✅ Espera: ~234.691 registros (se importall foi usado)

### Validação de Estados
```bash
# Verificar que todos os registros têm estados válidos
db.geolocations.distinct("estado").sort()
```
✅ Espera: 27 estados brasileiros (AC, AL, AP, AM, BA, CE, DF, ES, GO, MA, MT, MS, MG, PA, PB, PE, PI, RJ, RN, RS, RO, RR, SC, SP, SE, TO)

### Amostra de Cidades Principais
```bash
# Verificar capitais
db.geolocations.find({municipio: "São Paulo"}).pretty()
db.geolocations.find({municipio: "Rio de Janeiro"}).pretty()
db.geolocations.find({municipio: "Brasília"}).pretty()
```

## 🔍 Testes de URL Encoding

| Entrada | URL | Status |
|---------|-----|--------|
| `São Paulo` | `S%C3%A3o%20Paulo` | ✅ |
| `Rio de Janeiro` | `Rio%20de%20Janeiro` | ✅ |
| `Porto Alegre` | `Porto%20Alegre` | ✅ |
| `Brasília` | `Brasília` | ✅ |
| `Manaus` | `Manaus` | ✅ |

## ⚡ Performance

### Teste de Latência
```bash
time curl "http://localhost:8080/location/S%C3%A3o%20Paulo"
```
✅ Espera: < 100ms

### Teste de Carga
```bash
# 100 requisições, 10 concorrentes
ab -n 100 -c 10 http://localhost:8080/health
```
✅ Espera: Requests per second > 1000

## 🐳 Docker

### Health Check
```bash
docker ps | grep geolocation-api
```
✅ Espera: Container rodando

### Logs
```bash
docker-compose logs api | tail -20
```
✅ Espera: "Servidor iniciado na porta 8080"

## 📋 Troubleshooting

### Problema: "Localização não encontrada"
**Solução:**
1. Verifique URL encoding (espaços = %20)
2. Confirme que dados foram importados
3. Tente com estado: `?estado=SP`

### Problema: "Address already in use"
**Solução:**
```bash
# Matar processo na porta 8080
lsof -i :8080 | grep LISTEN | awk '{print $2}' | xargs kill -9

# Ou usar porta diferente
go run ./cmd -serve -port=3000
```

### Problema: "Failed to connect to MongoDB"
**Solução:**
1. Verificar se MongoDB está rodando: `systemctl status mongodb`
2. Iniciar MongoDB: `systemctl start mongodb`
3. Usar MongoDB remoto: `-mongo-uri="mongodb://host:27017"`

### Problema: Importação muito lenta
**Solução:**
1. Normal - importação de 235k registros leva 5-10 minutos
2. Verificar conexão de internet (download de BR.zip)
3. Verificar espaço em disco (arquivo temporário ~50MB)

## ✨ Recursos Adicionais

- [README.md](README.md) - Documentação completa
- [EXAMPLES.md](EXAMPLES.md) - Exemplos em múltiplas linguagens
- [example.md](example.md) - Quick start
- [CHANGELOG.md](CHANGELOG.md) - Histórico de mudanças
