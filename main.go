package main

import (
	"Tiendify/db"
	"Tiendify/handlers"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var conn *sql.DB

func main() {
	var err error
	conn, err = db.Connect()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	r := mux.NewRouter()
	// la ruta de navegacion para el sistema pero /
	r.HandleFunc("/login", handlers.LoginHandler(conn)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
	r.HandleFunc("/register", handlers.RegisterPage).Methods("GET")           // GET
	r.HandleFunc("/register", handlers.RegisterHandler(conn)).Methods("POST") // POST

	log.Println("Servidor de bd inicicalizado en el puerto 8080")

	// http/ListenAndServe inicia el servidor en el puerto 8080
	// nil es el manejador por defecto de mux.Router
	// las rutas con gorilla/mux
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

}
