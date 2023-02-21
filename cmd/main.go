package main

import (
	"log"

	"proxyserver/internal/server"
	"proxyserver/web"
)

const port = "8080"

func main() {
	server := new(server.Server)

	handler := web.NewHandler()

	if err := server.Run(port, handler.InitRoutes()); err != nil {
		log.Printf("error while running the server: %s\n", err)
		return
	}
}
