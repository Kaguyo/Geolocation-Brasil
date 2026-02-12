package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Kaguyo/Geolocation-Brasil/internal/application/services/interfaces"
	domain "github.com/Kaguyo/Geolocation-Brasil/internal/domain/entities"
)

type API struct {
	importService interfaces.IImportService
}

// NewAPI cria uma nova instância da API
func NewAPI(service interfaces.IImportService) *API {
	return &API{importService: service}
}

// GetLocationByNameHandler busca localização por nome
func (api *API) GetLocationByNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	municipio := vars["municipio"]
	estado := r.URL.Query().Get("estado")

	// Normalizar entrada do usuário: capitalize cada palavra
	municipio = strings.ToTitle(strings.ToLower(municipio))
	estado = strings.ToUpper(estado)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	location, err := api.importService.GetLocationByName(ctx, municipio, estado)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			respondWithError(w, http.StatusNotFound, "Localização não encontrada")
			return
		}
		log.Printf("Erro ao buscar localização: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro ao buscar localização")
		return
	}

	if location == nil {
		respondWithError(w, http.StatusNotFound, "Localização não encontrada")
		return
	}

	if len(location.Localizacao.Coordinates) < 2 {
		respondWithError(w, http.StatusInternalServerError, "Localização sem coordenadas válidas")
		return
	}

	response := domain.LocationResponse{
		Municipio: location.Municipio,
		Estado:    location.Estado,
		Latitude:  location.Localizacao.Coordinates[1],
		Longitude: location.Localizacao.Coordinates[0],
	}

	respondWithJSON(w, http.StatusOK, response)
}

// GetNearbyLocationsHandler busca localizações próximas
func (api *API) GetNearbyLocationsHandler(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	distStr := r.URL.Query().Get("distance")

	if latStr == "" || lonStr == "" {
		respondWithError(w, http.StatusBadRequest, "Parâmetros lat e lon são obrigatórios")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Latitude inválida")
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Longitude inválida")
		return
	}

	distance := 50.0 // padrão 50km
	if distStr != "" {
		distance, err = strconv.ParseFloat(distStr, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Distância inválida")
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	locations, err := api.importService.GetLocationsInKilometersRange(ctx, lon, lat, distance)
	if err != nil {
		log.Printf("Erro ao buscar localizações próximas: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro ao buscar localizações")
		return
	}

	responses := make([]domain.LocationResponse, len(*locations))
	for i, loc := range *locations {
		responses[i] = domain.LocationResponse{
			Municipio: loc.Municipio,
			Estado:    loc.Estado,
			Latitude:  loc.Localizacao.Coordinates[1],
			Longitude: loc.Localizacao.Coordinates[0],
		}
	}

	respondWithJSON(w, http.StatusOK, responses)
}

// HealthCheckHandler verifica se a API está funcionando
func (api *API) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "API de Geolocalização Brasil está funcionando!",
	})
}

// respondWithJSON envia resposta JSON
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao gerar JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// respondWithError envia resposta de erro
func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

// SetupRoutes configura as rotas da API
func (api *API) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Middleware de logging
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)

	// Rotas
	router.HandleFunc("/health", api.HealthCheckHandler).Methods("GET")
	router.HandleFunc("/location/{municipio}", api.GetLocationByNameHandler).Methods("GET")
	router.HandleFunc("/nearby", api.GetNearbyLocationsHandler).Methods("GET")

	return router
}

// loggingMiddleware registra todas as requisições
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("Concluído em %v", time.Since(start))
	})
}

// corsMiddleware adiciona headers CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
