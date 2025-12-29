package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"
)

type Item struct {
	TitleH5 string
	TitleH2 string
	Link    string
	Image   string
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
