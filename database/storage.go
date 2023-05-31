package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct {
	ConnectionString string
	Client           *mongo.Client
	UserModel        *mongo.Collection
	Ctx              context.Context
}

func NewConnection(connectionString string) *Storage {
	return &Storage{
		ConnectionString: connectionString,
		Ctx:              context.Background(),
	}
}

func (s *Storage) Connect() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(s.ConnectionString).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(s.Ctx, opts)
	if err != nil {
		log.Fatalf("Error while connecting to database in storage.go - %v", err)
	}

	if err = client.Ping(s.Ctx, nil); err != nil {
		log.Fatalf("Error while sending ping to database - %v", err)
	}

	fmt.Println("Ping Done!, database connected")

	s.Client = client
	s.UserModel = client.Database("nextagram").Collection("user")

	// Create unique index on email field
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = s.UserModel.Indexes().CreateOne(s.Ctx, indexModel)
	if err != nil {
		log.Fatalf("Error while creating unique index on email field - %v", err)
	}
}

func (s *Storage) Close() {
	if err := s.Client.Disconnect(s.Ctx); err != nil {
		log.Fatalf("Error while disconnecting database in storage.go - %v", err)
	}
}
