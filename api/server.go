package api

import (
	controller "github.com/Arjun259194/nextagram-backend/controllers"
	"github.com/Arjun259194/nextagram-backend/database"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	ListenAddr   string
	DbConnString string
}

var (
	storage *database.Storage
	ctrl    *controller.Controllers
)

func NewServer(listenAddr, dbConnString string) *Server {
	return &Server{
		ListenAddr:   listenAddr,
		DbConnString: dbConnString,
	}
}

func (s *Server) Start() error {
	// Creating new database connection from database package
	storage = database.NewConnection(s.DbConnString)

	storage.Connect()
	defer storage.Close()

	// Creating controllers
	ctrl = controller.NewController(storage)

	server := fiber.New()

	setRoutes(server)

	return server.Listen(s.ListenAddr)
}
