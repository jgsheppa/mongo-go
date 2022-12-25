package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Services struct{
	Magazine MagazineService
	mongo *mongo.Client
}

func NewServices(connectionString string) (*Services, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(connectionString).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	db, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &Services{
		Magazine: NewMagazineService(db),
		mongo: db,
	}, nil
}
