package main 
import (
"fmt")

func main (){
	
}

type Product struct{
	id string
	nombre string
	precio float64
		
}
type Client struct{
	nombre string
	apellido string
	cedula string
	correo string
}
type ShoppingCart struct{
	productos []Product 
	total float64
	
}