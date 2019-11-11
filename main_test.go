package main

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"

	"github.com/dixneuf19/go-crud-api/greetings"
	"github.com/dixneuf19/go-crud-api/server"
)

type testRequestParams struct {
	Method  string
	Path    string
	Headers http.Header
	Body    string
	ExpCode int
	ExpMsg  string
}

const msgNotFound = "404 page not found\n"

func testHandlerFunc(t *testing.T, handlerFunc http.HandlerFunc, params testRequestParams) {
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

func testHandler(t *testing.T, handler http.Handler, params testRequestParams) {
	req := httptest.NewRequest(params.Method, params.Path, strings.NewReader(params.Body))
	// req := httptest.NewRequest(params.Method, params.Path, nil)

	w := httptest.NewRecorder()

	for key, values := range params.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	handler.ServeHTTP(w, req)

	if status := w.Code; status != params.ExpCode {
		t.Errorf("expecting %d instead of %d", params.ExpCode, status)
	}

	if body := w.Body.String(); body != params.ExpMsg {
		t.Errorf("expecting body '%s', not '%s'", params.ExpMsg, body)
	}
}

func TestHelloHandlerFunc(t *testing.T) {

	requests := []testRequestParams{
		{
			Method:  http.MethodGet,
			Path:    "/",
			ExpCode: http.StatusOK,
			ExpMsg:  "Hello",
		},
		{
			Method:  http.MethodPost,
			Path:    "/",
			ExpCode: http.StatusNotFound,
			ExpMsg:  msgNotFound,
		},
	}

	for _, r := range requests {
		testHandlerFunc(t, server.HelloHandlerFunc, r)
	}

}

func TestGreetingsHandler(t *testing.T) {

	gp := greetings.NewGreetingsMap()
	gp.Add("en", "hello")
	gp.Add("fr", "bonjour")

	header := http.Header{}
	header.Add("Content-Type", "application/json")

	requests := []testRequestParams{
		{
			Method:  http.MethodGet,
			Path:    "/hello",
			ExpCode: http.StatusBadRequest,
			ExpMsg:  "Please provide a language as 'lang' query parameter. Ex: /hello?lang=en",
		},
		{
			Method:  http.MethodGet,
			Path:    "/hello?language=en",
			ExpCode: http.StatusBadRequest,
			ExpMsg:  "Please provide a language as 'lang' query parameter. Ex: /hello?lang=en",
		},
		{
			Method:  http.MethodGet,
			Path:    "/hello?lang=en",
			ExpCode: http.StatusOK,
			ExpMsg:  "hello",
		},
		{
			Method:  http.MethodGet,
			Path:    "/hello?lang=fr",
			ExpCode: http.StatusOK,
			ExpMsg:  "bonjour",
		},
		{
			Method:  http.MethodGet,
			Path:    "/hello?lang=it",
			ExpCode: http.StatusBadRequest,
			ExpMsg:  "I don't know how to greet in 'it'. Learn me how with a POST method",
		},
		{
			Method:  http.MethodPost,
			Path:    "/hello?lang=it",
			Headers: header,
			Body:    `{"language": "it", "hello": "buongiorno"}`,
			ExpCode: http.StatusCreated,
			ExpMsg:  "",
		},
		{
			Method:  http.MethodGet,
			Path:    "/hello?lang=it",
			ExpCode: http.StatusOK,
			ExpMsg:  "buongiorno",
		},
		{
			Method:  http.MethodDelete,
			Path:    "/hello?lang=it",
			ExpCode: http.StatusOK,
			ExpMsg:  "",
		},
		{
			Method:  http.MethodGet,
			Path:    "/hello?lang=it",
			ExpCode: http.StatusBadRequest,
			ExpMsg:  "I don't know how to greet in 'it'. Learn me how with a POST method",
		},
	}

	for _, r := range requests {
		testHandler(t, &server.GreetingsHandler{GP: gp}, r)
	}

}
