package mongodb

import (
	"context"
	"fmt"
	"log"

	domain "github.com/Kaguyo/Geolocation-Brasil/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GeoRepository struct {
	collection *mongo.Collection
}

func NewGeoRepository(db *mongo.Database) *GeoRepository {
	return &GeoRepository{
		collection: db.Collection("geolocations"),
	}
}

// CreateGeoIndex cria índice geoespacial
func (gr *GeoRepository) CreateGeoIndex(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "localizacao", Value: "2dsphere"}},
	}

	_, err := gr.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("erro ao criar índice geoespacial: %v", err)
	}

	log.Println("✅ Índice geoespacial criado!")
	return nil
}

// CreateTextIndex cria índice de texto para busca
func (gr *GeoRepository) CreateTextIndex(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "municipio", Value: "text"},
			{Key: "estado", Value: "text"},
		},
	}

	_, err := gr.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("erro ao criar índice de texto: %v", err)
	}

	log.Println("✅ Índice de texto criado!")
	return nil
}

// Inserts as many locations as given through parameter
func (gr *GeoRepository) InsertLocations(ctx context.Context, locationBuffer []domain.Location) error {
	if len(locationBuffer) == 0 {
		return nil
	}

	// Normalizar municipios antes de inserir
	for i := range locationBuffer {
		locationBuffer[i].Municipio = utils.NormalizeMunicipio(locationBuffer[i].Municipio)
	}

	// Converte []domain.Location para []interface{} exigido pelo InsertMany
	documents := make([]interface{}, len(locationBuffer))
	for i, loc := range locationBuffer {
		documents[i] = loc
	}

	_, err := gr.collection.InsertMany(ctx, documents)
	if err != nil {
		return err
	}

	return nil
}

// Inserts as many locations as given through parameter
func (gr *GeoRepository) GetLocationsInKilometersRange(ctx context.Context, longitude, latitude float64, maxDistanceKm float64) (*[]domain.Location, error) {
	filter := bson.M{
		"localizacao": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longitude, latitude},
				},
				"$maxDistance": maxDistanceKm * 1000, // converter km para metros
			},
		},
	}

	cursor, err := gr.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var locations []domain.Location
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, err
	}

	return &locations, nil
}

// Funcionalidade de teste de importação de localidades
func (gr *GeoRepository) ImportTest(ctx context.Context, locations []domain.Location) error {

	// Limpar coleção antes de importar
	if err := gr.collection.Drop(ctx); err != nil {
		log.Printf("Aviso ao limpar coleção: %v", err)
	}

	// Normalizar municipios antes de inserir
	for i := range locations {
		locations[i].Municipio = utils.NormalizeMunicipio(locations[i].Municipio)
	}

	// Converter para []interface{}
	docs := make([]interface{}, len(locations))
	for i, location := range locations {
		docs[i] = location
	}

	_, err := gr.collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("erro ao inserir cidades: %v", err)
	}

	log.Printf("✅ %d cidades importadas com sucesso!", len(locations))
	return nil
}

// GetLocationByName busca localização por nome de município ou estado
func (gr *GeoRepository) GetLocationByName(ctx context.Context, municipio, estado string) (*domain.Location, error) {
	// Construir filtro dinamicamente: se estado for vazio, apenas buscar por municipio
	filter := bson.M{
		"municipio": municipio,
	}

	// Se estado foi fornecido, adicionar ao filtro
	if estado != "" {
		filter["estado"] = estado
	}

	// Opções: ordenar por população descrescente para retornar o resultado mais relevante
	opts := options.FindOne().SetSort(bson.M{"populacao": -1})

	var location domain.Location

	err := gr.collection.FindOne(ctx, filter, opts).Decode(&location)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &location, nil
}

func (gr *GeoRepository) DropCollection(ctx context.Context, collection string) error {
	err := gr.collection.Drop(ctx)
	if err != nil {
		return fmt.Errorf("erro ao deletar dropar collection: %v", err)
	}
	log.Println("✅ Coleção deletada com sucesso!")
	return nil
}
