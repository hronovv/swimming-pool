package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RouteRegister interface {
	RegisterRoutes(r gin.IRouter)
}

type Server struct {
	server *http.Server
}

func New(addr string, handlers ...RouteRegister) *Server {
	engine := gin.Default()

	apiGroup := engine.Group("/api/v1")

	for _, handler := range handlers {
		handler.RegisterRoutes(apiGroup)
	}

	return &Server{
		server: &http.Server{
			Addr:         addr,
			Handler:      engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}
