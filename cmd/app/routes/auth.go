package routes

import (
	"net/http"

	"github.com/learnToCrypto/lakoposlati/internal/logger"
	"github.com/learnToCrypto/lakoposlati/internal/platform/postgres"
	"github.com/learnToCrypto/lakoposlati/internal/user"
)

// GET /login
// Show the login page
func Login(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "layout/login.base", "public/navbar", "login")
}

// GET /signup
// Show the signup page
func Signup(writer http.ResponseWriter, request *http.Request) {
	if request.PostFormValue("user-type") == "customer" {
		generateHTML(writer, nil, "layout/login.base", "public/navbar", "signup/customer")
	} else {
		generateHTML(writer, nil, "layout/login.base", "public/navbar", "signup/provider")
	}
}

// GET /signup
// Show the signup page
func SignupChoice(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "layout/login.base", "public/navbar", "signup/choice")
}

// POST /signup
// Create the user account
func SignupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		logger.Danger(err, "Cannot parse form")
	}
	user := user.User{
		FirstName: request.PostFormValue("first_name"),
		LastName:  request.PostFormValue("last_name"),
		Name:      request.PostFormValue("first_name") + " " + request.PostFormValue("last_name"),
		Email:     request.PostFormValue("email"),
		Password:  request.PostFormValue("password"),
	}
	if err := user.Create(); err != nil {
		logger.Danger(err, "Cannot create user")
	}
	http.Redirect(writer, request, "/login", 302)
}

// POST /authenticate
// Authenticate the user given the email and password
func Authenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	user, err := user.UserByEmail(request.PostFormValue("email"))
	if err != nil {
		http.Error(writer, "Cannot find user", http.StatusForbidden)
		//logger.Danger(err, "Cannot find user")
		http.Error(writer, "Username and/or password do not match", http.StatusForbidden)
	}
	if user.Password == postgres.Encrypt(request.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			logger.Danger(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
		//fmt.Println(user, "is logged in")
		http.Redirect(writer, request, "/", 302)
	} else {
		http.Error(writer, "Username and/or password do not match 2", http.StatusForbidden)
		http.Redirect(writer, request, "/login", 302)
	}

}

// GET /logout
// Logs the user out
func Logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		//logger.Warning(err, "Failed to get cookie")
		session := user.Session{Uuid: cookie.Value}
		//fmt.Println(session)
		session.DeleteByUUID()
	}
	http.Redirect(writer, request, "/", 302)
}
