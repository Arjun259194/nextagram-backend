package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Arjun259194/nextagram-backend/types"
	"github.com/Arjun259194/nextagram-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

//structs

// user data struct to send query search result
type SearchUserDataResponse struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Email  string             `json:"email" bson:"email"`
	Gender string             `json:"gender" bson:"gender"`
}

//Authorization handler

// "/auth/register"
func postRegisterHandler(c *fiber.Ctx) error {
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

	if _, err = Storage.CreateUser(user); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while inserting into database")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusCreated, nil, "User Created")
	return c.Status(fiber.StatusCreated).JSON(res)
}

// "/auth/login"
func postLoginHandler(c *fiber.Ctx) error {
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

	result := Storage.GetOneUser(filter)

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
func postLogoutHandler(c *fiber.Ctx) error {
	cookie := utils.EmptyCookie()
	c.Cookie(cookie)
	return c.SendStatus(fiber.StatusOK)
}

//User handler

// "/user/profile"
func getUserProfileHandler(c *fiber.Ctx) error {
	userID := c.Locals("id").(primitive.ObjectID)
	filter := bson.M{"_id": userID}
	result := Storage.GetOneUser(filter)

	var user types.User
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			var errRes types.ErrorResponse
			if err == mongo.ErrNoDocuments {
				errRes = types.NewErrorResponse(fiber.StatusNotFound, err, "User not found")
				return c.Status(fiber.StatusNotFound).JSON(errRes)
			} else {
				errRes = types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while fetching user from database")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}
	}

	res := types.NewSuccessResponse(fiber.StatusOK, user, "User found")
	return c.Status(fiber.StatusOK).JSON(res)
}

// "/user/:id"
func getUserHandler(c *fiber.Ctx) error {
	strID := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "User Id not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	filter := bson.M{"_id": userID}

	result := Storage.GetOneUser(filter)

	var user types.User
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			var errRes types.ErrorResponse
			if err == mongo.ErrNoDocuments {
				errRes = types.NewErrorResponse(fiber.StatusNotFound, err, "User not found")
				return c.Status(fiber.StatusNotFound).JSON(errRes)
			} else {
				errRes = types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while fetching user from database")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}
	}

	res := types.NewSuccessResponse(fiber.StatusOK, user, "User found")
	return c.Status(fiber.StatusOK).JSON(res)
}

// "/user/profile"
func putUserProfileUpdateHandler(c *fiber.Ctx) error {
	userID := c.Locals("id").(primitive.ObjectID)

	updateBytes := c.Body()

	var updateBody types.UpgradeRouteReqBody
	if err := json.Unmarshal(updateBytes, &updateBody); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "request body not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	updateQuery := bson.M{
		"$set": bson.M{
			"name":   updateBody.Name,
			"email":  updateBody.Email,
			"gender": updateBody.Gender,
		},
	}

	result := Storage.UpdateUserById(userID, updateQuery)

	var user types.User
	if err := result.Decode(&user); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while updating in database")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusOK, user, "User data updated")
	return c.Status(fiber.StatusOK).JSON(res)
}

// "/user/:id/follow"
func putUserFollowOrUnFollowHandler(c *fiber.Ctx) error {
	clientID := c.Locals("id").(primitive.ObjectID)
	strID := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "User Id not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	isFollowed, err := Storage.ClientIDExistsInFollowers(bson.M{"_id": userID}, clientID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			errRes := types.NewErrorResponse(fiber.StatusNotFound, err, "User not found")
			return c.Status(fiber.StatusNotFound).JSON(errRes)
		} else {
			errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Database error while fetching data")
			return c.Status(fiber.StatusInternalServerError).JSON(errRes)
		}
	}

	// START FROM HERE
	var userUpdateQuery bson.M   // update query for they user whose id was sent in url
	var clientUpdateQuery bson.M // update query for the client who is sending request
	if isFollowed == true {
		userUpdateQuery = bson.M{"$pull": bson.M{"followers": clientID}}
		clientUpdateQuery = bson.M{"$pull": bson.M{"following": userID}}
	} else {
		userUpdateQuery = bson.M{"$push": bson.M{"followers": clientID}}
		clientUpdateQuery = bson.M{"$push": bson.M{"following": userID}}
	}

	userResult := Storage.UpdateUserById(userID, userUpdateQuery)
	clientResult := Storage.UpdateUserById(clientID, clientUpdateQuery)

	if userResult.Err() != nil || clientResult.Err() != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, nil, "Database error while fetching data")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusOK, nil, "OK")
	return c.Status(fiber.StatusOK).JSON(res)
}

// "/users/search?q=search+query"
func getUserSearchHandler(c *fiber.Ctx) error {
	searchQueryString := c.Query("q")

	if searchQueryString == "" || len(searchQueryString) <= 3 {
		res := types.NewSuccessResponse(fiber.StatusBadRequest, nil, "query is not valid size")
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	filter := bson.M{
		"name": bson.M{
			"$regex": searchQueryString, "$options": "i",
		},
	}

	projection := bson.M{
		"name":   1,
		"email":  1,
		"gender": 1,
	}

	result, err := Storage.GetUsersWithProjection(filter, projection)

	if err != nil {
		var (
			status int
			errRes types.ErrorResponse
		)
		if err == mongo.ErrNoDocuments {
			status = fiber.StatusNotFound
			errRes = types.NewErrorResponse(status, err, "No user found")
		} else {
			status = fiber.StatusInternalServerError
			errRes = types.NewErrorResponse(status, err, "Error while fetching from database")
		}
		return c.Status(status).JSON(errRes)
	}

	var users []SearchUserDataResponse

	if err := result.All(context.Background(), &users); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while decoding data")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusOK, users, fmt.Sprintf("Users found with query string %s", searchQueryString))
	return c.Status(fiber.StatusOK).JSON(res)

}

// "user/passwordChange"
func putUserPasswordUpdateHandler(c *fiber.Ctx) error {

	userID := c.Locals("id").(primitive.ObjectID)

	updateBytes := c.Body()

	var updateBody types.PasswordUpdateRequestBody
	if err := json.Unmarshal(updateBytes, &updateBody); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "request body not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	filter := bson.M{
		"_id": userID,
	}
	result := Storage.GetOneUser(filter)

	var user types.User
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			var errRes types.ErrorResponse
			if err == mongo.ErrNoDocuments {
				errRes = types.NewErrorResponse(fiber.StatusNotFound, err, "User not found")
				return c.Status(fiber.StatusNotFound).JSON(errRes)
			} else {
				errRes = types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while fetching user from database")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}
	}

	if err := utils.ComparePasswords(user.Password, updateBody.OldPassword); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusUnauthorized, err, "Password not matching")
		return c.Status(fiber.StatusUnauthorized).JSON(errRes)
	}

	newPassword, err := utils.HashPassword(updateBody.NewPassword)
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Failed to encrypt password")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	update := bson.M{
		"$set": bson.M{
			"password": newPassword,
		},
	}

	updateResult := Storage.UpdateUserById(userID, update)

	if err := updateResult.Err(); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Failed to update password")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusOK, nil, "Password changed")
	return c.Status(fiber.StatusOK).JSON(res)
}
