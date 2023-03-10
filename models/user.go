package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var PasswordPepper string

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
}

type UserDB interface {
	ByEmail(email string) (*User, error)
}

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(db *mongo.Client) UserService {
	uM := &userMongo{db}

	return &userService{
		UserDB: uM,
	}
}

var _ UserDB = &userService{}

type userService struct {
	UserDB
}

var _ UserDB = &userMongo{}

type userMongo struct {
	db *mongo.Client
}

func (u *userMongo) ByEmail(email string) (*User, error) {
	user := User{}
	db := u.db.Database("users").Collection("authentication")

	collection := db.FindOne(context.TODO(), bson.M{"email": email})
	err := collection.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password+PasswordPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, bcrypt.ErrMismatchedHashAndPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}
