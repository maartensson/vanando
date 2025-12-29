package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/maartensson/ctxsrv/getenv"
	"github.com/maartensson/ctxsrv/jlog"
)

func main() {
	log := slog.New(jlog.New(nil))

	stateDir, err := getenv.StateDirectory()
	if err != nil {
		log.Error("failed to get stateDir", "error", err)
		os.Exit(1)
	}

	vans, err := getVans(context.Background(), log)
	if err != nil {
		log.Error("scraping failed", "error", err)
		os.Exit(1)
	}

	f, err := os.Create(filepath.Join(stateDir, "items.json"))
	if err != nil {
		log.Error("failed to open state file", "error", err)
		os.Exit(1)
	}

	if err := json.NewEncoder(f).Encode(vans); err != nil {
		log.Error("failed to save data to state", "error", err)
		os.Exit(1)
	}
}

type Item struct {
	TitleH5 string `json:"title_h_5"`
	TitleH2 string `json:"title_h_2"`
	Link    string `json:"link"`
	Image   string `json:"image"`
}

var re = regexp.MustCompile(`(?s)<h5[^>]*>.*?<strong>([^<]+)</strong>.*?</h5>.*?<h2[^>]*>.*?<span[^>]*>([^<]+)</span>.*?</h2>.*?<a[^>]+href="([^"]+)"[^>]*>.*?<img[^>]+src="(https://www\.vannado\.com/wp-content/uploads/[^"]+\.jpg)"[^>]*>`)

func getVans(ctx context.Context, log *slog.Logger) ([]Item, error) {

	start := time.Now()

	log.Info("Starting scrape - getting index")

	u := "https://www.vannado.com/van-conversion-inspirations/"
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call url: %w", err)
	}

	html, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	matches := re.FindAllStringSubmatch(string(html), -1)

	var items []Item

	log.Info("Index found", "van_count", len(matches), "elapsed", time.Since(start))

	for _, m := range matches {
		u, err := url.Parse(m[3])
		if err != nil {
			return nil, fmt.Errorf("failed to parse url: %w", err)
		}
		items = append(items, Item{
			TitleH5: m[1],
			TitleH2: m[2],
			Link:    fmt.Sprintf("/van/%s", path.Base(path.Clean(u.Path))),
			Image:   m[4],
		})
	}

	return items, nil
}
