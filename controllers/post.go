package controllers

import (
	"fmt"

	"github.com/Arjun259194/nextagram-backend/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (ctrl *Controllers) GetTopPosts(c *fiber.Ctx) error {

	options := options.Find().SetLimit(10)
	filter := bson.M{}

	result, err := ctrl.DB.PostModel.Find(ctrl.DB.Ctx, filter, options)

	if err != nil {
		var (
			status int
			errRes types.ErrorResponse
		)
		if err == mongo.ErrNoDocuments {
			status = fiber.StatusNotFound
			errRes = types.NewErrorResponse(status, err, "No post found")
		} else {
			status = fiber.StatusInternalServerError
			errRes = types.NewErrorResponse(status, err, "Error while fetching from database")
		}
		return c.Status(status).JSON(errRes)
	}

	var posts []types.Post
	if err := result.All(ctrl.DB.Ctx, &posts); err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while decoding data")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusOK, posts, "top 10 latest posts")
	return c.Status(fiber.StatusOK).JSON(res)

}

func (ctrl *Controllers) PostCreatePost(c *fiber.Ctx) error {
	userID := c.Locals("id").(primitive.ObjectID)

	file, err := c.FormFile("file")
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "Can't get file from formData")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	title := c.FormValue("title", "")

	if title == "" {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, nil, "Title not found")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	err = c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename))
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while saving file")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	// Create the URL for the saved image dynamically
	imageURL := fmt.Sprintf("http://%s/images/%s", c.Hostname(), file.Filename)

	newPost := types.NewPost(title, imageURL, userID)

	_, err = ctrl.DB.PostModel.InsertOne(ctrl.DB.Ctx, newPost)

	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusInternalServerError, err, "Error while inserting into database")
		return c.Status(fiber.StatusInternalServerError).JSON(errRes)
	}

	res := types.NewSuccessResponse(fiber.StatusCreated, nil, "Post Created")
	return c.Status(fiber.StatusCreated).JSON(res)
}
