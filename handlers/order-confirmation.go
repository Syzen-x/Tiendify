package handlers

import (
	"Tiendify/db"
	"Tiendify/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var orderTmpl = template.Must(template.New("order-confirmation.html").Funcs(template.FuncMap{
	"mul": func(a float64, b int) float64 {
		return a * float64(b)
	},
}).ParseFiles("views/order-confirmation.html"))

type OrderItem struct {
	models.Product
	Quantity int
}

func OrderConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}

	// Leer productos del carrito
	cartCookie, err := r.Cookie("cart")
	if err != nil || cartCookie.Value == "" {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}
	ids := strings.Split(cartCookie.Value, ",")

	// Calcular cantidades
	quantityMap := map[string]int{}
	for _, id := range ids {
		quantityMap[id]++
	}

	// IDs únicos
	uniqueIDs := []string{}
	for id := range quantityMap {
		uniqueIDs = append(uniqueIDs, id)
	}

	// Consulta SQL segura
	args := make([]interface{}, len(uniqueIDs))
	placeholders := make([]string, len(uniqueIDs))
	for i, id := range uniqueIDs {
		args[i] = id
		placeholders[i] = "?"
	}
	query := "SELECT id, name, description, price, stock, image_url FROM products WHERE id IN (" + strings.Join(placeholders, ",") + ")"

	dbConn, _ := db.Connect()
	defer dbConn.Close()

	rows, err := dbConn.Query(query, args...)
	if err != nil {
		http.Error(w, "Error cargando productos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []OrderItem
	var subtotal float64
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL); err == nil {
			qty := quantityMap[strconv.Itoa(p.ID)]
			p.ImageURL = "/static/images/" + p.ImageURL
			subtotal += p.Price * float64(qty)
			items = append(items, OrderItem{Product: p, Quantity: qty})
		}
	}

	shipping := 9.99
	tax := subtotal * 0.15
	discount := 10.00
	total := subtotal + shipping + tax - discount

	// Obtenemos el correo del usuario desde la cookie
	sessionCookie, err := r.Cookie("session")
	if err != nil || sessionCookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	email := sessionCookie.Value

	user, err := models.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "No se pudo obtener el usuario", http.StatusInternalServerError)
		return
	}

	// Limpiamos el carrito
	http.SetCookie(w, &http.Cookie{
		Name:   "cart",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Renderizamos
	orderTmpl.Execute(w, map[string]interface{}{
		"Items":     items,
		"Subtotal":  subtotal,
		"Envio":     shipping,
		"Tax":       tax,
		"Total":     total,
		"Direccion": user.Direccion,
	})
}
