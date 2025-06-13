package main

import (
	"Tiendify/models"
	"Tiendify/utils"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)


func main() {
    var usuarios []models.Usuario
    var productos []models.Producto
    var clientes []models.Cliente

    utils.LoadJSON("data/usuarios.json", &usuarios)
    utils.LoadJSON("data/productos.json", &productos)
    utils.LoadJSON("data/clientes.json", &clientes)

    fmt.Println("===== Bienvenido al sistema de e-commerce =====")
    fmt.Print("¿Tienes una cuenta? (s/n): ")
    op, _ := reader.ReadString('\n')
    op = strings.TrimSpace(strings.ToLower(op))

    if op == "n" {
        fmt.Println("Registrarse como cliente")
        fmt.Print("Elige un nombre de usuario: ")
        username, _ := reader.ReadString('\n')
        username = strings.TrimSpace(username)

        fmt.Print("Elige una contraseña: ")
        password, _ := reader.ReadString('\n')
        password = strings.TrimSpace(password)

        // Verificar si ya existe
        for _, u := range usuarios {
            if u.Username == username {
                fmt.Println("❌ Ese usuario ya existe.")
                return
            }
        }

        nuevo := models.Usuario{
            Username: username,
            Password: password,
            Rol:      "cliente",
        }
        usuarios = append(usuarios, nuevo)
        utils.SaveJSON("data/usuarios.json", usuarios)
        fmt.Println("✅ Cuenta creada. Inicia sesión para continuar.")
    }

    // Login normal
    fmt.Print("Usuario: ")
    username, _ := reader.ReadString('\n')
    username = strings.TrimSpace(username)

    fmt.Print("Contraseña: ")
    password, _ := reader.ReadString('\n')
    password = strings.TrimSpace(password)

    usuarioLogueado := login(usuarios, username, password)
    if usuarioLogueado == nil {
        fmt.Println("❌ Usuario o contraseña incorrecta.")
        return
    }

    if usuarioLogueado.Rol == "admin" {
        menuAdmin(&clientes, &productos, &usuarios)
    } else {
        menuCliente(&productos)
    }
}
func menuCliente(productos *[]models.Producto) {
    var carrito []models.Producto

    for {
        fmt.Println("\n--- Menú Cliente ---")
        fmt.Println("1. Ver productos")
        fmt.Println("2. Agregar producto al carrito")
        fmt.Println("3. Ver carrito")
        fmt.Println("4. Finalizar compra")
        fmt.Println("5. Salir")
        fmt.Print("Opción: ")

        op, _ := reader.ReadString('\n')
        op = strings.TrimSpace(op)

        switch op {
        case "1":
            for _, p := range *productos {
                fmt.Printf("ID: %d | %s | $%.2f | Stock: %d\n", p.ID, p.Nombre, p.Precio, p.Stock)
            }
        case "2":
            fmt.Print("ID del producto a agregar: ")
            idStr, _ := reader.ReadString('\n')
            idStr = strings.TrimSpace(idStr)

            for _, p := range *productos {
                if fmt.Sprintf("%d", p.ID) == idStr {
                    carrito = append(carrito, p)
                    fmt.Println("Producto agregado al carrito.")
                    break
                }
            }
        case "3":
            fmt.Println("Carrito:")
            for _, p := range carrito {
                fmt.Printf("- %s | $%.2f\n", p.Nombre, p.Precio)
            }
        case "4":
            var total float64
            for _, p := range carrito {
                total += p.Precio
            }
            fmt.Printf("Total a pagar: $%.2f\n", total)
            carrito = nil
        case "5":
            return
        default:
            fmt.Println("Opción inválida.")
        }
    }
}
func login(usuarios []models.Usuario, username, password string) *models.Usuario {
    for _, u := range usuarios {
        if u.Username == username && u.Password == password {
            return &u
        }
    }
    return nil
}
func menuAdmin(clientes *[]models.Cliente, productos *[]models.Producto, usuarios *[]models.Usuario) {
for {
fmt.Println("\n====== Menú Administrador ======")
fmt.Println("1. Ver Clientes")
fmt.Println("2. Agregar Cliente")
fmt.Println("3. Ver Productos")
fmt.Println("4. Agregar Producto")
fmt.Println("5. editar Usuarios")
fmt.Println("6. Salir al Login")
fmt.Print("Seleccione una opción: ")

opcionStr, _ := reader.ReadString('\n')
opcionStr = strings.TrimSpace(opcionStr)
opcion, _ := strconv.Atoi(opcionStr)

switch opcion {
case 1:
fmt.Println("\n--- Lista de Clientes ---")
for _, c := range *clientes {
fmt.Printf("ID: %d | Nombre: %s | Correo: %s\n", c.ID, c.Nombre, c.Correo)
}
case 2:
fmt.Print("Nombre del cliente: ")
nombre, _ := reader.ReadString('\n')
nombre = strings.TrimSpace(nombre)

fmt.Print("Correo del cliente: ")
correo, _ := reader.ReadString('\n')
correo = strings.TrimSpace(correo)

nuevo := models.Cliente{
ID:     len(*clientes) + 1,
Nombre: nombre,
Correo: correo,
}
*clientes = append(*clientes, nuevo)
utils.SaveJSON("data/clientes.json", clientes)
fmt.Println("✅ Cliente agregado.")
case 3:
fmt.Println("\n--- Lista de Productos ---")
for _, p := range *productos {
fmt.Printf("ID: %d | %s | $%.2f | Stock: %d\n", p.ID, p.Nombre, p.Precio, p.Stock)
}
case 4:
fmt.Print("Nombre del producto: ")
nombre, _ := reader.ReadString('\n')
nombre = strings.TrimSpace(nombre)

fmt.Print("Precio del producto: ")
precioStr, _ := reader.ReadString('\n')
precioStr = strings.TrimSpace(precioStr)
precio, _ := strconv.ParseFloat(precioStr, 64)

fmt.Print("Stock inicial: ")
stockStr, _ := reader.ReadString('\n')
stockStr = strings.TrimSpace(stockStr)
stock, _ := strconv.Atoi(stockStr)

nuevo := models.Producto{
ID:     len(*productos) + 1,
Nombre: nombre,
Precio: precio,
Stock:  stock,
}
*productos = append(*productos, nuevo)
utils.SaveJSON("data/productos.json", productos)
fmt.Println("✅ Producto agregado.")
case 5:
fmt.Println("\n--- Editar Usuarios ---")
editarUsuarios()
fmt.Println("✅ Usuarios actualizados.")
case 6:
fmt.Println("Saliendo al login...")
return
default:
fmt.Println("❌ Opción inválida.")
}
}
}
func editarUsuarios() {
    var usuarios []models.Usuario
    utils.LoadJSON("data/usuarios.json", &usuarios)

    fmt.Println("\n--- Lista de usuarios ---")
    for i, u := range usuarios {
        fmt.Printf("%d. %s (Rol: %s)\n", i+1, u.Username, u.Rol)
    }

    fmt.Print("Número de usuario a editar: ")
    indexStr, _ := reader.ReadString('\n')
    indexStr = strings.TrimSpace(indexStr)
    index, _ := strconv.Atoi(indexStr)
    index-- // para usar índice 0-based

    if index < 0 || index >= len(usuarios) {
        fmt.Println("❌ Índice inválido.")
        return
    }

    fmt.Print("Nuevo nombre (dejar vacío para mantener): ")
    nombre, _ := reader.ReadString('\n')
    nombre = strings.TrimSpace(nombre)

    fmt.Print("Nueva contraseña (dejar vacío para mantener): ")
    pass, _ := reader.ReadString('\n')
    pass = strings.TrimSpace(pass)

    fmt.Print("Nuevo rol (admin/cliente, dejar vacío para mantener): ")
    rol, _ := reader.ReadString('\n')
    rol = strings.TrimSpace(rol)

    if nombre != "" {
        usuarios[index].Username = nombre
    }
    if pass != "" {
        usuarios[index].Password = pass
    }
    if rol == "admin" || rol == "cliente" {
        usuarios[index].Rol = rol
    }

    utils.SaveJSON("data/usuarios.json", usuarios)
    fmt.Println("✅ Usuario actualizado.")
}