package interfaces

import (
	"context"

	domain "github.com/Kaguyo/Geolocation-Brasil/internal/domain/entities"
)

type IGeoRepository interface {
	// CreateGeoIndex cria índice geoespacial
	CreateGeoIndex(ctx context.Context) error
	// CreateTextIndex cria índice de texto para busca
	CreateTextIndex(ctx context.Context) error
	// Inserts as many locations as given through parameter
	InsertLocations(ctx context.Context, locationBuffer []domain.Location) error
	// GetNearbyLocations busca localizações próximas a um ponto
	GetLocationsInKilometersRange(ctx context.Context, longitude, latitude float64, maxDistanceKm float64) (*[]domain.Location, error)
	// GetLocationByName busca localização por nome de município
	GetLocationByName(ctx context.Context, municipio, estado string) (*domain.Location, error)
	// ImportBrazilianCities importa dados simplificados de cidades brasileiras
	ImportTest(ctx context.Context, locations []domain.Location) error
	// DropCollection recria a coleção, removendo todos os dados existentes
	DropCollection(ctx context.Context, collection string) error
}
