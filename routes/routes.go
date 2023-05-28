package routes

import (
	"github.com/Arjun259194/nextagram-backend/controller"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(group fiber.Router) {
  group.Post("/register", controller.RegisterController)
}

