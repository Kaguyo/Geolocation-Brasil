package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	domain "github.com/Kaguyo/Geolocation-Brasil/internal/domain/entities"
	domainIF "github.com/Kaguyo/Geolocation-Brasil/internal/domain/interfaces"
)

type ImportService struct {
	repo domainIF.IGeoRepository
}

func NewGeoService(repo domainIF.IGeoRepository) *ImportService {

	return &ImportService{
		repo: repo,
	}
}

// ImportData importa dados de um arquivo CSV
func (is *ImportService) ImportData(ctx context.Context, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // GeoNames usa tab como separador
	reader.LazyQuotes = true

	// Pular header se existir
	_, err = reader.Read()
	if err != nil {
		log.Println(fmt.Errorf("erro ao ler header: %v", err))
		return err
	}

	var locations []domain.Location

	count := 0
	rejectedTotal := 0
	rejectedCountry := 0
	rejectedState := 0
	rejectedCoords := 0
	rejectedBounds := 0

	// Mapa de conversão: admin1 code do GeoNames (número) -> código de estado (2 letras)
	stateCodeMap := map[string]string{
		"01": "DF", "02": "ES", "03": "BA", "04": "GO", "05": "MA", "06": "MT", "07": "MS",
		"08": "MG", "09": "PA", "10": "PB", "11": "PR", "12": "PE", "13": "PI",
		"14": "RJ", "15": "RN", "16": "RS", "17": "RO", "18": "RR", "19": "SC",
		"20": "SP", "21": "SE", "22": "TO", "23": "RS", "24": "RO", "25": "AC", "26": "SC", "27": "SP", "28": "AL", "29": "AP", "30": "AM", "31": "CE",
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Erro ao ler linha: %v", err)
			continue
		}

		rejectedTotal++

		// Formato GeoNames: geonameid, name, asciiname, alternatenames, latitude, longitude, ...
		if len(record) < 18 {
			continue
		}

		// Filtro crítico: apenas importar registros do Brasil (countryCode == "BR")
		if record[8] != "BR" {
			rejectedCountry++
			continue
		}

		// Validar estado: não pode estar vazio
		estadoCode := record[10]
		if estadoCode == "" {
			rejectedState++
			continue
		}

		// Converter admin1 code (número) para estado (2 letras)
		estado, exists := stateCodeMap[estadoCode]
		if !exists {
			log.Printf("⚠️ Estado inválido ignorado: %s (código não encontrado no mapa, município: %s)", estadoCode, record[1])
			rejectedState++
			continue
		}

		lat, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			rejectedCoords++
			continue
		}

		lon, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			rejectedCoords++
			continue
		}

		// Validar bounds de coordenadas brasileiras (segurança adicional)
		// Brasil: lat entre -33.7 e 5.3, lon entre -73.9 e -28.8
		if lat < -33.8 || lat > 5.4 || lon < -74.0 || lon > -28.7 {
			rejectedBounds++
			log.Printf("DEBUG: Coordenadas fora do bounds: lat=%v, lon=%v (municipio=%s, estado=%s)", lat, lon, record[1], estado)
			continue
		}

		population := 0
		if record[14] != "" {
			population, _ = strconv.Atoi(record[14])
		}

		// Usar o estado convertido (2 letras)
		location := domain.Location{
			Municipio: record[1], // name
			Estado:    estado,    // código de estado convertido (SP, BA, etc)
			Localizacao: domain.GeoJSON{
				Type:        "Point",
				Coordinates: [2]float64{lon, lat},
			},
			Populacao: population,
		}

		locations = append(locations, location)
		count++

		// Inserir em lotes de 1000
		if count%1000 == 0 {
			err := is.repo.InsertLocations(ctx, locations)
			if err != nil {
				log.Println(fmt.Errorf("erro ao inserir lote: %v", err))
				return err
			}
			locations = []domain.Location{} // Limpar o array após inserção
		}
	}

	// Inserir registros restantes
	if len(locations) > 0 {
		err := is.repo.InsertLocations(ctx, locations)
		if err != nil {
			return err
		}
	}

	log.Printf("✅ Importação concluída! Total: %d registros", count)
	log.Printf("📊 Estatísticas de rejeição:")
	log.Printf("   - Total de linhas processadas: %d", rejectedTotal)
	log.Printf("   - Rejeitadas por país diferente de BR: %d", rejectedCountry)
	log.Printf("   - Rejeitadas por estado inválido: %d", rejectedState)
	log.Printf("   - Rejeitadas por erro ao ler coordenadas: %d", rejectedCoords)
	log.Printf("   - Rejeitadas por coordenadas fora dos bounds: %d", rejectedBounds)
	log.Printf("   - ✓ Aceitas e importadas: %d", count)
	return nil
}

