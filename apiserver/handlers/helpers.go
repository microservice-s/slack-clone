package handlers

import (
	"encoding/json"
	"net/http"
)

// Respond writes data to a responseWriter
func Respond(w http.ResponseWriter, data interface{}, contentType string) {
	// add the header and encode the data as json
	w.Header().Add(headerContentType, contentTypeJSONUTF8)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}
