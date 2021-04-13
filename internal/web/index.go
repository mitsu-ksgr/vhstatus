package web

import (
	"html/template"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles(getTemplateDirPath() + "index.html")
	if err != nil {
		log.Print(err)
		render500(w)
		return
	}

	if err := temp.Execute(w, getVHStatusParams()); err != nil {
		log.Print(err)
		render500(w)
	}
}
