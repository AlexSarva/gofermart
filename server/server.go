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

// Server implementation of custom server
type Server struct {
	httpServer *http.Server
}

// NewServer Initializing new server instance
func NewServer(cfg *models.Config, database *app.Database) *Server {

	handler := handlers.MyHandler(database)
	server := http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      handler,
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
	}
	return &Server{
		httpServer: &server,
	}
}

// Run method that starts the server
func (a *Server) Run() error {
	addr := a.httpServer.Addr
	log.Printf("Web-server started at http://%s", addr)
	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil {
			log.Printf("Something wrong with server: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
