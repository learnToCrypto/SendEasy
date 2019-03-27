package routes

import (
	"net/http"

	"github.com/learnToCrypto/lakoposlati/internal/user"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	//fmt.Println("Session in index:", sess)
	if err != nil {
		generateHTML(writer, nil, "layout/base", "public/navbar", "index")
	} else {
		username, err := user.UsernamebySession(sess.UserId)
		if err != nil {
			username = "My Account"
		}
		d := struct {
			Username string
		}{
			Username: username,
		}
		generateHTML(writer, d, "layout/base", "private/navbar", "index")
	}

}

func HowItWorks(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout/base", "public/navbar", "howitworks")
	} else {
		username, err := user.UsernamebySession(sess.UserId)
		if err != nil {
			username = "My Account"
		}
		d := struct {
			Username string
		}{
			Username: username,
		}
		generateHTML(writer, d, "layout/base", "private/navbar", "howitworks")
	}

}
