package server

import (
	"github.com/TimothyYe/godns/internal/server/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	addr       string
	app        *fiber.App
	controller *controllers.Controller
}

func (s *Server) SetAddress(addr string) *Server {
	s.addr = addr
	return s
}

func (s *Server) Build() {
	config := fiber.Config{}
	s.app = fiber.New(config)
	s.controller = controllers.NewController()
}

func (s *Server) Start() error {
	s.initRoutes()

	log.Infof("Server is listening on port: %s", s.addr)
	return s.app.Listen(s.addr)
}

func (s *Server) initRoutes() {
	// set cross domain access rules
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST ,PUT ,DELETE, OPTIONS",
		AllowHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}))

	// Create routes group.
	route := s.app.Group("/api/v1")

	// authMiddleware := keyauth.New(keyauth.Config{
	// 	Validator: func(c *fiber.Ctx, key string) (bool, error) {
	// 		hashedAPIKey := sha256.Sum256([]byte(apiKey))
	// 		hashedKey := sha256.Sum256([]byte(key))

	// 		if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
	// 			return true, nil
	// 		}
	// 		return false, keyauth.ErrMissingOrMalformedAPIKey
	// 	},
	// })

	route.Get("/ping", s.controller.Ping)
}
