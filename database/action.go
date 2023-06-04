package database

import (
	"github.com/Arjun259194/nextagram-backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (s *Storage) CreateUser(user *types.User) (*mongo.InsertOneResult, error) {
	context := s.Ctx
	return s.UserModel.InsertOne(context, user)
}

func (s *Storage) GetOneUser(filter bson.M) *mongo.SingleResult {
	context := s.Ctx
	return s.UserModel.FindOne(context, filter)
}

func (s *Storage) GetOneUserWithProjection(filter bson.M, projection bson.M) *mongo.SingleResult {
	context := s.Ctx
	return s.UserModel.FindOne(context, filter, options.FindOne().SetProjection(projection))
}

func (s *Storage) GetUsers(filter bson.M) (*mongo.Cursor, error) {
	context := s.Ctx
	return s.UserModel.Find(context, filter)
}

func (s *Storage) GetUsersWithProjection(filter bson.M, projection bson.M) (*mongo.Cursor, error) {
	context := s.Ctx
	return s.UserModel.Find(context, filter, options.Find().SetProjection(projection))
}

func (s *Storage) UpdateUserById(id primitive.ObjectID, update bson.M) *mongo.SingleResult {
	filter := bson.M{"_id": id}
	context := s.Ctx
	return s.UserModel.FindOneAndUpdate(context, filter, update)
}

func (s *Storage) ClientIDExistsInFollowers(filter bson.M, clientID primitive.ObjectID) (bool, error) {
	result := s.UserModel.FindOne(s.Ctx, filter)

	if err := result.Err(); err != nil {
		return false, err
	}

	var user types.User // Assuming you have a User struct defined
	err := result.Decode(&user)
	if err != nil {
		return false, err
	}

	// Check if the client ID exists in the followers array
	for _, follower := range user.Followers {
		if follower == clientID {
			return true, nil
		}
	}

	return false, nil
}
