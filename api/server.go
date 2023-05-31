package api

import (
	"github.com/Arjun259194/nextagram-backend/database"
	"github.com/gofiber/fiber/v2"
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

	setRoutes(server)

	return server.Listen(s.ListenAddr)
}

func setRoutes(server *fiber.App) {
	server.Post("auth/register", PostRegisterHandler)
	server.Post("auth/login", PostLoginHandler)
}
