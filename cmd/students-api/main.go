package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imtiyazkumar/students-api/cmd/students-api/http/handlers/student"
	"github.com/imtiyazkumar/students-api/internal/config"
)

func main() {

	//load config
	cfg := config.MustLoad()

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	router.HandleFunc("POST /api/students", student.New())

	// Create the HTTP server
	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr, // e.g. ":8080"
		Handler: router,
	}

	slog.Info("Starting server...", slog.String("address", cfg.HTTPServer.Addr))

	//setup storage/database

	//gracefull shutdown
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	<-quit

	slog.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Server forced to shutdown:", slog.String("error", err.Error()))
	}

	slog.Info("Server exit sucessfully")
}
