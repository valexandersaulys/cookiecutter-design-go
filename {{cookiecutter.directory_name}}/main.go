package main

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed page.html.tmpl
var PageTemplateFS embed.FS

//go:embed static
var staticFilesFS embed.FS

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFS(PageTemplateFS, "*")
		if err != nil {
			panic(err)
		}
		template.ExecuteTemplate(w, "all", &struct {
			ProjectName string
		}{
			ProjectName: "{{cookiecutter.project_name}}",
		})
	})

	handleStaticFiles(r, "/static")

	panic(http.ListenAndServe("127.0.0.1:8081", r))
}

func handleStaticFiles(r chi.Router, path string) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	subFS, _ := fs.Sub(staticFilesFS, "static") // strip out the public/ from path
	fileServer := http.FileServer(http.FS(subFS))

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, fileServer)
		fs.ServeHTTP(w, r)
	})
}
