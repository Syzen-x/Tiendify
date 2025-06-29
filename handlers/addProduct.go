package handlers

import (
	"Tiendify/db"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var addProductTmpl = template.Must(template.ParseFiles("views/addProduct.html"))

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		addProductTmpl.Execute(w, nil)
		return
	}
	_, handler, err := r.FormFile("productImages")
	if err != nil {
		log.Println("No se pudo obtener la imagen:", err)
		http.Error(w, "Error subiendo la imagen", http.StatusBadRequest)
		return
	}
	imageName := filepath.Base(handler.Filename)
	os.MkdirAll("./static/images", os.ModePerm) //crear la carpeta, si existe no pasa nada y si no la crea
	os.Create("./static/images/" + imageName)

	if r.Method == http.MethodPost {

		name := r.FormValue("productName")
		description := r.FormValue("productDescription")
		price, _ := strconv.ParseFloat(r.FormValue("productPrice"), 64)
		stock, _ := strconv.Atoi(r.FormValue("productQuantity"))

		// Manejar imagen
		file, handler, err := r.FormFile("productImages")
		if err != nil {
			http.Error(w, "Error subiendo la imagen", http.StatusBadRequest)
			return
		}
		defer file.Close()

		imageName := filepath.Base(handler.Filename)
		dst, err := os.Create("./static/images/" + imageName)
		if err != nil {
			http.Error(w, "No se pudo guardar la imagen", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		io.Copy(dst, file)

		// Guardar en la base de datos
		conn, err := db.Connect()
		if err != nil {
			http.Error(w, "Error con la base de datos", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		_, err = conn.Exec("INSERT INTO products (name, description, price, stock, image_url) VALUES (?, ?, ?, ?, ?)",
			name, description, price, stock, imageName)
		if err != nil {
			http.Error(w, "Error insertando el producto", http.StatusInternalServerError)
			return
		}

		// Redirigir
		http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)
	}
}
func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID faltante", http.StatusBadRequest)
		return
	}

	dbConn, err := db.Connect()
	if err != nil {
		http.Error(w, "Error de conexión DB", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	_, err = dbConn.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error eliminando producto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
