package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/codeshubham/semaphore-demo-go-gin/handlers"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"golang.org/x/net/context"
)

var bindAddress = env.String("BIN_ADDRESS", false, ":8090", "Bind address for the server")

func main() {

	env.Parse()

	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	ph := handlers.NewProducts(l)

	//hh := handlers.NewHello(l)
	//gh := handlers.NewGoodbye(l)

	//create a new router using gorilla framework
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)
	//sm.Handle("/", hh)
	//sm.Handle("/goodbye", gh)
	//sm.Handle("/products", ph)
	s := &http.Server{
		Addr:         *bindAddress,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Receive terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
