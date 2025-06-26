package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/login.html")

	if err != nil {
		log.Fatalf("error al cargar el template: %v", err)

	}
	tmpl.ExecuteTemplate(w, "base", nil)

}
