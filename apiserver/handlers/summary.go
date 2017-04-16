package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"net/url"

	"path"

	"golang.org/x/net/html"
)

//openGraphPrefix is the prefix used for Open Graph meta properties
const openGraphPrefix = "og:"

//openGraphProps represents a map of open graph property names and values
type openGraphProps map[string]string

// fetchHTML fetches the html body from the given url
func fetchHTML(URL string) (io.ReadCloser, error) {
	//Get the URL
	//If there was an error, return it
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

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

	return res.Body, nil
}

// fetchOpenGraphProps tokenizes and stores the open graph properties
// from a html body
func fetchOpenGraphProps(body io.ReadCloser, URL string) (openGraphProps, error) {
	//create a new openGraphProps map instance to hold
	//the Open Graph properties you find
	//(see type definition above)
	props := make(openGraphProps)

	//create a new tokenizer over the response body
	tokenizer := html.NewTokenizer(body)
Loop:
	for {
		// get the next token type for the switch
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			// return the error if it is NOT io.EOF (EOF just means we hit the end of the page)
			eofErr := tokenizer.Err()
			if eofErr == io.EOF && eofErr != nil {
				break Loop
			}
			return nil, fmt.Errorf("error tokenizing HTML: %v", tokenizer.Err())
		// open graph properties only exist in the head
		// so we can break early
		case html.EndTagToken:
			token := tokenizer.Token()
			if token.Data == "head" {
				break Loop // using the go Label break "Loop"
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			// EXTRA CREDIT FALLBACK title
			// check if the props map doesn't already contain a title
			if _, contains := props["title"]; !contains && token.Data == "title" {
				//the next token should be the page title
				tt = tokenizer.Next()
				// make sure that it's actually a text token
				if tt == html.TextToken {
					props["title"] = tokenizer.Token().Data
					break // break to continue the loop
				}
				// look for property and content fields in only meta tags
			} else if token.Data == "meta" {
				var prop, cont, name string
				// get the content and property fields (handles order)
				for _, a := range token.Attr {
					switch a.Key {
					case "property":
						prop = a.Val
					case "content":
						cont = a.Val
					case "name":
						name = a.Val
					}
				}
				// trim the open graph meta tag property
				// and then add it to the map with the content
				if strings.HasPrefix(prop, openGraphPrefix) {
					ogProp := strings.TrimPrefix(prop, openGraphPrefix)
					props[ogProp] = cont
				} else if _, contains := props["description"]; !contains && name == "description" {
					props[name] = cont
				}
				// if we dont have an image tag yet, check the link icon tags
			} else if _, contains := props["image"]; !contains && token.Data == "link" {
				// store the rel and href of the link tags
				var rel, href string
				for _, a := range token.Attr {
					switch a.Key {
					case "rel":
						rel = a.Val
					case "href":
						href = a.Val
					}
				}
				// if we found a favicon in one of the link rel properties
				if rel == "icon" || rel == "shortcut icon" {
					url, _ := url.Parse(href)
					if !url.IsAbs() {
						urlSt := path.Join(URL, url.String())
						props["image"] = urlSt
					} else {
						urlSt := path.Join(URL, url.String())
						props["image"] = urlSt
					}
				}
			}
		}
	}

	// if we still didn't get an image, check the root dir for a favicon
	if _, contains := props["image"]; !contains {
		favicon := URL + "/favicon.ico"
		_, err := http.Get(favicon)
		if err == nil {
			props["image"] = favicon
		}
	}

	return props, nil
}

// getPageSummary fetches a webpage and returns it's open graph properties summary
func getPageSummary(url string) (openGraphProps, error) {

	// fetch the HTML body
	body, err := fetchHTML(url)
	//ensure that the response body stream is closed eventually
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// tokenize and fetch the open graph properties from the url body
	props, err := fetchOpenGraphProps(body, url)
	if err != nil {
		return nil, err
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
	if len(URL) == 0 {
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
