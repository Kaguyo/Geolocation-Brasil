package bootstrap

import (
	"net/http"

	handlers "github.com/Kaguyo/Geolocation-Brasil/internal/api"
	"github.com/Kaguyo/Geolocation-Brasil/internal/application/services"
	"github.com/Kaguyo/Geolocation-Brasil/internal/domain/interfaces"
	"github.com/Kaguyo/Geolocation-Brasil/internal/infrastructure/mongodb"
)

type Application struct {
	DB      *mongodb.Database
	Router  http.Handler
	Service services.ImportService
}

func Build(mongoURI, dbName, collection string) (*Application, error) {

	db, err := mongodb.ConnectDB(mongoURI, dbName)
	if err != nil {
		return nil, err
	}

	var geoRepository interfaces.IGeoRepository
	geoRepository = mongodb.NewGeoRepository(db.Database)
	geoService := services.NewGeoService(geoRepository)
	geoHandler := handlers.NewAPI(geoService)
	router := geoHandler.SetupRoutes()

	return &Application{
		DB:      db,
		Router:  router,
		Service: *geoService,
	}, nil
}
