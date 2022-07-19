package server

import (
	"AlexSarva/gofermart/handlers"
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *models.Config, database *app.Database) *Server {

	handler := handlers.MyHandler(database)
	server := http.Server{
		Addr:    cfg.ServerAddress,
		Handler: handler,
	}
	return &Server{
		httpServer: &server,
	}
}

func (a *Server) Run() error {
	addr := a.httpServer.Addr
	log.Printf("Web-server started at http://%s", addr)
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
