package main

import (
	"context"
	"log/slog"
	"net"
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

	port, err := getenv.NetworkPort()
	if err != nil {
		log.Error("failed to get valid port")
		os.Exit(1)
	}

	ln, err := server.ActivationListener(port, false)
	if err != nil {
		log.Error("failed to get activation socket")
		os.Exit(1)
	}

	stateDir, err := getenv.StateDirectory()
	if err != nil {
		log.Error("failed to get state dir", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := server.HTTP(
		ctx, cancel, router(stateDir),
		log, ln, time.Hour,
		func(ctx context.Context, c net.Conn) context.Context {
			return ctx
		},
	); err != nil {
		log.Error("HTTP server failed", "error", err)
	}
}
