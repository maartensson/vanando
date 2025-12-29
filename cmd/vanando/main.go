package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/maartensson/ctxsrv/jlog"
	"github.com/maartensson/ctxsrv/server"
)

func main() {
	log := slog.New(jlog.New(nil))

	ln, err := server.ActivationListener()
	if err != nil {
		log.Error("failed to get activation socket")
		os.Exit(1)
	}

	if err := server.HTTP(
		context.Background(),
		log,
		ln,
		router(),
		server.WithShutdownOnIdle(time.Hour),
	); err != nil {
		log.Error("HTTP server failed", "error", err)
	}
}
