package database

import (
	"context"

	"github.com/Arjun259194/nextagram-backend/types"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) CreateUser(user *types.User) (*mongo.InsertOneResult, error) {
	collection := s.Client.Database("nextagram").Collection("user")
	return collection.InsertOne(context.TODO(), user)
}
