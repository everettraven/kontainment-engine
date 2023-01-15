package main

import (
	"fmt"

	"github.com/kontainment/engine/api/server"
)

func main() {
	srv, err := server.NewKontainmentServer()
	if err != nil {
		fmt.Println("Error creating server:", err)
	}

	err = srv.Serve()

	if err != nil {
		fmt.Println("Error running server:", err)
	}
}
