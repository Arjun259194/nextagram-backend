package api

import (
	"github.com/Arjun259194/nextagram-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func setRoutes(server *fiber.App) {
	setAuthRoutes(server)
	setUserRoutes(server)
}

func setAuthRoutes(server *fiber.App) {
	//authorization routes
	server.Post("/auth/register", postRegisterHandler)
	server.Post("/auth/login", postLoginHandler)
	server.Post("/auth/logout", postLogoutHandler)
}

func setUserRoutes(server *fiber.App) {
	server.Get("/user/profile", utils.JWTMiddleware, getUserProfileHandler)
	server.Get("/user/:id", utils.JWTMiddleware, getUserHandler)
	server.Get("/users/search", utils.JWTMiddleware, getUserSearchHandler)
	server.Put("/user/:id/follow", utils.JWTMiddleware, putUserFollowOrUnFollowHandler)
	server.Put("/user/profile", utils.JWTMiddleware, putUserProfileUpdateHandler)
	server.Put("/user/password", utils.JWTMiddleware, putUserPasswordUpdateHandler)
}
