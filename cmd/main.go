package main

import (
	"github.com/mdreem/gopass-operator/gopass-server"
	"log"
)

func main() {
	log.Printf("starting server\n")
	gopass_server.Run()
}
