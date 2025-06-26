package handlers

import (
	"Tiendify/models"
	"database/sql"
	"html/template"
	"net/http"
)

var dashboardTmpl = template.Must(template.ParseFiles("views/dashboard.html"))

func Dashboard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verifica si hay cookie
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Cargar productos
		rows, err := db.Query("SELECT id, name, description, price, stock, image_url FROM products")
		if err != nil {

			http.Error(w, "Error cargando los productos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type Product struct {
			ID          int
			Name        string
			Description string
			Price       float64
			Stock       int
		}

		var products []models.Product
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL); err != nil {
				http.Error(w, "Error obteniendo el producto ", http.StatusInternalServerError)
				return
			}

			p.ImageURL = "static/images/" + p.ImageURL
			products = append(products, p)
		}

		var firstName string
		err = db.QueryRow("SELECT firstName FROM users WHERE email = ?", cookie.Value).Scan(&firstName)
		if err != nil {
			http.Error(w, "No se pudo obtener el nombre del usuario", http.StatusInternalServerError)
			return
		}
		// Mostrar la dashboard
		dashboardTmpl.Execute(w, map[string]interface{}{
			"UserEmail": cookie.Value,
			"firstName": firstName,
			"Products":  products,
		})
	}
}
