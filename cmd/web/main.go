package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	hugoPublicDir := "./content"

	// check if hugo public directory exists
	if _, err := os.Stat(hugoPublicDir); os.IsNotExist(err) {
		log.Printf("warning: hugo public directory '%s' not found. running without static files.", hugoPublicDir)
	}

	// serve hugo static files
	fileServer(r, hugoPublicDir)

	// todo: add chat routes here later
	// r.mount("/chat", chat.routes())

	log.Printf("starting on :8888...")
	if err := http.ListenAndServe(":8888", r); err != nil {
		log.Fatal(err)
	}
}

// fileserver sets up a static file server with proper 404 handling
func fileServer(r chi.Router, publicDir string) {
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		return
	}

	fs := http.FileServer(http.Dir(publicDir))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// check if file exists
		requestedFile := filepath.Join(publicDir, r.URL.Path)
		if _, err := os.Stat(requestedFile); os.IsNotExist(err) {
			// file doesn't exist - try serving index.html
			if r.URL.Path != "/" && filepath.Ext(requestedFile) == "" {
				http.ServeFile(w, r, filepath.Join(publicDir, "index.html"))
				return
			}
		}
		fs.ServeHTTP(w, r)
	})
}
