package controllers

import (
	"github.com/Arjun259194/nextagram-backend/database"
)

type Controllers struct {
	DB *database.Storage
}

func NewController(conn *database.Storage) *Controllers {
	return &Controllers{
		DB: conn,
	}
}
