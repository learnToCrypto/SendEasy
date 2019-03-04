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
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(writer, nil)
}

// GET /signup
// Show the signup page
func Signup(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "login.layout", "public.navbar", "signup")
}

// POST /signup
// Create the user account
func SignupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		logger.Danger(err, "Cannot parse form")
	}
	user := user.User{
		Name:     request.PostFormValue("name"),
		Email:    request.PostFormValue("email"),
		Password: request.PostFormValue("password"),
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
