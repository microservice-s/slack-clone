package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

//openGraphPrefix is the prefix used for Open Graph meta properties
const openGraphPrefix = "og:"

//openGraphProps represents a map of open graph property names and values
type openGraphProps map[string]string

func getPageSummary(url string) (openGraphProps, error) {
	//Get the URL
	//If there was an error, return it
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	//ensure that the response body stream is closed eventually
	defer res.Body.Close()

	//HINTS: https://gobyexample.com/defer
	//https://golang.org/pkg/net/http/#Response

	//if the response StatusCode is >= 400
	//return an error, using the response's .Status
	//property as the error message
	if res.StatusCode >= 400 {
		return nil, errors.New(res.Status)
	}

	//if the response's Content-Type header does not
	//start with "text/html", return an error noting
	//what the content type was and that you were
	//expecting HTML
	ctype := res.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, fmt.Errorf("content type: %q, expexted text/html", ctype)
	}

	//create a new openGraphProps map instance to hold
	//the Open Graph properties you find
	//(see type definition above)
	props := make(openGraphProps)

	//tokenize the response body's HTML and extract
	//any Open Graph properties you find into the map,
	//using the Open Graph property name as the key, and the
	//corresponding content as the value.
	//strip the openGraphPrefix from the property name before
	//you add it as a new key, so that the key is just `title`
	//and not `og:title` (for example).

	//create a new tokenizer over the response body
	tokenizer := html.NewTokenizer(res.Body)
Loop:
	for {
		tt := tokenizer.Next()
		token := tokenizer.Token()
		fmt.Println(token.Attr)
		switch tt {
		case html.ErrorToken:
			//log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
			err := tokenizer.Err()
			if err == io.EOF {
				break Loop
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
			//return nil, fmt.Errorf("error tokenizing HTML: %v", tokenizer.Err())``
		// open graph properties only exist in the head
		case html.EndTagToken:
			token := tokenizer.Token()
			if token.Data == "head" {
				break Loop // using the go Label break "Loop"
			}
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "meta" {
				for i, a := range token.Attr {
					if strings.HasPrefix(a.Val, openGraphPrefix) {
						ogKey := strings.TrimPrefix(a.Val, openGraphPrefix)
						ogVal := token.Attr[i+1].Val
						//fmt.Printf("Key: %v Val: %v", ogKey, ogVal)
						props[ogKey] = ogVal
					}
				}
			}
		}

	}

	//HINTS: https://info344-s17.github.io/tutorials/tokenizing/
	//https://godoc.org/golang.org/x/net/html
	return props, nil
}

//SummaryHandler fetches the URL in the `url` query string parameter, extracts
//summary information about the returned page and sends those summary properties
//to the client as a JSON-encoded object.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {

	//Add the following header to the response
	//   Access-Control-Allow-Origin: *
	//this will allow JavaScript served from other origins
	//to call this API
	w.Header().Add("Access-Control-Allow-Origin", "*")

	//get the `url` query string parameter
	//if you use r.FormValue() it will also handle cases where
	//the client did POST with `url` as a form field
	//HINT: https://golang.org/pkg/net/http/#Request.FormValue
	URL := r.FormValue("url")
	//if no `url` parameter was provided, respond with
	//an http.StatusBadRequest error and return
	//HINT: https://golang.org/pkg/net/http/#Error

	if URL == "" {
		http.Error(w, "no url parameter provided", http.StatusBadRequest)
		return
	}

	//call getPageSummary() passing the requested URL
	//and holding on to the returned openGraphProps map
	//(see type definition above)
	//if you get back an error, respond to the client
	//with that error and an http.StatusBadRequest code
	graphProps, err := getPageSummary(URL)
	if err != nil {
		http.Error(w, "error requesting page summary: "+err.Error(), http.StatusBadRequest)
		return
	}

	//otherwise, respond by writing the openGrahProps
	//map as a JSON-encoded object
	//add the following headers to the response before
	//you write the JSON-encoded object:
	//   Content-Type: application/json; charset=utf-8
	//this tells the client that you are sending it JSON
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	encoder := json.NewEncoder(w) // write back to the client
	if err := encoder.Encode(graphProps); err != nil {
		http.Error(w, "error encoding json: "+err.Error(), http.StatusInternalServerError)
	}

}
