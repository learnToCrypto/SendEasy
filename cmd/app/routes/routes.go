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
	//home when logged in
	mux.HandleFunc("/home", Index)
	// defined in route_qoute.go
	mux.HandleFunc("/demand/new", NewDemand)
	mux.HandleFunc("/demand/create/1", CreateDemand1)
	mux.HandleFunc("/demand/create/2", CreateDemandPriv)
	mux.HandleFunc("/demand/create/3", CreateDemandPub)
	mux.HandleFunc("/demand/post", PostDemand)
	mux.HandleFunc("/demand/read", ReadDemand)
	// Find Shipping
	mux.HandleFunc("/demand/list/", DemandList)
	// forum
	mux.HandleFunc("/forum", Forum)
	// about
	mux.HandleFunc("/howitworks", HowItWorks)
	// error
	mux.HandleFunc("/err", Err)

	// defined in auth.go
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/logout", Logout)
	mux.HandleFunc("/join", SignupChoice)
	mux.HandleFunc("/signup", Signup)
	mux.HandleFunc("/signup_account", SignupAccount)
	mux.HandleFunc("/signup_provider", SignupProvider)
	mux.HandleFunc("/authenticate", Authenticate)

	// defined in threads.go
	mux.HandleFunc("/thread/new", NewThread)
	mux.HandleFunc("/thread/create", CreateThread)
	mux.HandleFunc("/thread/post", PostThread)
	mux.HandleFunc("/thread/read", ReadThread)

	//
	mux.HandleFunc("/myshipments", MyShipments)
	mux.HandleFunc("/inbox", Inbox)
	return mux
}
