package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed public
var staticFiles embed.FS

var staticDir = "public"

func rootPath(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = fmt.Sprintf("/%s/the.html", staticDir)
		} else {
			b := strings.Split(r.URL.Path, "/")[0]
			if b != staticDir {
				r.URL.Path = fmt.Sprintf("/%s%s", staticDir, r.URL.Path)
			}
		}
		log.Printf("Serving %s\n", r.URL.Path)
		h.ServeHTTP(w, r) 
	})
}

func main() {

	var staticFS = http.FS(staticFiles)
	fs := rootPath(http.FileServer(staticFS))

	// Serve static files
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Listening on :%s...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
