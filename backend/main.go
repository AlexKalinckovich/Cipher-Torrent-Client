package main

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/routers"
	"log"
)

func main() {
	log.Printf("Server started")
	router := routers.SetupRouter()
	log.Fatal(router.Run(":8080"))
}
