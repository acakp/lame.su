package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/acakp/dumbchat/chat"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func renderHandler(chatApp chat.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatApp.Execute(w)
	}
}

func parseIndex() template.Template {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println("error parsing internal templates:", err)
		return *template.New("error parsing tmpl")
	}
	return *tmpl
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.Handle("/static/*", fs)

	db, err := sql.Open("sqlite3", "chat.db")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)

	chatApp, err := chat.New(chat.Config{
		DB:       db,
		BasePath: os.Getenv("CHAT_BASE_PATH"),
	})
	if err != nil {
		log.Fatal(err)
	}

	chatApp.CreateTables(db)

	indexTmpl := parseIndex()
	chatApp.AttachTemplates(&indexTmpl)

	r.Route("/chat", func(r chi.Router) {
		chatApp.RegisterRoutes(r)
	})
	r.Get("/", renderHandler(*chatApp))
	fmt.Println("starting TEST SERVER on :8880...")
	http.ListenAndServe(":8880", r)
}
