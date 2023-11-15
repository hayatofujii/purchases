package server

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	controllers "haf.systems/purchases/controllers"
	"haf.systems/purchases/routes"
)

type Server struct {
	port   interface{}
	server *gin.Engine
}

func NewServer(c *controllers.Controllers, port interface{}) *Server {
	engine := gin.New()

	engine.Use()

	routes.SetRoutes(c, engine)

	return &Server{
		port:   port,
		server: engine,
	}
}

func (s *Server) GetServer() *gin.Engine {
	return s.server
}

func (s *Server) CorsSetup() {
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("*")

	s.server.Use(cors.New(corsConfig))
}

func (s *Server) Run() {

	serverAddr := ":"

	switch v := s.port.(type) {

	case string:
		if v != "" {
			serverAddr = serverAddr + v
		} else {
			serverAddr = serverAddr + "8080"
		}

	case int:
		serverAddr = serverAddr + fmt.Sprint(v)

	default:
		fmt.Println(fmt.Errorf("port not specified as string or int"))
		fmt.Println("using port 8080")

		serverAddr = serverAddr + "8080"
	}

	s.server.Run(serverAddr)
}
