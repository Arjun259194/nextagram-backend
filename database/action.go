package database

import (
	"github.com/Arjun259194/nextagram-backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func (s *Storage) CreateUser(user *types.User) (*mongo.InsertOneResult, error) {
	return s.UserModel.InsertOne(s.Ctx, user)
}

func (s *Storage) SearchUserByEmail(email string) *mongo.SingleResult {
	filter := bson.M{"email": email}
	context := s.Ctx
	return s.UserModel.FindOne(context, filter)
}

func (s *Storage) SearchUserById(id primitive.ObjectID) *mongo.SingleResult {
	filter := bson.M{"_id": id}
	context := s.Ctx
	return s.UserModel.FindOne(context, filter)
}
