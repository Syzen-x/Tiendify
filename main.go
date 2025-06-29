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
	db.Setup(conn) // Configura las tablas de la base de datos

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // para poder obtener las imagenes desde otras clases
	// la ruta de navegacion para el sistema pero /
	r.HandleFunc("/login", handlers.LoginHandler(conn)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
	r.HandleFunc("/register", handlers.RegisterPage).Methods("GET")
	r.HandleFunc("/register", handlers.RegisterHandler(conn)).Methods("POST")
	r.HandleFunc("/dashboard", handlers.Dashboard(conn)).Methods("GET")

	//ruta static para las imagenes

	r.HandleFunc("/logout", handlers.LogoutHandler)
	r.HandleFunc("/cart", handlers.CartHandler)
	r.HandleFunc("/add-to-cart", handlers.AddToCartHandler)
	r.HandleFunc("/remove-from-cart", handlers.RemoveFromCartHandler)
	r.HandleFunc("/clear-cart", handlers.ClearCartHandler)
	r.HandleFunc("/update-cart", handlers.UpdateCartHandler)
	r.HandleFunc("/order-confirmation", handlers.OrderConfirmationHandler)
	r.HandleFunc("/admin/productos/agregar", handlers.AddProductHandler)
	r.HandleFunc("/admin/productos", handlers.AdminProductsHandler)
	r.HandleFunc("/admin/productos/delete", handlers.DeleteProductHandler)
	r.HandleFunc("/admin/usuarios", handlers.UsersHandler)
	r.HandleFunc("/admin/usuarios/agregar", handlers.AddUserHandler)
	r.HandleFunc("/admin/usuarios/delete", handlers.DeleteUsuarioHandler)

	log.Println("Servidor de bd inicicalizado en el puerto 8080")

	// http/ListenAndServe inicia el servidor en el puerto 8080
	// nil es el manejador por defecto de mux.Router
	// las rutas con gorilla/mux
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

}
