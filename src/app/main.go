package main

import (
	"net/http"
	"log"
	"server"
)

func main() {
	c := server.ShippingServer("myServer")
	c.Routes()

	log.Fatal(http.ListenAndServe(":8081", nil))
}