package controller

import (
	"encoding/json"
	"fmt"

	"github.com/Arjun259194/nextagram-backend/database"
	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

type registerInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterController(c *fiber.Ctx) error {

	requestBody := c.Body()

	registerationData := &registerInput{}

	if err := json.Unmarshal(requestBody, registerationData); err != nil {
		fmt.Printf("Error will unmarshaling request body - %v", err)
		return c.Status(400).SendString("failded to get data from request body. check request body")
	}

	if err := database.AddUser(registerationData.Name, registerationData.Email, registerationData.Password); err != nil {
		fmt.Printf("Error in AddUser function - %v", err)
		return c.Status(504).SendString("Error at database")
	}

	return c.Status(200).SendString("User registered")
}
