package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"encoding/json"

	"io"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
)

type testCase struct {
	method      string
	handler     http.HandlerFunc
	path        string
	body        interface{}
	expStatus   int
	expRespBody string
	jsonFlag    bool
	header      bool
	session     string
}

// create the handlers context for the tests
func NewContext() *Context {
	return &Context{
		SessionKey:   "supersecret",
		SessionStore: sessions.NewMemStore(-1),
		UserStore:    users.NewMemStore(),
	}
}

func newMongoContext() *Context {
	// Use the DBADDR to dial your MongoDB server
	dbAddr := os.Getenv("DBADDR")
	fmt.Printf("dialing mongo server at %s...\n", dbAddr)
	mongoSession, err := mgo.Dial(dbAddr)
	if err != nil {
		log.Fatalf("error dialing mongo: %v", err)
	}
	userStore, err := users.NewMongoStore(mongoSession, "production")
	if err != nil {
		log.Fatalf("error creating user store: %v", err)
	}
	return &Context{
		SessionKey:   "thisisasupersecretpasswordnobodyknowsit",
		SessionStore: sessions.NewMemStore(-1),
		UserStore:    userStore,
	}
}

func testCaseFunc(t *testing.T, c *testCase) {
	//fmt.Printf("Testing: %v\n", c.expRespBody)
	// defer wg.Done()

	// Create a request to pass to our handler.
	var body io.Reader
	if bod, ok := c.body.(string); ok {
		bodyStr := []byte(bod)
		body = bytes.NewBuffer(bodyStr)
	} else if c.body == nil {
		body = nil
	}

	req, err := http.NewRequest(c.method, c.path, body)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.handler)

	// add the session to the header if it was provided
	if c.session != "" {
		req.Header.Add("Authorization", c.session)
	}
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != c.expStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, c.expStatus)
	}

	if c.header {
		auth := rr.Header().Get("Authorization")
		if len(auth) == 0 {
			t.Errorf("handler didn't add auth header: got %v",
				auth)
		}
	} else {
		auth := rr.Header().Get("Authorization")
		if len(auth) != 0 {
			t.Errorf("handler added auth header, but shouldn't: got %v",
				auth)
		}
	}

	// Check the response body is what we expect.
	// MAKE SURE TO PUT A FUCKIN NEWLINE bcs encoder.Encode() writes a newline after every entry
	// lol
	var bodyRespStr string
	if c.jsonFlag {
		return
		buf := new(bytes.Buffer)
		if err := json.Compact(buf, []byte(c.expRespBody)); err != nil {
			t.Fatal(err)
		}
		bodyRespStr = buf.String() + "\n"
	} else {
		bodyRespStr = c.expRespBody + "\n"
	}

	// expBodyBytes := []byte(c.expRespBody)
	// expBodyStr := string(expBodyBytes)
	if rr.Body.String() != bodyRespStr {
		// fmt.Printf("exp: %v, %v\n", []byte(c.expRespBody), len(c.expRespBody))
		// fmt.Printf("got: %v, %v\n", []byte(rr.Body.String()), len(rr.Body.String()))
		t.Errorf("handler returned unexpected body: \ngot %v \nwant %v",
			rr.Body.String(), bodyRespStr)
	}

}

func TestUsersPOST(t *testing.T) {
	hctx := NewContext()
	// ----- test POST (sign up) -----
	cases := []testCase{
		//test invalid JSON
		testCase{
			method:      "POST",
			handler:     hctx.UsersHandler,
			path:        apiUsers,
			body:        `sdfjjsjdkkjlasldkfjllkl`,
			expStatus:   http.StatusBadRequest,
			expRespBody: "invalid JSON",
			jsonFlag:    false,
		},
		// TODO: test invalid email figure out what the response is
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiUsers,
			body: `{
				"email": "invalid",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "jim1111",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus:   http.StatusBadRequest,
			expRespBody: "error validating user: mail: missing phrase",
			jsonFlag:    false,
		},
		//test differing passConf
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiUsers,
			body: `{
				"email": "test@gmail.com",
				"password": "test1234",
				"passwordConf": "not1234",
				"userName": "jim1111",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus:   http.StatusBadRequest,
			expRespBody: "error validating user: Error: password and passwordConf don't match",
			jsonFlag:    false,
		},
		// test zero len username
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiUsers,
			body: `{
				"email": "test@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus:   http.StatusBadRequest,
			expRespBody: `error validating user: Error: username is zero length`,
			jsonFlag:    false,
		},
		// test < 6 password
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiUsers,
			body: `{
				"email": "test@gmail.com",
				"password": "test",
				"passwordConf": "test",
				"userName": "jim",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus:   http.StatusBadRequest,
			expRespBody: `error validating user: Error: password less than 6 chars`,
			jsonFlag:    false,
		},
		// make sure a valid new user can sign up
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiUsers,
			body: `{
				"email": "test@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "jim",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus: http.StatusOK,
			expRespBody: `{
				"email": "test@gmail.com",
				"userName": "jim",
				"firstName": "jimmy",
				"lastName": "jones",
				"photoURL":"https://www.gravatar.com/avatar/1aedb8d9dc4751e229a335e371db8058"
			}`,
			jsonFlag: true,
			header:   true,
		},
		// make sure a user can't sign up twice
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiUsers,
			body: `{
				"email": "test@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "jim",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus:   http.StatusBadRequest,
			expRespBody: `Error: email already exists in database`,
			jsonFlag:    false,
		},
		// using same email
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiUsers,
			body: `{
				"email": "test@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "notjim",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus:   http.StatusBadRequest,
			expRespBody: `Error: email already exists in database`,
			jsonFlag:    false,
		},
		// using same username
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiUsers,
			body: `{
				"email": "nottest@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "jim",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus:   http.StatusBadRequest,
			expRespBody: `Error: user name already exists in database`,
			jsonFlag:    false,
		},
	}

	for _, c := range cases {
		//fmt.Println("testing", c.expRespBody)
		testCaseFunc(t, &c)
	}

}

