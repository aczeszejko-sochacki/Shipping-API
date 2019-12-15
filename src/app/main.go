package main

import (
	"net/http"
	"log"
	"server"
)

const port = ":8081"

func main() {
	c := server.ShippingServer("myServer")
	c.Routes()

	log.Fatal(http.ListenAndServe(port, nil))
}