package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client   *mongo.Client
	Database *mongo.Database
}

// ConnectDB estabelece conexão com MongoDB
func ConnectDB(uri, dbName string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("erro ao verificar conexão: %w", err)
	}

	log.Println("✅ Conectado ao MongoDB com sucesso!")

	return &Database{
		client:   client,
		Database: client.Database(dbName),
	}, nil
}

// Close fecha a conexão com o banco
func (db *Database) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}
