package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dixneuf19/go-crud-api/greetings"
)

func main() {

	greetings.AddGreeting("en", "hello")
	greetings.AddGreeting("fr", "bonjour")

	http.HandleFunc("/", HelloHandler)

	http.HandleFunc("/hello", GreetingsHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

// HelloHandler greets you
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello")
}

// GreetingsHandler returns the appropriate greet for the given language
// By default, it replies in english
func GreetingsHandler(w http.ResponseWriter, req *http.Request) {
	language := req.URL.Query().Get("lang")
	// Default language is "en"
	if len(language) == 0 {
		fmt.Fprintf(w, "Please provide a language as query parameter. Ex: /hello?lang=en")
		return
	}

	greet, ok := greetings.GetGreetings()[language]
	if !ok {
		fmt.Fprintf(w, "I don't know how to greet in '%s'. Learn me how!", language)
		return
	}
	fmt.Fprintf(w, "%s", greet)
}
