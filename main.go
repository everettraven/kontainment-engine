package main

import (
	"fmt"

	"github.com/kontainment/engine/api/server"
)

func main() {
	srv := server.NewKontainmentServer()
	err := srv.Serve()

	if err != nil {
		fmt.Println("Error running server:", err)
	}
}
