package api

import (
	"encoding/json"
	"fmt"

	"github.com/Arjun259194/nextagram-backend/types"
	"github.com/Arjun259194/nextagram-backend/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	BAD_REQUEST           = 400 // Invalid request syntax or parameters
	NOT_FOUND             = 404 // Requested data is not found
	BAD_GATEWAY           = 502 // Invalid response from an upstream server
	OK                    = 200 // Successful request processing
	CREATED               = 201 // New resource successfully created
	INTERNAL_SERVER_ERROR = 500 // Unexpected error on the server side
	UNAUTHORIZED          = 401 // client is not authorized for the data it requires
)

func PostRegisterHandler(c *fiber.Ctx) error {
	fmt.Println("Register handler called")
	reqBodyByte := c.Body()
	var reqBody types.RegisterRequestBody

	if err := json.Unmarshal(reqBodyByte, &reqBody); err != nil {
		errRes := types.NewErrorResponse(BAD_REQUEST, err, "Request Body is not valid")
		return c.Status(BAD_REQUEST).JSON(errRes)
	}

	if err := reqBody.Validate(); err != nil {
		errRes := types.NewErrorResponse(BAD_REQUEST, err, "Request Body is not valid")
		return c.Status(BAD_REQUEST).JSON(errRes)
	}

	hashedPassword, err := utils.HashPassword(reqBody.Password)
	if err != nil {
		errRes := types.NewErrorResponse(INTERNAL_SERVER_ERROR, err, "Failed to encrypt password")
		return c.Status(INTERNAL_SERVER_ERROR).JSON(errRes)
	}

	user := types.NewUser(reqBody.Name, reqBody.Email, reqBody.Gender, hashedPassword)

	result, err := Storage.CreateUser(user)

	if err != nil {
		fmt.Printf("Error while inserting into database - %v\n", err)
		errRes := types.NewErrorResponse(INTERNAL_SERVER_ERROR, err, "Error while inserting into database")
		return c.Status(INTERNAL_SERVER_ERROR).JSON(errRes)
	}

	fmt.Printf("\nInserted with id - %v\n", result.InsertedID)

	res := types.NewSuccessResponse(CREATED, nil, "User Created")
	return c.Status(CREATED).JSON(res)
}

func PostLoginHandler(c *fiber.Ctx) error {
	fmt.Println("Login handler called")
	reqBodyByte := c.Body()
	var reqBody types.LoginRequestBody

	if err := json.Unmarshal(reqBodyByte, &reqBody); err != nil {
		errRes := types.NewErrorResponse(BAD_REQUEST, err, "Request Body is not valid")
		return c.Status(BAD_REQUEST).JSON(errRes)
	}

	if err := reqBody.Validate(); err != nil {
		errRes := types.NewErrorResponse(BAD_REQUEST, err, "Request Body is not valid")
		return c.Status(BAD_REQUEST).JSON(errRes)
	}

	result := Storage.SearchUserByEmail(reqBody.Email)

	var foundUser types.User
	err := result.Decode(&foundUser)
	if err != nil {
		fmt.Printf("\nerror while marshaling foundUser data - %v\n", err)
		errRes := types.NewErrorResponse(NOT_FOUND, err, "User not found")
		return c.Status(NOT_FOUND).JSON(errRes)
	}

	if err = utils.ComparePasswords(foundUser.Password, reqBody.Password); err != nil {
		fmt.Printf("Password not matched - %v", err)
		errRes := types.NewErrorResponse(UNAUTHORIZED, err, "Wrong password")
		return c.Status(UNAUTHORIZED).JSON(errRes)
	}

	token, err := utils.GenerateToken(foundUser.ID)
	if err != nil {
		errRes := types.NewErrorResponse(INTERNAL_SERVER_ERROR, err, "Error while generating token")
		c.Status(INTERNAL_SERVER_ERROR).JSON(errRes)
	}

	accessCookie := utils.NewHTTPOnlyCookie(token)

	c.Cookie(accessCookie)

	res := types.NewSuccessResponse(OK, nil, "Logged In")
	return c.Status(OK).JSON(res)
}
