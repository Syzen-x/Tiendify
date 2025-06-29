package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

var registerTmpl = template.Must(template.ParseFiles("views/register.html"))

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	registerTmpl.Execute(w, map[string]interface{}{
		"ErrorMessage": " ",
	})
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		direccion := r.FormValue("direccion")

		fmt.Print("hola2")
		// Validación básica
		if firstName == "" || lastName == "" || email == "" || password == "" {
			registerTmpl.Execute(w, map[string]interface{}{
				"ErrorMessage": "All fields are required.",
			})
			return
		}
		fmt.Print("hola1")
		// Verificar si el email ya existe
		var existing string
		err := db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existing)
		if err == nil {
			registerTmpl.Execute(w, map[string]interface{}{
				"ErrorMessage": "Email already registered.",
			})
			return
		}

		// Insertar nuevo usuario
		fmt.Print("hola")
		_, err = db.Exec("INSERT INTO users (firstName, lastName, email, password, role, direccion) VALUES (?, ?, ?, ?, ?,?)", firstName, lastName, email, password, "user", direccion)
		if err != nil {
			fmt.Println(err.Error())
			registerTmpl.Execute(w, map[string]interface{}{
				"ErrorMessage": "Something went wrong. Try again.",
			})
			return
		}

		// Redirigir al login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
