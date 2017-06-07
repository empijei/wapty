package main

import (
	"log"

	"github.com/empijei/wapty/mocksy"
)

func main() {
	const port = ":8082"

	log.Printf("Starting mocksy server at %s\n", port)
	if err := mocksy.StartServer(port); err != nil {
		panic(err)
	}
}
