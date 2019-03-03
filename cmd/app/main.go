package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/learnToCrypto/lakoposlati/cmd/app/routes"
	"github.com/learnToCrypto/lakoposlati/internal/config"
)

// Convenience function for printing to stdout
func p(a ...interface{}) {
	fmt.Println(a)
}

// version
func version() string {
	return "0.2"
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func main() {

	//var conf config.Configuration

	// =========================================================================
	// Configuration
	conf := config.LoadConfig()

	// =========================================================================
	// Logging

	file, err := os.OpenFile("lakoposlati.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger := log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)

	// =========================================================================
	// App Starting
	log.Printf("Lako Poslati : Started  at %v: Application Initializing version %q", conf.Address, version())
	defer log.Println("main : Completed")

	// starting up the server
	server := &http.Server{
		Addr:           conf.Address,
		Handler:        routes.API(logger, conf),
		ReadTimeout:    time.Duration(conf.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(conf.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	// Start the service listening for requests.

	log.Printf("main : Server Listening")
	server.ListenAndServe()
	//private key and cert generated using the following:
	//openssl ecparam -name prime256v1 -genkey -noout -out priv.pem
	// openssl req -new -x509 -key priv.pem -out EC_server.pem -days 365
	//	err = server.ListenAndServeTLS("EC_server.pem", "priv.pem")
	//	fmt.Println(err)
}
