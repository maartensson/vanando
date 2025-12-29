package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/maartensson/ctxsrv/getenv"
	"github.com/maartensson/ctxsrv/jlog"
	"github.com/maartensson/ctxsrv/server"
)

type Item struct {
	TitleH5 string `json:"title_h_5"`
	TitleH2 string `json:"title_h_2"`
	Link    string `json:"link"`
	Image   string `json:"image"`
}

func main() {
	log := slog.New(jlog.New(nil))

	ln, err := server.ActivationListener()
	if err != nil {
		log.Error("failed to get activation socket")
		os.Exit(1)
	}

	stateDir, err := getenv.StateDirectory()
	if err != nil {
		log.Error("failed to get state dir", "error", err)
		os.Exit(1)
	}

	if err := server.HTTP(
		context.Background(),
		log,
		ln,
		router(stateDir),
		server.WithShutdownOnIdle(time.Hour),
	); err != nil {
		log.Error("HTTP server failed", "error", err)
	}
}
