package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dixneuf19/go-crud-api/greetings"
)

// GreetingServer is a http server with its greeting provider attached
type GreetingServer struct {
	server           *http.Server
	greetingProvider *greetings.Provider
}

// NewGreetingServer returns a http server with its greeting provider attached
func NewGreetingServer(addr string, gp greetings.Provider) *GreetingServer {

	gp.Add("en", "hello")
	gp.Add("fr", "bonjour")

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", HelloHandlerFunc)
	serverMux.Handle("/hello", &GreetingsHandler{GP: gp})

	server := &http.Server{
		Addr:    addr,
		Handler: serverMux,
	}

	return &GreetingServer{
		server,
		&gp,
	}
}

// ListenAndServe always returns a non-nil error.
func (gs *GreetingServer) ListenAndServe() error {
	return gs.server.ListenAndServe()
}

// HelloHandlerFunc greets you
func HelloHandlerFunc(w http.ResponseWriter, req *http.Request) {
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
type GreetingsHandler struct {
	GP greetings.Provider
}

func (h *GreetingsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	switch req.Method {

	case http.MethodGet:
		h.GetGreetingsHandler(w, req)
		return

	case http.MethodPost:
		h.PostGreetingsHandler(w, req)
		return

	case http.MethodDelete:
		h.DeleteGreetingsHandler(w, req)

	default:
		http.NotFoundHandler().ServeHTTP(w, req)
		return
	}
}

// GetGreetingsHandler handles GET greetings requests
func (h *GreetingsHandler) GetGreetingsHandler(w http.ResponseWriter, req *http.Request) {
	language := req.URL.Query().Get("lang")
	if len(language) == 0 {
		BadRequestHandler(w, req, fmt.Errorf("Please provide a language as 'lang' query parameter. Ex: /hello?lang=en"))
		return
	}

	greet, ok := h.GP.Get(language)
	if !ok {
		BadRequestHandler(w, req, fmt.Errorf("I don't know how to greet in '%s'. Learn me how with a POST method", language))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", greet)
	return
}

// GreetRequest is the standart format to add a new greet
type GreetRequest struct {
	Language string `json:"language"`
	Greet    string `json:"hello"`
}

// PostGreetingsHandler handles POST greetings requests
func (h *GreetingsHandler) PostGreetingsHandler(w http.ResponseWriter, req *http.Request) {
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

	err = h.GP.Add(greetRequest.Language, greetRequest.Greet)
	if err != nil {
		BadRequestHandler(w, req, fmt.Errorf("Cannot add greet '%s' for language '%s': %s", greetRequest.Greet, greetRequest.Language, err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

// DeleteGreetingsHandler handles DELETE greetings requests
func (h *GreetingsHandler) DeleteGreetingsHandler(w http.ResponseWriter, req *http.Request) {
	language := req.URL.Query().Get("lang")
	if len(language) == 0 {
		BadRequestHandler(w, req, fmt.Errorf("Please provide a language as 'lang' query parameter. Ex: /hello?lang=en"))
		return
	}

	err := h.GP.Delete(language)
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
