//errores: go
//podemos imaginar donde puede ocurrrir el error y controlandolo
//manejar de mejor manera

package main

import (
	"fmt"
	"time"
)

//estructura de la clase Usuario

type Usuario struct {
	id        int
	nombres   string
	apellidos string
	edad      int
	email     string
}

//implementar la encapsulacion de la clase Usuario
// proteger la clase y dar a traves de metodos el acceso a los datos necesarios

func (u *Usuario) GetId() int {
	return u.id
}

func (u *Usuario) SetId(id int) {
	u.id = id
}

func (u *Usuario) GetNombres() string {
	return u.nombres
}
func (u *Usuario) SetNombres(nombres string) {
	u.nombres = nombres
}
func (u *Usuario) GetApellidos() string {
	return u.apellidos
}
func (u *Usuario) SetApellidos(apellidos string) {
	u.apellidos = apellidos
}
func (u *Usuario) GetEdad() int {
	return u.edad

}
func (u *Usuario) SetEdad(edad int) {
	u.edad = edad
}
func (u *Usuario) GetEmail() string {
	return u.email
}
func (u *Usuario) SetEmail(email string) {
	u.email = email
}

//una estructura para almacenar libros

type Libro struct {
	id       int
	titulo   string
	autor    string
	isbn     string
	prestado bool
}

// funciones para acceder a los datos de la clase Libro
func (l *Libro) GetId() int {
	return l.id
}

func (l *Libro) SetId(id int) {

	l.id = id
}

func (l *Libro) GetTitulo() string {
	return l.titulo
}
func (l *Libro) SetTitulo(titulo string) {
	l.titulo = titulo
}
func (l *Libro) GetAutor() string {
	return l.autor
}

func (l *Libro) SetAutor(autor string) {
	l.autor = autor
}

func (l *Libro) GetIsbn() string {

	return l.isbn
}

func (l *Libro) SetIsbn(isbn string) {
	l.isbn = isbn
}

func (l *Libro) EstaPrestado() bool {
	return l.prestado
}
func (l *Libro) SetPrestado(prestado bool) {
	l.prestado = prestado
}

// una estructura para almacenar prestamo

type Prestamo struct {
	id              int
	usuario         *Usuario
	libros          []*Libro
	fechaPrestamo   time.Time
	fechaDevolucion time.Time
}

// funciones para acceder a los datos de la clase Prestamo
func (p *Prestamo) GetId() int {
	return p.id
}
func (p *Prestamo) SetId(id int) {
	p.id = id
}
func (p *Prestamo) GetUsuario() *Usuario {
	return p.usuario
}
func (p *Prestamo) SetUsuario(usuario *Usuario) {
	p.usuario = usuario
}
func (p *Prestamo) GetLibros() []*Libro {
	return p.libros
}

func (p *Prestamo) SetLibros(libros []*Libro) {
	p.libros = libros
}
func (p *Prestamo) GetFechaPrestamo() time.Time {
	return p.fechaPrestamo
}
func (p *Prestamo) SetFechaPrestamo(fecha time.Time) {
	p.fechaPrestamo = fecha
}
func (p *Prestamo) GetFechaDevolucion() time.Time {
	return p.fechaDevolucion
}
func (p *Prestamo) SetFechaDevolucion(fecha time.Time) {
	p.fechaDevolucion = fecha
}

//una clase que permita gestionar todo el recurso de bioblioteca

type Biblioteca struct {
	usuarios  []*Usuario
	libros    []*Libro
	prestamos []*Prestamo
	lastId    int
}

// los metodos que permiten gestionar los recursos de la biblioteca
// el manejo de los errores

func (b *Biblioteca) AgregarUsuario(nombres, apellidos, email string, edad int) error {
	//verificar si las variables son validas
	if nombres == "" || apellidos == "" || email == "" || edad <= 0 {
		return fmt.Errorf("Datos del usuario invalidos")
	}

	b.lastId++
	usuario := &Usuario{
		nombres:   nombres,
		apellidos: apellidos,
		email:     email,
		edad:      edad,
	}

	//asignamos la id
	usuario.SetId(b.lastId)
	b.usuarios = append(b.usuarios, usuario)
	return nil
}

func (b *Biblioteca) AgregarLibro(titulo, autor, isbn string) error {
	if titulo == "" || autor == "" || isbn == "" {
		return fmt.Errorf("Datos del libro invalidos")
	}

	b.lastId++
	libro := &Libro{
		titulo: titulo,
		autor:  autor,
		isbn:   isbn,
	}

	libro.SetId(b.lastId)
	b.libros = append(b.libros, libro)
	return nil
}

