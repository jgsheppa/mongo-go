package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	Create(magazine Magazine) (*mongo.InsertOneResult, error)
	UpdateById(magazine Magazine) (*mongo.UpdateResult, error)
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
	db := mM.db.Database("library").Collection("magazines")

	collection := db.FindOne(context.TODO(), bson.M{"_id": objectId})
	err = collection.Decode(&magazine)
	if err != nil {
		return nil, err
	}

	return &magazine, nil
}

func (mM *mongoMagazine) FindByTitle(title string) (*Magazine, error) {
	magazine := Magazine{}
	db := mM.db.Database("library").Collection("magazines")

	collection := db.FindOne(context.TODO(), bson.M{"title": title})
	err := collection.Decode(&magazine)
	if err != nil {
		return nil, err
	}

	return &magazine, nil
}

func (mM *mongoMagazine) FindAll() (*[]Magazine, error) {
	magazines := make([]Magazine, 2)
	db := mM.db.Database("library").Collection("magazines")

	collection, err := db.Find(context.TODO(), bson.M{})
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

	db := mM.db.Database("library").Collection("magazines")

	res, err := db.DeleteOne(context.Background(), bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (mM *mongoMagazine) Create(magazine Magazine) (*mongo.InsertOneResult, error) {
	db := mM.db.Database("library").Collection("magazines")

	res, err := db.InsertOne(context.Background(), magazine)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (mM *mongoMagazine) UpdateById(magazine Magazine) (*mongo.UpdateResult, error) {
	db := mM.db.Database("library").Collection("magazines")
	payload := bson.D{{Key: "$set", Value: magazine}}

	res, err := db.UpdateByID(context.Background(), magazine.ID, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}
