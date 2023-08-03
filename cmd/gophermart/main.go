package main

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app"
	"log"
)

func main() {
	server := app.NewServer()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
