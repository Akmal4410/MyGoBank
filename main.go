package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Working on Go bank...")

	storage, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(":3000", storage)
	server.Run()
}
