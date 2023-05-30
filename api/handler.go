package api

import (
	"encoding/json"
	"fmt"

	"github.com/Arjun259194/nextagram-backend/types"
	"github.com/Arjun259194/nextagram-backend/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	BAD_REQUEST           = 400
	BAD_GATEWAY           = 502
	OK                    = 200
	CREATED               = 201
	INTERNAL_SERVER_ERROR = 500
)

type registerRequestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
}

func PostRegisterHandler(c *fiber.Ctx) error {
	reqBodyByte := c.Body()
	var reqBody registerRequestBody

	if err := json.Unmarshal(reqBodyByte, &reqBody); err != nil {
		return c.Status(BAD_REQUEST).SendString("Request Body is not valid")
	}

	hashedPassword, err := utils.HashPassword(reqBody.Password)
	if err != nil {
		return c.Status(INTERNAL_SERVER_ERROR).SendString("Failed to encrypt password")
	}

	user := types.NewUser(reqBody.Name, reqBody.Email, reqBody.Gender, hashedPassword)

	fmt.Println("This is the user: ", user)

	result, err := Storage.CreateUser(user)

	if err != nil {
		fmt.Printf("Error while inserting into database - %v\n", err)
		fmt.Printf(err.Error())
		return c.Status(INTERNAL_SERVER_ERROR).SendString("Error while inserting into database")
	}

	fmt.Printf("Inserted with id - %v\n", result.InsertedID)

	return c.Status(CREATED).SendString("User Created")
}
