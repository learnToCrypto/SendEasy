package routes

import (
	"net/http"

	"github.com/learnToCrypto/lakoposlati/internal/user"
)

func Forum(writer http.ResponseWriter, request *http.Request) {
	threads, err := user.Threads()
	if err != nil {
		error_message(writer, request, "Cannot get threads")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, threads, "layout/base", "public/navbar", "forum")
		} else {
			generateHTML(writer, threads, "layout/base", "private/navbar", "forum")
		}
	}
}
