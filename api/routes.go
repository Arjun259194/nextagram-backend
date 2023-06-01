package api

import (
	"github.com/gofiber/fiber/v2"
)

func setRoutes(server *fiber.App) {
	setAuthRoutes(server)
}

func setAuthRoutes(server *fiber.App) {
	//authorization routes
	server.Post("/auth/register", PostRegisterHandler)
	server.Post("/auth/login", PostLoginHandler)
	server.Post("/auth/logout", PostLogoutHandler)
}
