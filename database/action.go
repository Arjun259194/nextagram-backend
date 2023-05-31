package database

import (
	"fmt"

	"github.com/Arjun259194/nextagram-backend/types"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func (s *Storage) CreateUser(user *types.User) (*mongo.InsertOneResult, error) {
	return s.UserModel.InsertOne(s.Ctx, user)
}

func (s *Storage) SearchUserByEmail(email string) *mongo.SingleResult {
	result := s.UserModel.FindOne(s.Ctx, bson.M{"email": email})
	fmt.Println(result)
	return result
}
