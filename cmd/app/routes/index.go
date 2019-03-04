package routes

import (
	"net/http"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "public.navbar", "index")
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "index")
	}

}

func About(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		t := parseTemplateFiles("layout", "public.navbar", "about")
		t.Execute(writer, nil)
	} else {
		t := parseTemplateFiles("layout", "private.navbar", "about")
		t.Execute(writer, nil)
	}

}
