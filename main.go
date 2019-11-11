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

type greetingsProvider interface {
	Get(string) (string, bool)
	Add(string, string) error
	Delete(string) error
}

var g greetingsProvider

func main() {
	g = greetings.NewGreetings()
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
		GetGreetingsHandler(w, req)
		return

	case http.MethodPost:
		PostGreetingsHandler(w, req)
		return

	case http.MethodDelete:
		DeleteGreetingsHandler(w, req)

	default:
		http.NotFoundHandler().ServeHTTP(w, req)
		return
	}
}

// GetGreetingsHandler handles GET greetings requests
func GetGreetingsHandler(w http.ResponseWriter, req *http.Request) {
	language := req.URL.Query().Get("lang")
	if len(language) == 0 {
		BadRequestHandler(w, req, fmt.Errorf("Please provide a language as 'lang' query parameter. Ex: /hello?lang=en"))
		return
	}

	greet, ok := g.Get((language))
	if !ok {
		BadRequestHandler(w, req, fmt.Errorf("I don't know how to greet in '%s'. Learn me how with a POST method", language))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", greet)
	return
}

// PostGreetingsHandler handles POST greetings requests
func PostGreetingsHandler(w http.ResponseWriter, req *http.Request) {
	if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
		BadRequestHandler(w, req, fmt.Errorf("only accepts %s for body, not %s", "application/json", contentType))
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		BadRequestHandler(w, req, fmt.Errorf("unable to read body: %s", err))
		return
	}

	var greetRequest GreetRequest
	err = json.Unmarshal(reqBody, &greetRequest)
	if err != nil {
		BadRequestHandler(w, req, fmt.Errorf(
			`unable to parse the body: %s
			 please provide a greeting like {"language": "fr", "hello": "bonsoir"}`,
			string(reqBody)))
		return
	}

	if len(greetRequest.Greet) == 0 || len(greetRequest.Language) == 0 {
		BadRequestHandler(w, req, fmt.Errorf(`Please provide a new greeting, for example: {"language": "fr", "hello": "bonsoir"}`))
		return
	}

	err = g.Add(greetRequest.Language, greetRequest.Greet)
	if err != nil {
		BadRequestHandler(w, req, fmt.Errorf("Cannot add greet '%s' for language '%s': %s", greetRequest.Greet, greetRequest.Language, err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

// DeleteGreetingsHandler handles DELETE greetings requests
func DeleteGreetingsHandler(w http.ResponseWriter, req *http.Request) {
	language := req.URL.Query().Get("lang")
	if len(language) == 0 {
		BadRequestHandler(w, req, fmt.Errorf("Please provide a language as 'lang' query parameter. Ex: /hello?lang=en"))
		return
	}

	err := g.Delete(language)
	if err != nil {
		BadRequestHandler(w, req, fmt.Errorf("unable to delete '%s' entry: %s", language, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

// BadRequestHandler handle incorrect requests
func BadRequestHandler(w http.ResponseWriter, req *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, err.Error())
	return
}
