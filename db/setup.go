package db

import (
	"database/sql"
	"log"
)

func Setup(db *sql.DB) {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		firstName TEXT NOT NULL,
		lastName TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		direccion TEXT
	);`

	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		stock INTEGER NOT NULL,
		image_url TEXT
	);`

	if _, err := db.Exec(createUsersTable); err != nil {
		log.Fatalf("❌ Error creating users table: %v", err)
	}
	if _, err := db.Exec(createProductsTable); err != nil {
		log.Fatalf("❌ Error creating products table: %v", err)
	}

	log.Println("✅ Tablas creadas o ya existentes.")

}
