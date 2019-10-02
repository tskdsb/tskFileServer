package main

import (
	"log"
	"os"

	"github.com/tskdsb/tskFileServer/server"
)

func main() {



	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Printf("pwd: %s\n", dir)

	go server.RunServer()
	select {}
}
