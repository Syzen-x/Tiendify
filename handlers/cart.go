package handlers

import (
	"Tiendify/db"
	"Tiendify/models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var cartTmpl = template.Must(template.New("cart.html").Funcs(template.FuncMap{
	"mul": func(a float64, b int) float64 {
		return a * float64(b)
	},
}).ParseFiles("views/cart.html"))

type CartItem struct {
	models.Product
	Quantity int
}

func CartHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("cart")
	if err != nil || cookie.Value == "" {
		cartTmpl.Execute(w, map[string]interface{}{
			"Items": nil,
		})
		return
	}

	ids := strings.Split(cookie.Value, ",")

	// Mapa de cantidad por ID
	quantityMap := make(map[string]int)
	for _, id := range ids {
		quantityMap[id]++
	}

	// IDs únicos
	uniqueIDs := []string{}
	for id := range quantityMap {
		uniqueIDs = append(uniqueIDs, id)
	}

	// Consulta con placeholders
	query := "SELECT id, name, description, price, stock, image_url FROM products WHERE id IN (?" + strings.Repeat(",?", len(uniqueIDs)-1) + ")"
	args := make([]interface{}, len(uniqueIDs))
	for i, id := range uniqueIDs {
		args[i] = id
	}

	db, _ := db.Connect()
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Error al cargar el carrito", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cartItems []CartItem
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL)
		if err == nil {
			p.ImageURL = "/static/images/" + p.ImageURL
			qty := quantityMap[strconv.Itoa(p.ID)]
			cartItems = append(cartItems, CartItem{
				Product:  p,
				Quantity: qty,
			})
		}
	}
	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}

	envio := 9.99
	var subtotal float64
	subtotal = total / 1.15

	tax := subtotal * 0.15
	discount := 0.00 // si lo aplicás fijo

	total = subtotal + envio + tax - discount

	cartTmpl.Execute(w, map[string]interface{}{
		"Items":      cartItems,
		"TotalItems": len(ids),
		"Subtotal":   subtotal,
		"Envio":      envio,
		"Tax":        tax,
		"Discount":   discount,
		"Total":      total,
	})
}

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Obtenemos la cookie actual del carrito
	cookie, err := r.Cookie("cart")
	var cart []string
	if err == nil && cookie.Value != "" {
		cart = strings.Split(cookie.Value, ",")
	} else {
		cart = []string{} // inicializa sin nada
	}

	// Se agrega nuevo producto al carrito
	productID := r.FormValue("product_id")
	fmt.Println("Producto agregado al carrito:", productID)

	cart = append(cart, productID)

	// guardar la cookie actualizada
	http.SetCookie(w, &http.Cookie{
		Name:  "cart",
		Value: strings.Join(cart, ","),
		Path:  "/",
	})
	fmt.Println("Nuevo valor de la cookie:", strings.Join(cart, ","))

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
func RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}

	productID := r.FormValue("product_id")

	// Leer la cookie
	cookie, err := r.Cookie("cart")
	if err != nil {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}

	// Separar productos y filtrar el que se quiere eliminar
	var updated []string
	for _, id := range strings.Split(cookie.Value, ",") {
		if id != productID {
			updated = append(updated, id)
		}
	}

	// Actualizar cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "cart",
		Value: strings.Join(updated, ","),
		Path:  "/",
	})

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	// Se elimina la cookie que usamos para los productos del carrito de compras
	http.SetCookie(w, &http.Cookie{
		Name:   "cart",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func UpdateCartHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.FormValue("product_id")
	action := r.FormValue("action")

	if productID == "" || (action != "increase" && action != "decrease") {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("cart")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}

	ids := strings.Split(cookie.Value, ",")
	var updated []string
	count := 0

	for _, id := range ids {
		if id == productID {
			count++
		} else {
			updated = append(updated, id)
		}
	}

	// Ajustamos la cantidad
	switch action {
	case "increase":
		count++
	case "decrease":
		if count > 1 {
			count--
		}
	}

	// Agregamos el producto actualizado
	for i := 0; i < count; i++ {
		updated = append(updated, productID)
	}

	// Guardar cookie actualizada
	http.SetCookie(w, &http.Cookie{
		Name:  "cart",
		Value: strings.Join(updated, ","),
		Path:  "/",
	})

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}
