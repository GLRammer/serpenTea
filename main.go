package main

import (
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Title string
	MOTD  string
	Style string
}

var templates map[string]*template.Template

var liveDev = true

func initTemplates() {
	templates = make(map[string]*template.Template)
	pages := []string{
		"home.html",
		"404.html",
	}
	for _, page := range pages {
		t, err := template.ParseFiles(
			"templates/layout.html",
			"templates/"+page,
		)
		if err != nil {
			log.Fatal(err)
		}
		templates[page] = t
	}
}

func render(w http.ResponseWriter, name string, data PageData) {
	var t *template.Template
	var err error

	if liveDev {
		files := []string{
			"templates/layout.html",
			"templates/" + name,
		}

		t, err = template.ParseFiles(files...)
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	} else {
		var ok bool
		t, ok = templates[name]
		if !ok {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Render error", http.StatusInternalServerError)
		log.Println(err)
	}
}

func nfHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := PageData{
		Title: "404",
		Style: "404.css",
	}

	render(w, "404.html", data)
}

func homeHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		nfHandle(w, r)
		return
	}
	data := PageData{
		Title: "Home",
		MOTD:  "Hello from a template-based Go web server!",
	}

	render(w, "home.html", data)
}

func killHandle(srv *http.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Shutting down server...\n"))

		go func() {
			// shutdown must run in a goroutine so response can complete
			if err := srv.Shutdown(r.Context()); err != nil {
				log.Println("Shutdown error:", err)
			}
		}()
	}
}

func main() {

	if !liveDev {
		initTemplates()
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.HandleFunc("GET /", homeHandle)
	mux.HandleFunc("GET /killer", killHandle(srv))

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	log.Println("Server running at http://localhost:8080")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
