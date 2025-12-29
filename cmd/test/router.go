package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func router(items []Item) http.Handler {

	tmpl := template.Must(template.New("grid").Parse(tpl))
	tmpl2 := template.Must(template.New("grid").Parse(tpl2))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
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
			Images: uniqueStrings(ExtractGalleryImages(string(body))),
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
