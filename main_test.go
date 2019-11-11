package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRequestParams struct {
	Method  string
	Path    string
	ExpCode int
	ExpMsg  string
}

const msgNotFound = "404 page not found\n"

func testHandler(t *testing.T, handlerFunc http.HandlerFunc, params testRequestParams) {
	req := httptest.NewRequest(params.Method, params.Path, nil)
	w := httptest.NewRecorder()

	handlerFunc(w, req)

	if status := w.Code; status != params.ExpCode {
		t.Errorf("expecting %d instead of %d", params.ExpCode, status)
	}

	if body := w.Body.String(); body != params.ExpMsg {
		t.Errorf("expecting body '%s', not '%s'", params.ExpMsg, body)
	}
}

func TestHelloHandler(t *testing.T) {

	requests := []testRequestParams{
		{
			http.MethodGet,
			"/",
			http.StatusOK,
			"Hello",
		},
		{
			http.MethodPost,
			"/",
			http.StatusNotFound,
			msgNotFound,
		},
	}

	for _, r := range requests {
		testHandler(t, HelloHandler, r)
	}

}

func TestGetGreetingsHandler(t *testing.T) {

	requests := []testRequestParams{
		{
			http.MethodGet,
			"/hello",
			http.StatusBadRequest,
			"Please provide a language as 'lang' query parameter. Ex: /hello?lang=en",
		},
		{
			http.MethodGet,
			"/hello?language=en",
			http.StatusBadRequest,
			"Please provide a language as 'lang' query parameter. Ex: /hello?lang=en",
		},
	}

	for _, r := range requests {
		testHandler(t, GreetingsHandler, r)
	}

}
