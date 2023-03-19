// Packge classification of Product API

package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/codeshubham/semaphore-demo-go-gin/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling GET Request")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling POST Request")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
	}
	data.AddProduct(prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, errr := strconv.Atoi(vars["id"])
	if errr != nil {
		http.Error(rw, "Cannot convert id value to an integer", http.StatusInternalServerError)
		return
	}
	p.l.Println("Handle PUT Request")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	//fmt.Printf("%#v", prod)
	e := data.UpdateProduct(id, &prod)
	if e == data.ErrProductNotFound {
		http.Error(rw, "product not found", http.StatusNotFound)
		return
	}
	if e != nil {
		http.Error(rw, "Product not retreived", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct {
}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		//fmt.Printf("%#v", prod)
		err := (prod).FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			p.l.Println("Error in deserializing the product ", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
