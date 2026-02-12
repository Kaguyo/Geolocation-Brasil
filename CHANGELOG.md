# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

## [v1.1.0] - 2026-02-12

### ✨ Novo
- **Importação de 234.691 registros brasileiros**: Anteriormente importava ~5.570, agora importa todos os municípios, bairros e localidades do Brasil (99.6% do arquivo GeoNames)
- **Validação rigorosa de dados**: Filtra automaticamente registros não-brasileiros e com estados inválidos
- **Busca sem estado obrigatório**: Agora é possível buscar apenas por nome de município. Retorna o resultado mais populoso quando há duplicados
- **Ordenação por população**: Quando há múltiplas cidades com o mesmo nome, retorna a mais populosa

### 🔧 Melhorado
- Melhoria na documentação de endpoints (clarificando quando estado é obrigatório ou opcional)
- Exemplos de URL encoding adicionados à documentação
- Logs mais detalhados durante importação com estatísticas de rejeição
- Variação de códigos de admin1 agora mapeados corretamente (01-31)

### 🐛 Corrigido
- Problema onde cidades com nomes duplicados em diferentes estados retornavam resultado incorreto
- Filtro de estado agora é opcional (antes era sempre obrigatório)
- URL encoding documentado para espaços e caracteres especiais

### 📚 Documentação
- [README.md](README.md) - Atualizado com novas características
- [example.md](example.md) - Quick start simplificado
- [EXAMPLES.md](EXAMPLES.md) - Exemplos em múltiplas linguagens com URL encoding
- [CHANGELOG.md](CHANGELOG.md) - Este arquivo

## Estatísticas de Importação

**Dados Importados:**
- Total de linhas processadas: 235.522
- Total aceitas e importadas: **234.691**
- Rejeitadas por país diferente de BR: 0
- Rejeitadas por estado inválido: 828
- Rejeitadas por coordenadas fora dos bounds: 3
- **Taxa de sucesso: 99,6%**

**Tempo de importação:** ~5-10 minutos

## Mapa de Códigos Admin1 (GeoNames → Estado)

| Código | Estado | Código | Estado |
|--------|--------|--------|--------|
| 01 | DF | 17 | RO |
| 02 | ES | 18 | RR |
| 03 | BA | 19 | SC |
| 04 | GO | 20 | SP |
| 05 | MA | 21 | SE |
| 06 | MT | 22 | TO |
| 07 | MS | 23 | RS |
| 08 | MG | 24 | RO |
| 09 | PA | 25 | AC |
| 10 | PB | 26 | SC |
| 11 | PR | 27 | SP |
| 12 | PE | 28 | AL |
| 13 | PI | 29 | AP |
| 14 | RJ | 30 | AM |
| 15 | RN | 31 | CE |
| 16 | RS |  |  |

## Como Atualizar

Se você está usando uma versão anterior:

```bash
# Fazer pull das mudanças
git pull origin main

# Limpar dados antigos
rm BR.zip BR.txt 2>/dev/null

# Reimportar dados (recomendado)
go run ./cmd -importall -serve
```

## Notas Importantes

- ⚠️ A busca agora é **case-sensitive** para nomes de cidades
- 📍 Sempre use URL encoding para espaços: `São%20Paulo`
- 🔍 Para buscas exatas, sempre especifique o estado: `?estado=SP`
- 🌍 Bounds de validação geográfica: lat [-33.8, 5.4], lon [-74.0, -28.7]
