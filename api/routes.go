package api

import (
	"github.com/gofiber/fiber/v2"
)

func setRoutes(server *fiber.App) {
	setAuthRoutes(server)
}

func setAuthRoutes(server *fiber.App) {
	//authorization routes
	server.Post("/auth/register", postRegisterHandler)
	server.Post("/auth/login", postLoginHandler)
	server.Post("/auth/logout", postLogoutHandler)
}

func setUserRoutes(server *fiber.App) {
	server.Get("/user/profile", getUserProfileHandler)
	server.Get("/user/:id", getUserHandler)
	server.Get("/user/search", getUserSearchHandler)
	server.Put("/user/:id/follow", putUserFollowOrUnFollowHandler)
	server.Put("/user/profile", putUserProfileUpdateHandler)
	server.Put("/user/password", putUserPasswordUpdateHandler)
}
