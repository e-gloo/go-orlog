package main

import (
	"context"
	"flag"
	"log/slog"

	"github.com/e-gloo/orlog/internal/pkg/logging"
	"github.com/e-gloo/orlog/internal/server"
)

func main() {
	// Server run config
	addr := flag.String("addr", "localhost", "http service address")
	port := flag.String("port", "8080", "http service port")
	dev := flag.Bool("dev", false, "Running in development mode")
	flag.Parse()

	ctx := context.Background()
	logging.InitLogger(*dev, nil, nil)

	// Create and run server
	srv := server.NewServer(ctx, *addr, *port)
	slog.Info("Listening", "addr", *addr, "port", *port)
	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("Error running server", "err", err)
	}
}
