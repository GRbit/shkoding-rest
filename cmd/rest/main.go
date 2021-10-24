package main

import (
	"log"

	"github.com/GRbit/shkoding-rest/internal/server"
)

func main() {
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
