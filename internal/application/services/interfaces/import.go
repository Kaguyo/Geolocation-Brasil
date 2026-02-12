package interfaces

import (
	"context"

	domain "github.com/Kaguyo/Geolocation-Brasil/internal/domain/entities"
)

type IImportService interface {
	ImportBrazilianCitiesExampleTest(ctx context.Context) error
	ImportData(ctx context.Context, filename string) error
	GetLocationByName(ctx context.Context, municipio, estado string) (*domain.Location, error)
	GetLocationsInKilometersRange(ctx context.Context, longitude, latitude float64, rangeInKilometers float64) (*[]domain.Location, error)
	// CreateGeoIndex cria índice geoespacial
	CreateGeoIndex(ctx context.Context) error
	// CreateTextIndex cria índice de texto para busca
	CreateTextIndex(ctx context.Context) error
	// ResetCollection recria a coleção, removendo todos os dados existentes
	ResetCollection(ctx context.Context, collection string) error
}
