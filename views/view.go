package views

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func baseFiles(viewsDir string) []string {
	files, _ := filepath.Glob(viewsDir + "/*.gohtml")

	return files
}

func New(viewsDir string, layout string, files ...string) *View {
	files = append(baseFiles(viewsDir), files...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Print("error parsing files: ", err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}
