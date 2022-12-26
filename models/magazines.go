package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Magazine struct {
	ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title string             `bson:"title" json:"title"`
	Price string             `bson:"price" json:"price"`
}

type MagazineDB interface {
	FindById(id string) (*Magazine, error)
	FindAll() (*[]Magazine, error)
	Delete(id string) (*mongo.DeleteResult, error)
	Create(magazine Magazine) (*mongo.UpdateResult, error)
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

func (mM *mongoMagazine) FindAll() (*[]Magazine, error) {
	magazines := make([]Magazine, 2)

	collection, err := mM.db.Database("library").Collection("magazines").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	err = collection.All(context.TODO(), &magazines)
	if err != nil {
		return nil, err
	}

	return &magazines, nil
}

func (mM *mongoMagazine) Delete(id string) (*mongo.DeleteResult, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res, err := mM.db.Database("library").Collection("magazines").DeleteOne(context.Background(), bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (mM *mongoMagazine) Create(magazine Magazine) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	
	payload := bson.D{{Key: "$set", Value: magazine}}

	res, err := mM.db.Database("library").Collection("magazines").UpdateOne(context.Background(), bson.M{"_id": magazine.ID}, payload, opts)

	if err != nil {
		return nil, err
	}

	return res, nil
}
