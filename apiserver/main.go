package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v5"

	"github.com/aethanol/challenges-aethanol/apiserver/handlers"
	"github.com/aethanol/challenges-aethanol/apiserver/middleware"
	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
)

const defaultPort = "443"

const (
	//     /v1/users: UsersHandler
	//     /v1/sessions: SessionsHandler
	//     /v1/sessions/mine: SessionsMineHandler
	//     /v1/users/me: UsersMeHandler
	apiRoot         = "/v1/"
	apiSummary      = apiRoot + "summary"
	apiUsers        = apiRoot + "users"
	apiSessions     = apiRoot + "sessions"
	apiSessionsMine = apiSessions + "/mine"
	apiUsersMe      = apiUsers + "/me"
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
	if len(tlsCertPath) == 0 || len(tlsCertPath) == 0 {
		log.Fatal("one of the TLS env variable not set")
	}
	// 	Read the following new environment variables:

	// read and use the following environment variables
	// when initalizing the handlers context for authorization
	//     SESSIONKEY: a string to use as the session ID signing key
	//     REDISADDR: the address of your redis session store
	//     DBADDR: the address of your database server

	sessionKey := os.Getenv("SESSIONKEY")
	if len(sessionKey) == 0 {
		log.Fatal("no SESSIONKEY env variable set")
	}
	// Use the REDISADDR to create a new redis Client
	reddisAddr := os.Getenv("REDISADDR")
	roptions := redis.Options{
		Addr: reddisAddr,
	}
	reddisClient := redis.NewClient(&roptions)

	// pass the client to a new redis store -1 session duration for default duration
	sesStore := sessions.NewRedisStore(reddisClient, -1)

	// Use the DBADDR to dial your MongoDB server
	dbAddr := os.Getenv("DBADDR")
	fmt.Printf("dialing mongo server at %s...\n", dbAddr)
	mongoSession, err := mgo.Dial(dbAddr)
	if err != nil {
		log.Fatalf("error dialing mongo: %v", err)
	}

	// use the mongo session to create a new user store
	userStore, err := users.NewMongoStore(mongoSession, "production")
	if err != nil {
		log.Fatalf("error creating user store: %v", err)
	}

	// Create and initialize a new handlers.Context with the signing key,
	// the session store, and the user store.
	hctx := &handlers.Context{
		SessionKey:   sessionKey,
		SessionStore: sesStore,
		UserStore:    userStore,
	}

	// Create a new mux handlers to it
	mux := http.NewServeMux()
	mux.HandleFunc(apiUsers, hctx.UsersHandler)
	mux.HandleFunc(apiSessions, hctx.SessionsHandler)
	mux.HandleFunc(apiSessionsMine, hctx.SessionsMineHandler)
	mux.HandleFunc(apiUsersMe, hctx.UsersMeHanlder)
	//add your handlers.SummaryHandler function as a handler
	//for the apiSummary route
	//HINT: https://golang.org/pkg/net/http/#HandleFunc
	mux.HandleFunc(apiSummary, handlers.SummaryHandler)

	// apiRoot as the path, and the result of calling Adapt() on your new mux.
	// Specify the middleware.CORS() adapter as the only adapter
	http.Handle(apiRoot, middleware.Adapt(mux, middleware.CORS("", "", "", "")))

	//start your web server and use log.Fatal() to log
	//any errors that occur if the server can't start
	//HINT: https://golang.org/pkg/net/http/#ListenAndServe
	fmt.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, nil))
}
