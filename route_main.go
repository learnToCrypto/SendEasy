package main

import (
	"github.com/learnToCrypto/lakoposlati/data"
	"net/http"
)

// GET /err?msg=
// shows the error message page
func err(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		generateHTML(writer, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}

func forum(writer http.ResponseWriter, request *http.Request) {
	threads, err := data.Threads()
	if err != nil {
		error_message(writer, request, "Cannot get threads")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, threads, "layout", "public.navbar", "forum")
		} else {
			generateHTML(writer, threads, "layout", "private.navbar", "forum")
		}
	}
}

func index(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "public.navbar", "index")
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "index")
	}

}

func about(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		t := parseTemplateFiles("layout", "public.navbar", "about")
		t.Execute(writer, nil)
	} else {
		t := parseTemplateFiles("layout", "private.navbar", "about")
		t.Execute(writer, nil)
	}

}
