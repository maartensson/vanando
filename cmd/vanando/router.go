package main

import (
	"io"
	"net/http"
	"regexp"
)

func router() http.Handler {
	re := regexp.MustCompile(`https://www\.vannado\.com[^\s"']*kopia\.jpg`)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	})

	mux.HandleFunc("POST /{$}", func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue("url")

		req, err := http.NewRequestWithContext(r.Context(), "GET", url, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var images []string
		matches := re.FindAll(body, -1)
		for _, match := range matches {
			images = append(images, string(match))
		}

		tpl.Execute(w, uniqueStrings(images))
	})

	return mux
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, str := range input {
		if _, ok := seen[str]; !ok {
			seen[str] = struct{}{}
			result = append(result, str)
		}
	}

	return result
}
