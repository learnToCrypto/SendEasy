package routes

import (
	"log"
	"net/http"

	"github.com/learnToCrypto/lakoposlati/internal/config"
)

func API(log *log.Logger, conf config.Configuration) http.Handler {
	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(conf.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//
	// all route patterns matched here
	// route handler functions defined in other files
	//

	// index
	mux.HandleFunc("/", Index)
	// defined in route_qoute.go
	mux.HandleFunc("/demand/new", NewDemand)
	mux.HandleFunc("/demand/create", CreateDemand)
	mux.HandleFunc("/demand/post", PostDemand)
	mux.HandleFunc("/demand/read", ReadDemand)
	// Find Shipping
	mux.HandleFunc("/demand/list/", DemandList)
	// forum
	mux.HandleFunc("/forum", Forum)
	// about
	mux.HandleFunc("/about", About)
	// error
	mux.HandleFunc("/err", Err)

	// defined in route_auth.go
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/logout", Logout)
	mux.HandleFunc("/join", Signup)
	mux.HandleFunc("/signup", Signup)
	mux.HandleFunc("/signup_account", SignupAccount)
	mux.HandleFunc("/authenticate", Authenticate)

	// defined in route_thread.go
	mux.HandleFunc("/thread/new", NewThread)
	mux.HandleFunc("/thread/create", CreateThread)
	mux.HandleFunc("/thread/post", PostThread)
	mux.HandleFunc("/thread/read", ReadThread)

	return mux
}