// metodo para prestar un libro y registrar el prestamo
func (b *Biblioteca) PrestarLibro(usuarioId int, librosIds ...int) error {
	var usuario *Usuario
	var libros []*Libro
	for _, u := range b.usuarios {
		if u.GetId() == usuarioId {
			usuario = u
			break
		}
	}

	if usuario == nil {
		return fmt.Errorf("Usuario no encontrado")
	}
	// buscar los libros para prestar
	for _, libroId := range librosIds {
		var libro *Libro
		for _, l := range b.libros {
			if l.GetId() == libroId {
				libro = l
				break
			}
		}

		if libro == nil {
			return fmt.Errorf("Libro con ID %d no encontrado", libroId)
		}

		if libro.EstaPrestado() {
			return fmt.Errorf("El libro '%s' ya esta prestado", libro.GetTitulo())
		}

		libro.SetPrestado(true)
		libros = append(libros, libro)
	}
	b.lastId++
	prestamo := &Prestamo{
		id:            b.lastId,
		usuario:       usuario,
		libros:        libros,
		fechaPrestamo: time.Now(),
	}

	for _, libro := range libros {

		libro.SetPrestado(true)
	}
	//agregar el prestamo a la lista de prestamos de la biblioteca
	b.prestamos = append(b.prestamos, prestamo)
	return nil
}

func (b *Biblioteca) DevolverLibro(prestamoId int) error {
	var prestamo *Prestamo
	for _, p := range b.prestamos {
		if p.GetId() == prestamoId {
			prestamo = p
			break
		}
	}

	if prestamo == nil {
		return fmt.Errorf("Prestamo con ID %d no encontrado", prestamoId)
	}

	for _, libro := range prestamo.GetLibros() {
		libro.SetPrestado(false)
	}

	prestamo.SetFechaDevolucion(time.Now())
	return nil
}

func main() {
	biblioteca := &Biblioteca{}
	// menu
	fmt.Println("=== Bienvenido a la Biblioteca ===")
	for {
		fmt.Println("1. Agregar Usuario")
		fmt.Println("2. Agregar Libro")
		fmt.Println("3. Prestar Libro")
		fmt.Println("4. Devolver Libro")
		fmt.Println("5. Listar Usuarios")
		fmt.Println("6. Listar Libros")
		fmt.Println("7. Salir")
		var opcion int
		fmt.Print("Seleccione una opcion: ")
		fmt.Scan(&opcion)
		switch opcion {
		case 1:

			var nombres, apellidos, email string
			var edad int
			fmt.Print("Ingrese los nombres del usuario: ")
			fmt.Scan(&nombres)
			fmt.Print("Ingrese los apellidos del usuario: ")
			fmt.Scan(&apellidos)
			fmt.Print("Ingrese el email del usuario: ")
			fmt.Scan(&email)
			fmt.Print("Ingrese la edad del usuario: ")
			fmt.Scan(&edad)
			err := biblioteca.AgregarUsuario(nombres, apellidos, email, edad)
			if err != nil {
				fmt.Println("Error al agregar usuario:", err)
			} else {
				fmt.Println("Usuario agregado exitosamente.")
			}

		case 2:
			var titulo, autor, isbn string
			fmt.Print("Ingrese el titulo del libro: ")
			fmt.Scan(&titulo)
			fmt.Print("Ingrese el autor del libro: ")
			fmt.Scan(&autor)

			fmt.Print("Ingrese el ISBN del libro: ")
			fmt.Scan(&isbn)
			err := biblioteca.AgregarLibro(titulo, autor, isbn)
			if err != nil {
				fmt.Println("Error al agregar libro:", err)
			} else {
				fmt.Println("Libro agregado exitosamente.")
			}
		case 3:
			var usuarioId int
			var librosIds []int
			numeroLibros := 0
			fmt.Print("Ingrese el ID del usuario: ")
			fmt.Scan(&usuarioId)

			fmt.Print("Ingrese cuantos libros van a ser prestados: ")
			fmt.Scan(&numeroLibros)
			for i := 0; i < numeroLibros; i++ {
				var libroId int
				fmt.Printf("Ingrese el ID del libro %d: ", i+1)
				fmt.Scan(&libroId)
				librosIds = append(librosIds, libroId)
			}
			err := biblioteca.PrestarLibro(usuarioId, librosIds...)
			if err != nil {
				fmt.Println("Error al prestar libro:", err)
			} else {
				fmt.Println("Libro prestado exitosamente.")
			}

		case 4:
			var prestamoId int
			fmt.Print("Ingrese el ID del prestamo a devolver: ")
			fmt.Scan(&prestamoId)
			err := biblioteca.DevolverLibro(prestamoId)
			if err != nil {
				fmt.Println("Error al devolver libro:", err)
			} else {
				fmt.Println("Libro devuelto exitosamente.")
			}
		case 5:
			fmt.Println("Usuarios registrados:")
			for _, usuario := range biblioteca.usuarios {
				fmt.Printf("ID: %d, Nombres: %s, Apellidos: %s, Email: %s, Edad: %d\n",
					usuario.GetId(), usuario.GetNombres(), usuario.GetApellidos(), usuario.GetEmail(), usuario.GetEdad())
			}
		case 6:
			fmt.Println("Libros registrados:")
			for _, libro := range biblioteca.libros {
				fmt.Printf("ID: %d, Titulo: %s, Autor: %s, ISBN: %s, Prestado: %t\n",
					libro.GetId(), libro.GetTitulo(), libro.GetAutor(), libro.GetIsbn(), libro.EstaPrestado())
			}
		case 7:
			fmt.Println("Saliendo de la biblioteca. ¡Hasta luego!")
			return
		default:
			fmt.Println("Opcion no valida, por favor intente de nuevo.")
		}

	}
}
