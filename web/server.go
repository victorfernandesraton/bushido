package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed pages
var webFS embed.FS

//go:embed assets
var assetsFS embed.FS

func NewWebServer(port int) *http.Server {

	assetsFS, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`{"status": "OK"}`))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title string
		}{
			"Teste",
		}
		err := template.
			Must(template.New("").ParseFS(webFS, "pages/*")).
			ExecuteTemplate(w, "index.tmpl.html", data)
		if err != nil {
			log.Fatal(err)
		}
	})
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS))))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}

	return server

}
