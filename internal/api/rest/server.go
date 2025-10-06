package rest

import (
	"context"
	"crypto-rates/internal/service"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func NewServer(port string, rateService *service.RateService) *Server {
	handlers := NewHandlers(rateService)

	mux := http.NewServeMux()

	// Раздельные маршруты
	mux.HandleFunc("/rates", handlers.GetRates) // Только /rates
	mux.HandleFunc("/rates/", handlers.GetRate) // /rates/что-угодно

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		server: server,
	}
}

func (s *Server) Start() error {
	zap.L().Info("Запуск REST API сервера", zap.String("address", s.server.Addr))

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("REST API server error: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	zap.L().Info("Остановка REST API сервера")
	return s.server.Shutdown(ctx)
}
