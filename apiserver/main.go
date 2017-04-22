package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aethanol/challenges-aethanol/apiserver/handlers"
)

const defaultPort = "443"

const (
	apiRoot    = "/v1/"
	apiSummary = apiRoot + "summary"
)

//main is the main entry point for this program
func main() {
	//read and use the following environment variables
	//when initializing and starting your web server
	// PORT - port number to listen on for HTTP requests (if not set, use defaultPort)
	// HOST - host address to respond to (if not set, leave empty, which means any host)
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if len(port) == 0 {
		port = defaultPort
	}
	// concat the host and port to a valid address
	addr := fmt.Sprintf("%s:%s", host, port)

	//get the TLS key and cert paths from environment variables
	//this allows us to use a self-signed cert/key during development
	//and the Let's Encrypt cert/key in production
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	//add your handlers.SummaryHandler function as a handler
	//for the apiSummary route
	//HINT: https://golang.org/pkg/net/http/#HandleFunc
	http.HandleFunc(apiSummary, handlers.SummaryHandler)

	//start your web server and use log.Fatal() to log
	//any errors that occur if the server can't start
	//HINT: https://golang.org/pkg/net/http/#ListenAndServe
	fmt.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, nil))
}