// ImportBrazilianCities importa dados simplificados de cidades brasileiras
func (is *ImportService) ImportBrazilianCitiesExampleTest(ctx context.Context) error {
	// Dados de exemplo das capitais brasileiras
	cities := []domain.Location{
		{Municipio: "São Paulo", Estado: "SP", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-46.6333, -23.5505}}},
		{Municipio: "Rio de Janeiro", Estado: "RJ", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-43.1729, -22.9068}}},
		{Municipio: "Brasília", Estado: "DF", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-47.9292, -15.7801}}},
		{Municipio: "Salvador", Estado: "BA", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-38.5108, -12.9714}}},
		{Municipio: "Fortaleza", Estado: "CE", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-38.5434, -3.7172}}},
		{Municipio: "Belo Horizonte", Estado: "MG", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-43.9378, -19.9208}}},
		{Municipio: "Manaus", Estado: "AM", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-60.0217, -3.1190}}},
		{Municipio: "Curitiba", Estado: "PR", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-49.2643, -25.4284}}},
		{Municipio: "Recife", Estado: "PE", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-34.8813, -8.0476}}},
		{Municipio: "Goiânia", Estado: "GO", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-49.2532, -16.6864}}},
		{Municipio: "Porto Aleise", Estado: "RS", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-51.2302, -30.0346}}},
		{Municipio: "Belém", Estado: "PA", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-48.5044, -1.4558}}},
		{Municipio: "Guarulhos", Estado: "SP", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-46.5333, -23.4625}}},
		{Municipio: "Campinas", Estado: "SP", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-47.0608, -22.9099}}},
		{Municipio: "São Luís", Estado: "MA", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-44.3028, -2.5387}}},
		{Municipio: "São Gonçalo", Estado: "RJ", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-43.0539, -22.8268}}},
		{Municipio: "Maceió", Estado: "AL", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-35.7353, -9.6658}}},
		{Municipio: "Duque de Caxias", Estado: "RJ", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-43.3055, -22.7858}}},
		{Municipio: "Natal", Estado: "RN", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-35.2094, -5.7945}}},
		{Municipio: "Teresina", Estado: "PI", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-42.8034, -5.0892}}},
		{Municipio: "Campo isande", Estado: "MS", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-54.6295, -20.4697}}},
		{Municipio: "João Pessoa", Estado: "PB", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-34.8631, -7.1195}}},
		{Municipio: "Jaboatão dos Guararapes", Estado: "PE", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-35.0147, -8.1130}}},
		{Municipio: "Osasco", Estado: "SP", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-46.7917, -23.5329}}},
		{Municipio: "Santo André", Estado: "SP", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-46.5386, -23.6639}}},
		{Municipio: "São Bernardo do Campo", Estado: "SP", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-46.5650, -23.6914}}},
		{Municipio: "Ribeirão Preto", Estado: "SP", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-47.8103, -21.1704}}},
		{Municipio: "Uberlândia", Estado: "MG", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-48.2772, -18.9186}}},
		{Municipio: "Contagem", Estado: "MG", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-44.0539, -19.9320}}},
		{Municipio: "Aracaju", Estado: "SE", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-37.0731, -10.9091}}},
	}

	err := is.repo.ImportTest(ctx, cities)
	if err != nil {
		return err
	}

	return nil
}

// GetLocationByName busca localização por nome de município ou estado
func (is *ImportService) GetLocationByName(ctx context.Context, municipio, estado string) (*domain.Location, error) {
	loc, err := is.repo.GetLocationByName(ctx, municipio, estado)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

// GetNearbyLocations busca localizações próximas a um ponto
func (is *ImportService) GetLocationsInKilometersRange(ctx context.Context, longitude, latitude float64, rangeInKilometers float64) (*[]domain.Location, error) {
	loc, err := is.repo.GetLocationsInKilometersRange(ctx, longitude, latitude, rangeInKilometers)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

// CreateGeoIndex cria índice geoespacial
func (is *ImportService) CreateGeoIndex(ctx context.Context) error {
	err := is.repo.CreateGeoIndex(ctx)
	if err != nil {
		return err
	}
	return nil
}

// CreateTextIndex cria índice de texto para busca
func (is *ImportService) CreateTextIndex(ctx context.Context) error {
	err := is.repo.CreateTextIndex(ctx)
	if err != nil {
		return err
	}
	return nil
}

// ResetCollection recria a coleção, removendo todos os dados existentes
func (is *ImportService) ResetCollection(ctx context.Context, collection string) error {
	err := is.repo.DropCollection(ctx, collection)
	if err != nil {
		return err
	}

	err = is.repo.CreateGeoIndex(ctx)
	if err != nil {
		return err
	}

	err = is.repo.CreateTextIndex(ctx)
	if err != nil {
		return err
	}

	return nil
}
