package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Location representa uma localização geográfica
type Location struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Municipio   string             `json:"municipio" bson:"municipio"`
	Estado      string             `json:"estado" bson:"estado"`
	Localizacao GeoJSON            `json:"localizacao" bson:"localizacao"`
	Populacao   int                `json:"populacao,omitempty" bson:"populacao,omitempty"`
}

// GeoJSON representa um ponto geográfico no formato GeoJSON
type GeoJSON struct {
	Type        string     `json:"type" bson:"type"`
	Coordinates [2]float64 `json:"coordinates" bson:"coordinates"` // [longitude, latitude]
}

// LocationResponse é a resposta da API
type LocationResponse struct {
	Municipio string  `json:"municipio"`
	Estado    string  `json:"estado"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ErrorResponse é a resposta de erro da API
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
