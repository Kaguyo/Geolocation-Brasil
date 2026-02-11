package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	Client *mongo.Client
}

// ConnectDB estabelece conexão com MongoDB
func ConnectDB(uri, dbName string) (*MongoConnection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no MongoDB: %v", err)
	}

	// Verificar conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar conexão: %v", err)
	}

	log.Println("✅ Conectado ao MongoDB com sucesso!")

	return &MongoConnection{
		Client: client,
	}, nil
}

// Close fecha a conexão com o banco
func (db *MongoConnection) Close(ctx context.Context) error {
	return db.Client.Disconnect(ctx)
}
