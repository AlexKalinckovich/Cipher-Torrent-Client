package main

import (
	openapi "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/go"
	"log"
)

func main() {
	routes := openapi.ApiHandleFunctions{}
	log.Printf("Server started")
	router := openapi.NewRouter(routes)
	log.Fatal(router.Run(":8080"))
}
