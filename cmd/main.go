package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	httpHandler "github.com/HolySxn/cv-hash-parser/internal/http"
	"github.com/HolySxn/cv-hash-parser/pkg/config"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	handler := httpHandler.NewHandler(logger)
	server := httpHandler.NewServer(handler)

	run(ctx, cfg, logger, server)
}

func run(
	ctx context.Context,
	cfg *config.Config,
	logger *slog.Logger,
	srv http.Handler,
) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler: srv,
	}

	go func() {
		logger.Info("server listening", "address", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error listening and serving", "error", err)
			cancel()
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutdownCancel()

		logger.Info("Gracefully shutting down...")

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("error shutting down server", "error", err)
		}
		logger.Info("HTTP server stopped.")
	}()
	wg.Wait()
	logger.Info("Shutdown complete.")
}
