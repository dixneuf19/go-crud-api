package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dixneuf19/go-crud-api/greetings"
)

// GreetRequest is the standart format to add a new greet
type GreetRequest struct {
	Language string `json:"language"`
	Greet    string `json:"hello"`
}

var g greetings.Greetings = greetings.NewGreetings()

func main() {
	g.Add("en", "hello")
	g.Add("fr", "bonjour")

	http.HandleFunc("/", HelloHandler)
	http.HandleFunc("/hello", GreetingsHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

// HelloHandler greets you
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {

	case http.MethodGet:
		fmt.Fprint(w, "Hello")
		return

	default:
		http.NotFoundHandler().ServeHTTP(w, req)
		return
	}

}

// GreetingsHandler returns the appropriate greet for the given language
// By default, it replies in english
func GreetingsHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {

	case http.MethodGet:
		language := req.URL.Query().Get("lang")
		if len(language) == 0 {
			fmt.Fprintf(w, "Please provide a language as query parameter. Ex: /hello?lang=en")
			return
		}

		greet, ok := g[language]
		fmt.Printf("greet[%s]=%s", language, greet)
		if !ok {
			fmt.Fprintf(w, "I don't know how to greet in '%s'. Learn me how!", language)
			return
		}

		fmt.Fprintf(w, "%s", greet)
		return

	case http.MethodPost:
		if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
			fmt.Fprintf(w, "only accepts %s for body, not %s", "application/json", contentType)
			return
		}

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		var greetRequest GreetRequest
		err = json.Unmarshal(reqBody, &greetRequest)
		if err != nil {
			fmt.Fprintf(w,
				`unable to parse the body: %s
				 please provide a greeting like {"language": "fr", "hello": "bonsoir"}`,
				reqBody)
		}

		if len(greetRequest.Greet) == 0 || len(greetRequest.Language) == 0 {
			fmt.Fprintf(w,
				`Please provide a new greeting, for example: {"language": "fr", "hello": "bonsoir"}`)
			return
		}

		err = g.Add(greetRequest.Language, greetRequest.Greet)
		if err != nil {
			fmt.Printf("Cannot add greet '%s' for language '%s': %s", greetRequest.Greet, greetRequest.Language, err)
			return
		}
		w.WriteHeader(201)
		return

	default:
		http.NotFoundHandler().ServeHTTP(w, req)
		return
	}

	// func BadRequestHandler(w http.ResponseWriter, req *http.Request, err error) {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf
	// 	return
	// }

}
