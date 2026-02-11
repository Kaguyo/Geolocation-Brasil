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

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Erro ao ler linha: %v", err)
			continue
		}

		// Formato GeoNames: geonameid, name, asciiname, alternatenames, latitude, longitude, ...
		if len(record) < 18 {
			continue
		}

		lat, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			continue
		}

		lon, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			continue
		}

		population := 0
		if record[14] != "" {
			population, _ = strconv.Atoi(record[14])
		}

		// record[10] é o código admin1 (estado)
		location := domain.Location{
			Municipio: record[1],  // name
			Estado:    record[10], // admin1 code
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
				log.Println(fmt.Errorf(""))
				return err
			}

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
		{Municipio: "Porto Alegre", Estado: "RS", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-51.2302, -30.0346}}},
		{Municipio: "Belém", Estado: "PA", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-48.5044, -1.4558}}},
		{Municipio: "Guarulhos", Estado: "SP", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-46.5333, -23.4625}}},
		{Municipio: "Campinas", Estado: "SP", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-47.0608, -22.9099}}},
		{Municipio: "São Luís", Estado: "MA", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-44.3028, -2.5387}}},
		{Municipio: "São Gonçalo", Estado: "RJ", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-43.0539, -22.8268}}},
		{Municipio: "Maceió", Estado: "AL", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-35.7353, -9.6658}}},
		{Municipio: "Duque de Caxias", Estado: "RJ", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-43.3055, -22.7858}}},
		{Municipio: "Natal", Estado: "RN", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-35.2094, -5.7945}}},
		{Municipio: "Teresina", Estado: "PI", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-42.8034, -5.0892}}},
		{Municipio: "Campo Grande", Estado: "MS", Localizacao: domain.GeoJSON{Type: "Point", Coordinates: [2]float64{-54.6295, -20.4697}}},
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
	loc, err := is.repo.GetLocationsInKilometerRange(ctx, longitude, latitude, rangeInKilometers)
	if err != nil {
		return nil, err
	}

	return loc, nil
}
