package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// Sahi packages ko import karein
	"github.com/keshav78-78/ECOM/service/cart"
	"github.com/keshav78-78/ECOM/service/order"
	"github.com/keshav78-78/ECOM/service/product"
	"github.com/keshav78-78/ECOM/service/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// User service setup
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// Product service setup
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	// Order service setup
	// Yahaan 'order' package ka istemaal karein
	orderStore := order.NewStore(s.db)
	// (Agar order ke bhi routes hain, toh unhe yahaan register karein)
	// orderHandler := order.NewHandler(orderStore)
	// orderHandler.RegisterRoutes(subrouter)

	// Cart service setup
	// Ab 'cart' package se NewHandler ko call karein
	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
