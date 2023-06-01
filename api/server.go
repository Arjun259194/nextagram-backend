package api

import (
	"github.com/Arjun259194/nextagram-backend/database"
	"github.com/Arjun259194/nextagram-backend/types"
	"github.com/Arjun259194/nextagram-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	ListenAddr   string
	DbConnString string
}

var Storage *database.Storage

func NewServer(listenAddr, dbConnString string) *Server {
	return &Server{
		ListenAddr:   listenAddr,
		DbConnString: dbConnString,
	}
}

func (s *Server) Start() error {
	Storage = database.NewConnection(s.DbConnString)

	Storage.Connect()
	defer Storage.Close()

	server := fiber.New()

	server.Get("/test", utils.JWTMiddleware, func(c *fiber.Ctx) error {
		userID := c.Locals("id").(primitive.ObjectID)

		result := Storage.SearchUserById(userID)

		var user types.User
		if err := result.Decode(&user); err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(fiber.StatusNotFound)
			} else {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		return c.Status(fiber.StatusOK).JSON(user)
	})

	setRoutes(server)

	return server.Listen(s.ListenAddr)
}
