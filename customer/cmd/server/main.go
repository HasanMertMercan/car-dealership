package main

import (
	"car-dealership/customer/internal"
	"car-dealership/customer/transport/restserver"
	"fmt"
	"net/http"
)

func main() {
	application := internal.NewCustomerService()
	handler := restserver.NewHandler(application)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	fmt.Println("HTTP server for REST API listening on :8081")

	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}

}
