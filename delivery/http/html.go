package http

import (
	"embed"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/hytaoist/faw-vw-auto/internal/log"
)

type html struct {
	tmpl *template.Template
}

func newHTML(assets embed.FS) *html {
	path := filepath.Join("assets", "static", "index.html")
	tmpl := template.Must(template.ParseFS(assets, path))
	return &html{tmpl}
}

func (html *html) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := html.tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		err = errors.WithStack(err)
		log.Error(err)
		http.Error(w, "Internal Server Error", 500)
	}
}
