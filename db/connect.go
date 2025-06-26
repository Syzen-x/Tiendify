package db

// Conectamos ala base de datos MySQL
import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Importamos el driver de MySQL
	"github.com/joho/godotenv"
)

func Connect() (*sql.DB, error) {
	// cargamos los datos del archivo .env

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// root:root@tcp(localhost:3306)/tiendify
	dns := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println(dns)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	//probamos la conexion haciendo un ping ala base de datos
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Conectado a la base de datos MySQL")

	return db, nil

}
