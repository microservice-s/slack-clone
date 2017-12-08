package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
)

// func (ctx *Context) getServiceProxy(w http.ResponseWriter, r *http.Request) *httputil.ReverseProxy {
// 	// return a reverse proxy
// 	return &httputil.ReverseProxy{
// 		Director: func(r *http.Request) {
// 			//user := getUser(r)
// 			// get the session state

// 			// reset the scheme and host of the request url
// 			r.URL.Scheme = "http" // terminate the https (so http faster to microservice)
// 			r.URL.Host = ctx.SvcAddr
// 			// i++
// 			// i = i % len(instances)
// 			j, _ := json.Marshal(state.User)
// 			r.Header.Add("X-User", string(j))
// 		},
// 	}
// }

func (ctx *Context) ChatbotHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session state
	state := &SessionState{}

	// get the state of the browser that is accessing their page
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state)
	if err != nil {
		http.Error(w, "error getting session state "+err.Error(),
			http.StatusForbidden)
		return
	}
	// proxy := ctx.getServiceProxy(w, r)
	// proxy.ServeHTTP(w, r)
	// get a proxy
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = ctx.SvcAddr
			j, _ := json.Marshal(state.User)
			r.Header.Add("X-User", string(j))
		},
	}

	proxy.ServeHTTP(w, r)
	//http.Handle("/v1/bot", ctx.getServiceProxy(w, r, ctx.SvcAddr))
}
