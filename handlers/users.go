package handlers

import (
	"Tiendify/db"
	"Tiendify/models"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var usersTmpl = template.Must(template.ParseFiles("views/users.html"))

func UsersHandler(w http.ResponseWriter, r *http.Request) {

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

	// Obtener el usuario logueado
	conn, err := db.Connect()

	if err != nil {
		http.Error(w, "No se pudo obtener el usuario actual", http.StatusInternalServerError)
		return
	}

	// Conectar a la base de datos

	rows, err := conn.Query("SELECT id, firstName, lastName, email, role FROM users")

	if err != nil {
		http.Error(w, "Error conectando con la base de datos", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	defer rows.Close()
	var admins = 0
	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role)
		if err == nil {
			if u.Role == "admin" {
				admins++
			}
			users = append(users, u)
		}
	}

	// Renderizar template
	usersTmpl.Execute(w, map[string]interface{}{
		"User":   user,
		"Users":  users,
		"Admins": admins,
	})

}

var AddUserTmpl = template.Must(template.ParseFiles("views/addUsuario.html"))

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		AddUserTmpl.Execute(w, nil)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "No se pudieron procesar los datos", http.StatusBadRequest)
			return
		}

		firstName := strings.TrimSpace(r.FormValue("firstName"))
		lastName := strings.TrimSpace(r.FormValue("lastName"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")
		address := strings.TrimSpace(r.FormValue("address"))
		role := r.FormValue("role")

		// Validaciones básicas
		if firstName == "" || lastName == "" || email == "" || password == "" || address == "" || role == "" {
			http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
			return
		}

		if password != confirmPassword {
			http.Error(w, "Las contraseñas no coinciden", http.StatusBadRequest)
			return
		}

		if len(password) < 6 {
			http.Error(w, "La contraseña debe tener al menos 6 caracteres", http.StatusBadRequest)
			return
		}

		// Encriptar contraseña
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error al encriptar la contraseña", http.StatusInternalServerError)
			return
		}

		// Insertar en base de datos
		conn, err := db.Connect()
		if err != nil {
			http.Error(w, "Error con la base de datos", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		_, err = conn.Exec(`INSERT INTO users (firstName, lastName, email, password, direccion, role) VALUES (?, ?, ?, ?, ?, ?)`,
			firstName, lastName, email, hashedPassword, address, role)

		if err != nil {
			http.Error(w, "Error al guardar el usuario: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirigir al listado de usuarios
		http.Redirect(w, r, "/admin/usuarios", http.StatusSeeOther)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func DeleteUsuarioHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = dbConn.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error eliminando usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
