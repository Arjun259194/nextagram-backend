package controllers

import (
	"fmt"

	"github.com/Arjun259194/nextagram-backend/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ctrl *Controllers) GetTopPosts(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (ctrl *Controllers) PostCreatePost(c *fiber.Ctx) error {
	userID := c.Locals("id").(primitive.ObjectID)

	file, err := c.FormFile("file")
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "Can't get file from formdata")
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
