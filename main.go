package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type httpError struct {
	Cause  error
	Detail string
	Status int
}

func (e *httpError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Detail + " : " + e.Cause.Error()
}

type rootHandler func(w http.ResponseWriter, r *http.Request) *httpError

func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpErr := fn(w, r)
	if httpErr.Cause == nil {
		return
	}
	log.Printf("An error occured: %v", httpErr.Cause.Error())
	body, err := json.Marshal(httpErr)
	if err != nil {
		log.Printf("Convert to json failed: %v", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(httpErr.Status)
	w.Write(body)
}

func NewHTTPError(err error, detail string, status int) *httpError {
	return &httpError{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}

func welcome(w http.ResponseWriter, r *http.Request) *httpError {
	err := fmt.Errorf("very bad error")
	return NewHTTPError(err, "too: bad", 405)
}

func main() {
	http.Handle("/welcome", rootHandler(welcome))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
