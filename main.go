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
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		return c.Status(502).Send([]byte("error getting file from form"))
// 	}

// 	err = c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename))
// 	if err != nil {
// 		return c.Status(502).Send([]byte("Error saving file"))
// 	}

// 	// Create the URL for the saved image dynamically
// 	imageURL := fmt.Sprintf("http://%s/images/%s", c.Hostname(), file.Filename)

// 	// Send the URL in the response
// 	return c.JSON(fiber.Map{
// 		"url": imageURL,
// 	})
// })
