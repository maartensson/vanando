package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/maartensson/ctxsrv/getenv"
	"github.com/maartensson/ctxsrv/jlog"
	"github.com/maartensson/ctxsrv/server"
)

func ExtractGalleryImages(html string) []string {
	var urls []string

	// Match each <div class="fusion-gallery-image"> block
	divRe := regexp.MustCompile(`(?s)<div class="fusion-gallery-image">(.*?)</div>\s*</div>`)
	divMatches := divRe.FindAllStringSubmatch(html, -1)

	// Regexes for data-orig-src and background-image
	origRe := regexp.MustCompile(`data-orig-src="(https?:\/\/www\.vannado\.com\/wp-content\/uploads\/[^\s"']+\.jpg)"`)
	bgRe := regexp.MustCompile(`background-image:\s*url\(&quot;(https?:\/\/www\.vannado\.com\/wp-content\/uploads\/[^\s"']+\.jpg)&quot;\)`)

	for _, div := range divMatches {
		content := div[1]
		var img string

		if match := origRe.FindStringSubmatch(content); match != nil {
			img = match[1] // prefer data-orig-src
		} else if match := bgRe.FindStringSubmatch(content); match != nil {
			img = match[1] // fallback to background-image
		}

		if img != "" {
			urls = append(urls, img)
		}
	}

	return urls
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	items, err := getVans(ctx, log)
	if err != nil {
		log.Error("failed to get vans", "error", err)
	}

	if err := server.HTTP(
		ctx, cancel, router(items),
		log, ln, time.Hour,
		func(ctx context.Context, c net.Conn) context.Context {
			return ctx
		},
	); err != nil {
		log.Error("HTTP server failed", "error", err)
	}
}
