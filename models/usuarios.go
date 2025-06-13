package models

type Usuario struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Rol      string `json:"rol"`      // "admin" o "cliente"
}