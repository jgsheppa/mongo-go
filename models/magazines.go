package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Magazine struct {
	Title string `json:"title"`
	Price string `json:"price"`
}

type MagazineDB interface {
	FindById(id string) (*Magazine, error)
}

type MagazineService interface {
	MagazineDB
}

func NewMagazineService(db *mongo.Client) MagazineService {
	mDb := &mongoMagazine{db}

	return &magazineService{
		MagazineDB: mDb,
	}
}

var _ MagazineDB = &magazineService{}

type magazineService struct {
	MagazineDB
}

var _ MagazineDB = &mongoMagazine{}

type mongoMagazine struct {
	db *mongo.Client
}

func (mM *mongoMagazine) FindById(id string) (*Magazine, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	magazine := Magazine{}

	collection := mM.db.Database("library").Collection("magazines").FindOne(context.TODO(), bson.M{"_id": objectId})
	err = collection.Decode(&magazine)
	if err != nil {
		return nil, err
	}

	return &magazine, nil
}
