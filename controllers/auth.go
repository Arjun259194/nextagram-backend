package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/Arjun259194/nextagram-backend/types"
	"github.com/Arjun259194/nextagram-backend/utils"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2/bson"
)

// Authorization handler
// "/auth/register"
func (ctrl *Controllers) PostRegisterHandler(c *fiber.Ctx) error {
	reqBodyByte := c.Body()
	var reqBody types.RegisterRequestBody

	if err := json.Unmarshal(reqBodyByte, &reqBody); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "Request Body is not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	if err := reqBody.Validate(); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "Request Body is not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	hashedPassword, err := utils.HashPassword(reqBody.Password)
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Failed to encrypt password")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	user := types.NewUser(reqBody.Name, reqBody.Email, reqBody.Gender, hashedPassword)

	if _, err = ctrl.DB.CreateUser(user); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while inserting into database")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusCreated, nil, "User Created")
	return c.Status(fiber.StatusCreated).JSON(res)
}

// "/auth/login"
func (ctrl *Controllers) PostLoginHandler(c *fiber.Ctx) error {
	reqBodyByte := c.Body()
	var reqBody types.LoginRequestBody

	if err := json.Unmarshal(reqBodyByte, &reqBody); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "Request Body is not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	if err := reqBody.Validate(); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "Request Body is not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	filter := bson.M{"email": reqBody.Email}

	result := ctrl.DB.GetOneUser(filter)

	var foundUser types.User

	if err := result.Decode(&foundUser); err != nil {
		fmt.Printf("\nerror while marshaling foundUser data - %v\n", err)
		errRes := types.NewErrorResponse(fiber.StatusNotFound, err, "User not found")
		return c.Status(fiber.StatusNotFound).JSON(errRes)
	}

	if err := utils.ComparePasswords(foundUser.Password, reqBody.Password); err != nil {
		fmt.Printf("Password not matched - %v", err)
		errRes := types.NewErrorResponse(fiber.StatusUnauthorized, err, "Wrong password")
		return c.Status(fiber.StatusUnauthorized).JSON(errRes)
	}

	token, err := utils.GenerateToken(foundUser.ID)
	if err != nil {
		fmt.Printf("Error while creating json web token - %v", err)
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while generating token")
		c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	accessCookie := utils.NewHTTPOnlyCookie(token)

	c.Cookie(accessCookie)

	res := types.NewSuccessResponse(fiber.StatusOK, nil, "Logged In")

	return c.Status(fiber.StatusOK).JSON(res)
}

// "/auth/logout"
func (ctrl *Controllers) PostLogoutHandler(c *fiber.Ctx) error {
	cookie := utils.EmptyCookie()
	c.Cookie(cookie)
	return c.SendStatus(fiber.StatusOK)
}
