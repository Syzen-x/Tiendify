package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("views/login.html"))

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, map[string]interface{}{
		"ErrorMessage": "",
	})
}

func LoginHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		email := r.FormValue("email")
		fmt.Println("Email:", email)
		password := r.FormValue("password")
		fmt.Println("Contraseña:", password)

		var dbPassword string
		err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&dbPassword)

		if err != nil {
			fmt.Print(err.Error())
			tmpl.Execute(w, map[string]interface{}{
				"ErrorMessage": "Email or Password is incorrect",
			})
			return
		}

		if password != dbPassword {
			tmpl.Execute(w, map[string]interface{}{
				"ErrorMessage": "Email or Password is incorrect",
			})

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: email,
			Path:  "/",
		})
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

		fmt.Fprintf(w, "¡Bienvenido, %s!", email)
	}
}