func TestUsersGET(t *testing.T) {
	// ----- test GET (get all users) -----
	hctx := NewContext()

	// add two valid users and verify that the response is valid
	cases := []testCase{
		// test case where there are no entries
		testCase{
			method:      "GET",
			handler:     hctx.UsersHandler,
			path:        apiSessions,
			body:        nil,
			expStatus:   http.StatusOK,
			expRespBody: "[]",
			jsonFlag:    false,
		},
		// test valid
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiSessions,
			body: `{
				"email": "real@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "test",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus: http.StatusOK,
			expRespBody: `{
				"email": "real@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "test",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			jsonFlag: true,
			header:   true,
		},
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiSessions,
			body: `{
				"email": "test@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "test2",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus: http.StatusOK,
			expRespBody: `{
				"email": "test@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "test2",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			jsonFlag: true,
			header:   true,
		},
		testCase{
			method:    "GET",
			handler:   hctx.UsersHandler,
			path:      apiSessions,
			body:      nil,
			expStatus: http.StatusOK,
			expRespBody: `[
				{
				"email": "real@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "test",
				"firstName": "jimmy",
				"lastName": "jones"
				},
				{
				"email": "test@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "test2",
				"firstName": "jimmy",
				"lastName": "jones"
			}
			]`,
			jsonFlag: true,
		},
	}

	for _, c := range cases {
		//fmt.Println("testing", c.expRespBody)
		testCaseFunc(t, &c)
	}
}

