package controllers

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

// structs
type searchUserDataResponse struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Email  string             `json:"email" bson:"email"`
	Gender string             `json:"gender" bson:"gender"`
}

type userDataWithoutPassword struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id"`
	Name      string               `json:"name" bson:"name"`
	Email     string               `json:"email" bson:"email"`
	Gender    string               `json:"gender" bson:"gender"`
	Followers []primitive.ObjectID `json:"followers" bson:"followers"`
	Following []primitive.ObjectID `json:"following" bson:"following"`
}

//User handler

// "/user/profile"
func (ctrl *Controllers) GetUserProfileHandler(c *fiber.Ctx) error {
	userID := c.Locals("id").(primitive.ObjectID)
	filter := bson.M{"_id": userID}
	projection := bson.M{
		"name":      1,
		"email":     1,
		"gender":    1,
		"followers": 1,
		"following": 1,
	}
	result := ctrl.DB.GetOneUserWithProjection(filter, projection)

	var user userDataWithoutPassword
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
func (ctrl *Controllers) GetUserHandler(c *fiber.Ctx) error {
	strID := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "User Id not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	filter := bson.M{"_id": userID}

	result := ctrl.DB.GetOneUser(filter)

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
func (ctrl *Controllers) PutUserProfileUpdateHandler(c *fiber.Ctx) error {
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

	result := ctrl.DB.UpdateUserById(userID, updateQuery)

	var user types.User
	if err := result.Decode(&user); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while updating in database")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusOK, user, "User data updated")
	return c.Status(fiber.StatusOK).JSON(res)
}

// "/user/:id/follow"
func (u *Controllers) PutUserFollowOrUnFollowHandler(c *fiber.Ctx) error {
	clientID := c.Locals("id").(primitive.ObjectID)
	strID := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "User Id not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	isFollowed, err := u.DB.ClientIDExistsInFollowers(bson.M{"_id": userID}, clientID)
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

	userResult := u.DB.UpdateUserById(userID, userUpdateQuery)
	clientResult := u.DB.UpdateUserById(clientID, clientUpdateQuery)

	if userResult.Err() != nil || clientResult.Err() != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, nil, "Database error while fetching data")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusOK, nil, "OK")
	return c.Status(fiber.StatusOK).JSON(res)
}

// "/users/search?q=search+query"
func (u *Controllers) GetUserSearchHandler(c *fiber.Ctx) error {
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

	result, err := u.DB.GetUsersWithProjection(filter, projection)

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

	var users []searchUserDataResponse

	if err := result.All(context.Background(), &users); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while decoding data")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusOK, users, fmt.Sprintf("Users found with query string %s", searchQueryString))
	return c.Status(fiber.StatusOK).JSON(res)

}

// "user/passwordChange"
func (u *Controllers) PutUserPasswordUpdateHandler(c *fiber.Ctx) error {

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
	result := u.DB.GetOneUser(filter)

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

	updateResult := u.DB.UpdateUserById(userID, update)

	if err := updateResult.Err(); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Failed to update password")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusOK, nil, "Password changed")
	return c.Status(fiber.StatusOK).JSON(res)
}
