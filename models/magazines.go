package models

import (
	"context"
	"fmt"
	"time"

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
	AggregateByPrice(price string) (*[]Magazine, error)
	// CRUD operations
	Create(magazine Magazine) (*mongo.InsertOneResult, error)
	FindById(id string) (*Magazine, error)
	FindBySlug(slug string) (*Magazine, error)
	FindAll() (*[]Magazine, error)
	UpdateById(magazine Magazine) (*mongo.UpdateResult, error)
	Delete(id string) (*mongo.DeleteResult, error)
	// Search
	Search(field, term string) (*[]Magazine, error)
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

func (mM *mongoMagazine) FindBySlug(slug string) (*Magazine, error) {
	magazine := Magazine{}
	db := mM.db.Database("library").Collection("magazines")

	collection := db.FindOne(context.TODO(), bson.M{"title": slug})
	err := collection.Decode(&magazine)
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

func (mM *mongoMagazine) AggregateByPrice(price string) (*[]Magazine, error) {
	db := mM.db.Database("library").Collection("magazines")
	groupStage := bson.D{{Key: "$match", Value: bson.D{{Key: "price", Value: price}}}}

	res, err := db.Aggregate(context.Background(), mongo.Pipeline{groupStage})
	if err != nil {
		return nil, err
	}

	magazines := make([]Magazine, 2)

	err = res.All(context.Background(), &magazines)
	if err != nil {
		return nil, err
	}

	return &magazines, nil
}

func (mM *mongoMagazine) Search(field, term string) (*[]Magazine, error) {
	searchQuery := bson.D{{Key: "index", Value: "magazine_title"},
		{Key: "autocomplete", Value: bson.D{
			{Key: "path", Value: field},
			{Key: "query", Value: term},
		}}}
	searchStage := bson.D{{Key: "$search", Value: searchQuery}}
	limitStage := bson.D{{Key: "$limit", Value: 5}}
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{{Key: "score", Value: bson.D{{Key: "$meta", Value: "searchScore"}}},
			{Key: field, Value: 1}, {Key: "_id", Value: 1}, {Key: "price", Value: 1},
			{Key: "highlight", Value: bson.D{{Key: "$meta", Value: "searchHighlights"}}}}}}

	opts := options.Aggregate().SetMaxTime(5 * time.Second)
	// run pipeline
	db := mM.db.Database("library").Collection("magazines")
	res, err := db.Aggregate(context.Background(), mongo.Pipeline{searchStage, limitStage, projectStage}, opts)
	if err != nil {
		return nil, err
	}

	magazines := make([]Magazine, 2)

	err = res.All(context.Background(), &magazines)
	if err != nil {
		return nil, err
	}

	fmt.Printf("magazines: %v", magazines)
	return &magazines, nil
}
