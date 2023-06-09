package main

import (
	"log"

	"github.com/Arjun259194/nextagram-backend/api"
)

const (
	PORT           string = ":8080"
	DB_CONN_STRING string = "mongodb+srv://mongo:arjun259@cluster0.12gakmk.mongodb.net/mgmt_db"
)

func main() {
	server := api.NewServer(PORT, DB_CONN_STRING)
	log.Fatal(server.Start())
}

// app.Post("/upload", func(c *fiber.Ctx) error {
// })


