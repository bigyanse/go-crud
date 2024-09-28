package main

import (
	"fmt"
	"net/http"

	"github.com/bigyanse/go-crud/controllers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", controllers.GetProducts)
	mux.HandleFunc("GET /products/{id}", controllers.GetProduct)
	mux.HandleFunc("POST /products", controllers.CreateProduct)
	mux.HandleFunc("PUT /products/{id}", controllers.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", controllers.DeleteProduct)

	fmt.Println("Server listening to :8080")
	http.ListenAndServe(":8080", mux)
}
