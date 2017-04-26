package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResetCodes(t *testing.T) {

	hctx := NewContext()
	handler := http.HandlerFunc(hctx.ResetCodesHandler)
	rr := httptest.NewRecorder()
	// add a new user to the userstore and get the auth token
	// body := `{
	// 			"email": "ethan.anderson6@gmail.com",
	// 		}`

	// bodyStr := []byte(body)
	req, err := http.NewRequest("POST", "/v1/resetcodes", nil)
	if nil != err {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

}
