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
}
