package main

import (
	"log"

	"github.com/dixneuf19/go-crud-api/greetings"
	"github.com/dixneuf19/go-crud-api/server"
)

func main() {
	gp := greetings.NewGreetingsMap()

	greetServer := server.NewGreetingServer(":8080", gp)

	err := greetServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