func TestSessionshandler(t *testing.T) {
	hctx := NewContext()
	cases := []testCase{
		// test NOT POST case
		testCase{
			method:      "GET",
			handler:     hctx.SessionsHandler,
			path:        apiSessions,
			body:        nil,
			expStatus:   http.StatusMethodNotAllowed,
			expRespBody: "request method must be POST",
			jsonFlag:    false,
		},
		//test invalid JSON
		testCase{
			method:      "POST",
			handler:     hctx.SessionsHandler,
			path:        apiSessions,
			body:        `sdfjjsjdkkjlasldkfjllkl`,
			expStatus:   http.StatusBadRequest,
			expRespBody: "invalid JSON",
			jsonFlag:    false,
		},
		// test valid user and credentials
		// create new user and then attempt to sign in
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiUsers,
			body: `{
				"email": "real@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "test",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus:   http.StatusOK,
			expRespBody: "",
			jsonFlag:    true,
			header:      true,
		},
		// then sign in
		testCase{
			method:  "POST",
			handler: hctx.SessionsHandler,
			path:    apiSessions,
			body: `{
				"email": "real@gmail.com",
				"password": "test1234"
				}`,
			expStatus:   http.StatusOK,
			expRespBody: "",
			jsonFlag:    true,
			header:      true,
		},
		// then sign in with bad password!
		testCase{
			method:  "POST",
			handler: hctx.SessionsHandler,
			path:    apiSessions,
			body: `{
				"email": "real@gmail.com",
				"password": "BADPASSWORD"
				}`,
			expStatus:   http.StatusUnauthorized,
			expRespBody: "error authenticating user",
			jsonFlag:    false,
		},
		// test sign in with wrong email
		testCase{
			method:  "POST",
			handler: hctx.SessionsHandler,
			path:    apiSessions,
			body: `{
				"email": "notReal@gmail.com",
				"password": "BADPASSWORD"
				}`,
			expStatus:   http.StatusUnauthorized,
			expRespBody: "error authenticating user",
			jsonFlag:    false,
		},
	}

	for _, c := range cases {
		testCaseFunc(t, &c)
	}
}

func TestSessionsHeaders(t *testing.T) {
	hctx := NewContext()
	cases := []testCase{
		// test NOT DELETE case
		testCase{
			method:      "GET",
			handler:     hctx.SessionsMineHandler,
			path:        apiSessionsMine,
			body:        nil,
			expStatus:   http.StatusMethodNotAllowed,
			expRespBody: "request method must be DELETE",
			jsonFlag:    false,
		},
		// add a valid user so we can test deleting their session
		testCase{
			method:  "POST",
			handler: hctx.UsersHandler,
			path:    apiUsers,
			body: `{
				"email": "real@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "test",
				"firstName": "jimmy",
				"lastName": "jones"
				}`,
			expStatus:   http.StatusOK,
			expRespBody: "",
			jsonFlag:    true,
			header:      true,
		},
		// test deleting with no sessionID in header
		testCase{
			method:      "DELETE",
			handler:     hctx.SessionsMineHandler,
			path:        apiSessionsMine,
			body:        nil,
			expStatus:   http.StatusBadRequest,
			expRespBody: "error getting sessionID: no session ID found in Authorization header",
			jsonFlag:    false,
		},
		// test invalid session scheme
		testCase{
			method:      "DELETE",
			handler:     hctx.SessionsMineHandler,
			path:        apiSessionsMine,
			body:        nil,
			expStatus:   http.StatusBadRequest,
			expRespBody: "error getting sessionID: scheme used in Authorization header is not supported",
			jsonFlag:    false,
			session:     "thisisobviouslygarbage",
		},
		// test invalid session
		testCase{
			method:      "DELETE",
			handler:     hctx.SessionsMineHandler,
			path:        apiSessionsMine,
			body:        nil,
			expStatus:   http.StatusBadRequest,
			expRespBody: "error getting sessionID: Invalid Session ID",
			jsonFlag:    false,
			// this was a previously generated bearer token
			session: "Bearer d2l6FFNU7aiZrF60gdvsQNetRkf9SjYoJdW5ll1qlILZ3eiW24DR6v46tsvx99cWsINNg1b6dYRmdQACwMMZbw==",
		},
	}
	for _, c := range cases {
		testCaseFunc(t, &c)
	}
}

func TestSessionsMine(t *testing.T) {
	hctx := NewContext()

	handler := http.HandlerFunc(hctx.UsersHandler)
	rr := httptest.NewRecorder()
	// add a new user to the userstore and get the auth token
	body := `{
				"email": "test@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "jim",
				"firstName": "jimmy",
				"lastName": "jones"
			}`

	bodyStr := []byte(body)
	req, err := http.NewRequest("POST", apiUsers, bytes.NewBuffer(bodyStr))
	if nil != err {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// get the sessionID from the header
	// and attempt to delete the session
	auth := rr.Header().Get("Authorization")
	handler = http.HandlerFunc(hctx.SessionsMineHandler)
	rr = httptest.NewRecorder()

	req, err = http.NewRequest("DELETE", apiSessionsMine, nil)
	if nil != err {
		t.Fatal(err)
	}
	// add the auth to the header
	req.Header.Add("Authorization", auth)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// check the body
	if rr.Body.String() != "user signed out\n" {
		t.Errorf("handler returned unexpected body: \ngot %v \nwant %v",
			rr.Body.String(), "user signed out\n")
	}

}

func TestUsersMe(t *testing.T) {
	hctx := NewContext()

	handler := http.HandlerFunc(hctx.UsersHandler)
	rr := httptest.NewRecorder()
	// add a new user to the userstore and get the auth token
	body := `{
				"email": "test@gmail.com",
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "jim",
				"firstName": "jimmy",
				"lastName": "jones"
			}`

	bodyStr := []byte(body)
	req, err := http.NewRequest("POST", apiUsers, bytes.NewBuffer(bodyStr))
	if nil != err {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// get the sessionID from the header
	auth := rr.Header().Get("Authorization")

	// then check to see if we can access the /me
	handler = http.HandlerFunc(hctx.UsersMeHanlder)
	rr = httptest.NewRecorder()

	req, err = http.NewRequest("GET", apiUsersMe, nil)
	if nil != err {
		t.Fatal(err)
	}
	// add the auth to the header
	req.Header.Add("Authorization", auth)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// check the body
	if rr.Body.String() != "user signed out\n" {
		t.Errorf("handler returned unexpected body: \ngot %v \nwant %v",
			rr.Body.String(), "user signed out\n")
	}

}
