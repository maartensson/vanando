package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func router(stateDir string) http.Handler {

	tmpl := template.Must(template.New("grid").Parse(tpl))
	tmpl2 := template.Must(template.New("grid").Parse(tpl2))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(filepath.Join(stateDir, "items.json"))
		if err != nil {
			http.Error(w, "missing data", http.StatusInternalServerError)
			return
		}
		var items []Item
		if err := json.NewDecoder(f).Decode(&items); err != nil {
			http.Error(w, "invalid data on disk", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, items)
	})

	mux.HandleFunc("GET /van/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		resp, err := http.Get(fmt.Sprintf("https://www.vannado.com/%s", id)) //"")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tmpl2.Execute(w, Content{
			Title:  getTitle(string(body)),
			Images: uniqueStrings(extractGalleryImages(string(body))),
		})
	})

	return mux
}

func getTitle(html string) string {
	re := regexp.MustCompile(`<div[^>]*fusion-title-3[^>]*>.*?<span[^>]*>([^<]+)</span>.*?</div>`)

	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimPrefix(match[1], "VANNADO ")
	} else {
		return "No title found"
	}
}

type Content struct {
	Title  string
	Images []Img
}

type Img struct {
	Link string
}

func uniqueStrings(input []string) []Img {
	seen := make(map[string]struct{})
	var result []Img

	for _, str := range input {
		if _, ok := seen[str]; !ok {
			seen[str] = struct{}{}
			result = append(result, Img{Link: str})
		}
	}

	return result
}

func extractGalleryImages(html string) []string {
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
