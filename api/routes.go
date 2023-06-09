package api

import (
	"github.com/Arjun259194/nextagram-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func setRoutes(server *fiber.App) {
	server.Static("/images", "./uploads") // static path to server images from server

	setAuthRoutes(server)
	setUserRoutes(server)
	setPostRoutes(server)
}

func setAuthRoutes(server *fiber.App) {
	//authorization routes
	server.Post("/auth/register", ctrl.PostRegisterHandler)
	server.Post("/auth/login", ctrl.PostLoginHandler)
	server.Post("/auth/logout", ctrl.PostLogoutHandler)
}

func setUserRoutes(server *fiber.App) {
	server.Get("/user/profile", utils.JWTMiddleware, ctrl.GetUserProfileHandler)
	server.Get("/user/:id", utils.JWTMiddleware, ctrl.GetUserHandler)
	server.Get("/users/search", utils.JWTMiddleware, ctrl.GetUserSearchHandler)
	server.Put("/user/:id/follow", utils.JWTMiddleware, ctrl.PutUserFollowOrUnFollowHandler)
	server.Put("/user/profile", utils.JWTMiddleware, ctrl.PutUserProfileUpdateHandler)
	server.Put("/user/password", utils.JWTMiddleware, ctrl.PutUserPasswordUpdateHandler)
}

func setPostRoutes(server *fiber.App) {
	server.Get("/posts", utils.JWTMiddleware, ctrl.GetTopPosts)
  server.Post("/post", utils.JWTMiddleware, ctrl.PostCreatePost)
}
