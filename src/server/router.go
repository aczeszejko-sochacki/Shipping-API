package server

import "net/http"

type Server interface {
	routes()
}

type ShippingServer string

// One can add more routes here
func(s ShippingServer) Routes() {
	http.HandleFunc("/routes", sortDestinations)
}