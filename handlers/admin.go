package handlers

import (
	"Tiendify/db"
	"Tiendify/models"
	"html/template"
	"net/http"
)

var adminProductsTmpl = template.Must(template.ParseFiles("views/products.html"))

func AdminProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar si tiene la session activa
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := models.GetUserByEmail(cookie.Value)

	if err != nil || user.ID == 0 || user.Role != "admin" {
		http.Error(w, "Acceso no autorizado", http.StatusForbidden)
		return
	}

	// Cargar productos realesa
	dbConn, err := db.Connect()
	if err != nil {
		http.Error(w, "Error de conexión", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("SELECT id, name, description, price, stock, image_url FROM products")
	if err != nil {
		http.Error(w, "Error cargando productos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL)

		if err == nil {
			p.ImageURL = "/static/images/" + p.ImageURL

			products = append(products, p)
		}
	}

	// Renderizar
	adminProductsTmpl.Execute(w, map[string]interface{}{
		"User":     user,
		"Products": products,
	})
}
