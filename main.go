package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gergenus/GoMockServer/src/config"
	"github.com/Gergenus/GoMockServer/src/handler"
)

func main() {
	cfg, err := config.LoadConfig("./conf.yaml")
	if err != nil {
		panic(err)
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	s := handler.NewServer(cfg, log)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: http.HandlerFunc(s.HandleRequests),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info("Server listening", slog.String("address", srv.Addr))
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Error("Server error", slog.String("error", err.Error()))
		}

	case sig := <-shutdown:
		log.Info("Shutdown signal received", slog.Any("signal", sig))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("Graceful shutdown failed", slog.String("error", err.Error()))
			srv.Close()
		}
	}
}
