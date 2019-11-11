package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", HelloHandler)

	http.ListenAndServe(":8080", nil)

}

// HelloHandler greats you
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello world!")
}
