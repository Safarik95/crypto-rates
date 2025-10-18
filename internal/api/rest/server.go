package rest

import (
	"context"
	_ "crypto-rates/docs"
	"crypto-rates/internal/config"
	"crypto-rates/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	router *gin.Engine
	server *http.Server
}

func NewServer(port string, rateService *service.RateService) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	handlers := NewHandlers(rateService)

	router.GET("/rates", handlers.GetRates)
	router.GET("/rates/:cryptocurrency", handlers.GetRate)
	router.GET("/health", handlers.HealthCheck)
	router.GET("/swagger.json", handlers.SwaggerJSON)

	return &Server{
		router: router,
		server: &http.Server{
			Addr:         ":" + port,
			Handler:      router,
			ReadTimeout:  config.ServerReadTimeout,
			WriteTimeout: config.ServerWriteTimeout,
			IdleTimeout:  config.ServerIdleTimeout,
		},
	}
}

func (s *Server) Start() error {
	zap.L().Info("starting REST API server", zap.String("address", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	zap.L().Info("stopping REST API server")
	return s.server.Shutdown(ctx)
}
