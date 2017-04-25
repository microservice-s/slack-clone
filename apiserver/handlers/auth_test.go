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

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
	"github.com/aethanol/challenges-aethanol/apiserver/sessions"
)

type usersTestCase struct {
	method      string
	body        string
	expStatus   int
	expRespBody string
	jsonFlag    bool
}

// create the handlers context for the tests
func newContext() *Context {
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
		SessionKey:   "supersecret",
		SessionStore: sessions.NewMemStore(-1),
		UserStore:    userStore,
	}
}

func testPOSTUsersCase(t *testing.T, hctx *Context, c *usersTestCase) {
	//fmt.Printf("Testing: %v\n", c.expRespBody)
	// defer wg.Done()

	// Create a request to pass to our handler.
	bodyStr := []byte(c.body)
	req, err := http.NewRequest(c.method, apiUsers, bytes.NewBuffer(bodyStr))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hctx.UsersHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != c.expStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, c.expStatus)
	}

	// Check the response body is what we expect.
	// MAKE SURE TO PUT A FUCKIN NEWLINE bcs encoder.Encode() writes a newline after every entry
	// lol
	var bodyRespStr string
	if c.jsonFlag {
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

func TestUsersHandler(t *testing.T) {

	// ----- test POST (sign up) -----
	cases := []usersTestCase{
		//test invalid JSON
		usersTestCase{
			method:      "POST",
			body:        `sdfjjsjdkkjlasldkfjllkl`,
			expStatus:   http.StatusBadRequest,
			expRespBody: "invalid JSON",
			jsonFlag:    false,
		},
		//test invalid email TODO figure out what the response is
		usersTestCase{
			method: "POST",
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
		usersTestCase{
			method: "POST",
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
		usersTestCase{
			method: "POST",
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
		usersTestCase{
			method: "POST",
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
		usersTestCase{
			method: "POST",
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
				"password": "test1234",
				"passwordConf": "test1234",
				"userName": "jim",
				"firstName": "jimmy",
				"lastName": "jones"
			}`,
			jsonFlag: true,
		},
		// make sure a user can't sign up twice
		usersTestCase{
			method: "POST",
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
		usersTestCase{
			method: "POST",
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
		usersTestCase{
			method: "POST",
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
	hctx := newContext()
	for _, c := range cases {
		fmt.Println("testing", c.expRespBody)
		testPOSTUsersCase(t, hctx, &c)
	}

	// TEST VALIDATE

	// ----- test GET (get all users) -----

	// // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// // pass 'nil' as the third parameter.
	// req, err := http.NewRequest("GET", apiUsers, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	// rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(hctx.UsersHandler)

	// // Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// // directly and pass in our Request and ResponseRecorder.
	// handler.ServeHTTP(rr, req)

	// // Check the status code is what we expect.
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusOK)
	// }

	// // Check the response body is what we expect.
	// // MAKE SURE TO PUT A FUCKIN NEWLINE bcs encoder.Encode() writes a newline after every entry
	// // lol
	// expected := "[]\n"
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }

}

func TestSessionshandler(t *testing.T) {

}

func TestSessionsMineHandler(t *testing.T) {

}

func TestUsersMeHanlder(t *testing.T) {

}
