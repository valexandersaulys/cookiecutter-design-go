package main

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

//go:embed page.html.tmpl
var PageTemplateFS embed.FS

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

	panic(http.ListenAndServe("127.0.0.1:8081", r))
}
