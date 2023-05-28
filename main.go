package main

import (
	"log"

	"github.com/Arjun259194/nextagram-backend/database"
	"github.com/Arjun259194/nextagram-backend/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	err := database.Connect("./database.db")
	if err != nil {
		log.Printf("Error connecting databse - %v", err)
	}

	defer database.Close()

	app := fiber.New()

	authGroup := app.Group("/auth")
	routes.AuthRoutes(authGroup)

	log.Fatal(app.Listen(":8080"))
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
