package main

import (
	"flag"
	"log"

	"github.com/dixneuf19/go-crud-api/greetings"
	"github.com/dixneuf19/go-crud-api/server"
)

var host = flag.String("host", "localhost", "the host of the server, default localhost")
var port = flag.String("port", "8080", "the port of the server, default 8080")
var provider = flag.String("provider", "map", "the greetings provider: map (in memory), redis")
var redisHost = flag.String("redisHost", "127.0.0.1:6379", "address of the redis cluster")

func getProvider(key string) greetings.Provider {
	switch key {
	case "map":
		log.Println("use 'map' greetings provider")
		return greetings.NewGreetingsMap()
	case "redis":
		log.Println("use 'redis' greetings provider")
		return greetings.NewGreetingsRedis(*redisHost)
	default:
		return greetings.NewGreetingsMap()
	}
}

func main() {

	flag.Parse()

	gp := getProvider(*provider)
	address := *host + ":" + *port

	greetServer := server.NewGreetingServer(address, gp)

	log.Println("Launching the server!")
	err := greetServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
