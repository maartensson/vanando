package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
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

func getVans() ([]Item, error) {
	r, err := http.Get("https://www.vannado.com/van-conversion-inspirations/")
	if err != nil {
		return nil, fmt.Errorf("failed to call url: %w", err)
	}

	html, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	re := regexp.MustCompile(`(?s)<h5[^>]*>.*?<strong>([^<]+)</strong>.*?</h5>.*?<h2[^>]*>.*?<span[^>]*>([^<]+)</span>.*?</h2>.*?<a[^>]+href="([^"]+)"[^>]*>.*?<img[^>]+src="(https://www\.vannado\.com/wp-content/uploads/[^"]+\.jpg)"[^>]*>`)

	matches := re.FindAllStringSubmatch(string(html), -1)

	var items []Item

	for _, m := range matches {
		u, err := url.Parse(m[3])
		if err != nil {
			log.Fatal(err)
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

type Item struct {
	TitleH5 string
	TitleH2 string
	Link    string
	Image   string
}

func main() {
	// Example data, in practice fill this with your regex extraction
	items, err := getVans()
	if err != nil {
		log.Fatal(err)
	}
	const tpl = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<script src="https://cdn.tailwindcss.com"></script>
<title>Image Grid</title>
</head>
<body class="bg-gray-900 text-white p-4">

<h1 class="text-3xl font-bold mb-6">Gallery</h1>

<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
{{range .}}
  <div class="bg-gray-800 rounded-lg overflow-hidden shadow-lg">
    <a href="{{.Link}}">
      <img src="{{.Image}}" alt="{{.TitleH2}}" class="w-full h-64 object-cover hover:scale-105 transition-transform duration-300">
    </a>
    <div class="p-4">
      <h5 class="text-sm font-semibold text-gray-300">{{.TitleH5}}</h5>
      <h2 class="text-lg font-bold">{{.TitleH2}}</h2>
    </div>
  </div>
{{end}}
</div>

</body>
</html>
`
	const tpl2 = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<script src="https://cdn.tailwindcss.com"></script>
<title>{{.Title}}</title>
</head>
<body class="bg-gray-900 text-white p-4">

<h1 class="text-3xl font-bold mb-6">{{.Title}}</h1>

<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
{{range .Images}}
  <div class="bg-gray-800 rounded-lg overflow-hidden shadow-lg">
    <a href="{{.Link}}" download>
      <img src="{{.Link}}" class="w-full h-64 object-cover hover:scale-105 transition-transform duration-300">
    </a>
  </div>
{{end}}
</div>

</body>
</html>
`

	tmpl := template.Must(template.New("grid").Parse(tpl))
	tmpl2 := template.Must(template.New("grid").Parse(tpl2))
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, items)
		if err != nil {
			log.Fatal(err)
		}
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
		if err := tmpl2.Execute(w, Content{
			Title:  GetTitle(string(body)),
			Images: uniqueStrings(ExtractGalleryImages(string(body))),
		}); err != nil {
			log.Fatal(err)
		}
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
}

func GetTitle(html string) string {
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
