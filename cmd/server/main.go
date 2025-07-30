package main

import (
	"net/http"
	"os"

	"github.com/lumiknit/can/internal/logger"
	"github.com/lumiknit/can/internal/server"
)

func main() {
	cfg := HandleFlags()

	// Setup structured logging
	logger := logger.New(cfg.ReleaseMode)

	r := server.SetupRoutes(logger, cfg.ReleaseMode)

	httpServer := &http.Server{
		Addr:    cfg.Addr(),
		Handler: r,
	}

	logger.Info("Server starting", "addr", httpServer.Addr, "release_mode", cfg.ReleaseMode)
	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
