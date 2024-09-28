package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Product struct {
	Id    int         `json:"id"`
	Title string      `json:"title"`
	Price json.Number `json:"price"`
}

var products = make(map[int]Product)
var mutex sync.RWMutex

func GetProducts(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	productsArray := []Product{}
	for _, v := range products {
		productsArray = append(productsArray, v)
	}
	mutex.RUnlock()

	if len(products) < 1 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no data"))
		return
	}

	data, err := json.Marshal(productsArray)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.RLock()
	product, ok := products[id]
	mutex.RUnlock()

	if !ok {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "please provide all fields as designated (title as string and price as number)", http.StatusBadRequest)
		return
	}

	if product.Title == "" || product.Price == "" {
		http.Error(w, "all fields required (title, price)", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	product.Id = len(products) + 1
	products[product.Id] = product
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("done"))
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.RLock()
	_, ok := products[id]
	mutex.RUnlock()

	if !ok {
		http.Error(w, "product not found", http.StatusBadRequest)
		return
	}

	var product Product

	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "please provide all fields as designated (title as string and price as number)", http.StatusBadRequest)
		return
	}

	if product.Title == "" || product.Price == "" {
		http.Error(w, "all fields required (title, price)", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	product.Id = id
	products[id] = product
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("done"))
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.RLock()
	_, ok := products[id]
	mutex.RUnlock()

	if !ok {
		http.Error(w, "product not found", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	delete(products, id)
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("done"))
}
