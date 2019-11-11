package main

import (
	"flag"
	"log"

	"github.com/dixneuf19/go-crud-api/greetings"
	"github.com/dixneuf19/go-crud-api/server"
)

func getProvider(key string) greetings.Provider {
	switch key {
	case "map":
		return greetings.NewGreetingsMap()
	// case "redis":
	// 	return
	default:
		return greetings.NewGreetingsMap()
	}
}

func main() {
	var host = flag.String("host", "localhost", "the host of the server, default localhost")
	var port = flag.String("port", "8080", "the port of the server, default 8080")
	var provider = flag.String("provider", "map", "the greetings provider: map (in memory), redis")
	flag.Parse()

	gp := getProvider(*provider)
	address := *host + ":" + *port

	greetServer := server.NewGreetingServer(address, gp)

	err := greetServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
