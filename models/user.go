package models

import (
	"Tiendify/db"
	"database/sql"
	"log"
)

type User struct {
	ID        int
	firstName string
	lastName  string
	Email     string
	Password  string
	Role      string
}

//funciones que permitan a traves de la base de datos todos los elementos de el crud

func GetUserByID(id int) (User, error) {
	var user User
	DB, err := db.Connect()
	if err != nil {
		log.Print("Error conectando a la base de datos:", err)
		return user, err
	}

	defer DB.Close()
	//consulta a la base de datos
	stmt, err := DB.Prepare("SELECT id, name, email, password, role FROM users WHERE id = ?")
	if err != nil {
		log.Print("Error preparando la consulta:", err)
		return user, err
	}

	row := stmt.QueryRow(id)
	err = row.Scan(&user.ID, &user.firstName, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No se encontro el usuario con el id: ", id)
			return user, nil
		}
		log.Print("Error al obtener usuario:", err)
	}
	log.Println("Usuario obtenido con exito: ", user)

	return user, nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	DB, err := db.Connect()
	if err != nil {
		log.Print("Error conectando a la base de datos:", err)
		return users, err
	}

	defer DB.Close()

	stmt, err := DB.Prepare("SELECT id, name, email, password, role FROM users")

	if err != nil {
		log.Print("Error preparando la consulta:", err)
		return users, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.firstName, &user.Email, &user.Password, &user.Role)
		if err != nil {
			log.Print("Error al obtener el usuario:", err)
			return users, err
		}
		users = append(users, user)
	}

	if err != nil {
		log.Print("Error al iterar sobre los resultados:", err)
		return users, err
	}

	log.Println("Usuarios obtenidos con exito: ", users)
	return users, nil

}
