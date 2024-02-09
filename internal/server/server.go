package server

import (
	"github.com/TimothyYe/godns/internal/server/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	addr       string
	username   string
	password   string
	app        *fiber.App
	controller *controllers.Controller
}

func (s *Server) SetAddress(addr string) *Server {
	s.addr = addr
	return s
}

func (s *Server) SetAuthInfo(username, password string) *Server {
	s.username = username
	s.password = password
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

func (s *Server) Stop() error {
	return s.app.Shutdown()
}

func (s *Server) initRoutes() {
	// set cross domain access rules
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST ,PUT ,DELETE, OPTIONS",
		AllowHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}))

	s.app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			s.username: s.password,
		},
	}))

	// Create routes group.
	route := s.app.Group("/api/v1")
	route.Get("/auth", s.controller.Auth)
	route.Get("/info", s.controller.GetBasicInfo)
}
