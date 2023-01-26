package main

import (
	"fmt"
	"log"
	"os"

	"01.alem.school/git/aseitkha/net-cat/server"
)

func main() {
	if len(os.Args) > 2 {
		log.Println("Wrong usage")
	}
	port := "8989"

	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	addr := fmt.Sprintf("localhost:%s", port)
	serv, err := server.CreateNewServer(addr)
	defer serv.Listener.Close() // Defer closing the listener
	if err != nil {
		log.Println("Error while starting the server")
		return
	}
	log.Println("Started the server at", addr)
	serv.RunServer()
}
