package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/amartin3659/VacationHomeRental/internal/driver"
	"github.com/amartin3659/VacationHomeRental/internal/models"
)

type postData struct {
	key   string
	value string
}

var allTheHandlerTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"eremite", "/eremite", "GET", http.StatusOK},
	{"couple", "/couple", "GET", http.StatusOK},
	{"family", "/family", "GET", http.StatusOK},
	{"reservation", "/reservation", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"not-existing-route", "/not-existing-dummy", "GET", http.StatusNotFound},
}

//type DB struct {
//  SQL *sql.DB
//}

func TestNewRepo(t *testing.T) {
	var testdb driver.DB
	testRepo := NewRepo(&app, &testdb)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func TestAllTheHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range allTheHandlerTests {
		if test.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("%s: expected %d, got %d", test.name, test.expectedStatusCode, response.StatusCode)
			}
		}
	}
}

// TestRepository_MakeReservation tests the MakeReservation get-request handle
func TestRepository_MakeReservation(t *testing.T) {

	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservaton", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// rr means "request recorder" and is an initialized response recorder for http requests builtin the test
	// basically "faking" a client and to provide a valid request/response-cycle during tests
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	// turning handler into a function
	handler := http.HandlerFunc(Repo.MakeReservation)

	// calling handler  to test a function
	handler.ServeHTTP(rr, req)

	// the test itself as a condition
	if rr.Code != http.StatusOK {
		t.Errorf("handler MakeReservation failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case without a reservation in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler MakeReservation failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test error returned from database query function
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.BungalowID = 99
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler MakeReservation failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// TestRepository_PostMakeReservation tests the PostMakeReservation post-request handler
func TestRepository_PostMakeReservation(t *testing.T) {

	// case #1: reservation works fine

	postedData := url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout := "2006-01-02"
	sd, _ := time.Parse(layout, "2037-01-01")
	ed, _ := time.Parse(layout, "2037-01-02")
	bungalowId, _ := strconv.Atoi("1")

	reservation := models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case #2: missing post body

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #3: missing session data

	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("1")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code for missing session data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #4: invalid/insufficient data

	postedData = url.Values{}
	postedData.Add("full_name", "P")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("1")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("PostMakeReservation handler returned wrong response code invalid/insufficient data: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// case #5:  failure inserting reservation into database

	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("99")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to inserting a reservation into the database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #6: failure to inserting restriction into database

	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("999")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to inserting a reservation into the database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// TestRepository_ReservationJSON tests the ReservationJSON POST-request handler
func TestRepository_ReservationJSON(t *testing.T) {

	// -- test variables
	var j jsonResponse
	var postData url.Values
	var req *http.Request
	var ctx context.Context
	var rr *httptest.ResponseRecorder
	var handler http.Handler
	var err error

	// case #1: bungalow ID invalid

	// -- create request body
	postData = url.Values{}
	postData.Add("start", "2036-01-01")
	postData.Add("end", "2036-01-02")
	postData.Add("bungalow_id", "invalid")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation-json", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.ReservationJSON)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusTemporaryRedirect, rr.Code)
	}
	// case #2: start date invalid

	// -- create request body
	postData = url.Values{}
	postData.Add("start", "invalid")
	postData.Add("end", "2036-01-02")
	postData.Add("bungalow_id", "1")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation-json", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.ReservationJSON)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #3: end date invalid

	// -- create request body
	postData = url.Values{}
	postData.Add("start", "2036-01-02")
	postData.Add("end", "invalid")
	postData.Add("bungalow_id", "1")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation-json", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.ReservationJSON)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #4: bungalow not available
	// 2037-01-01

	// -- create request body
	postData = url.Values{}
	postData.Add("start", "2037-01-01")
	postData.Add("end", "2037-01-01")
	postData.Add("bungalow_id", "1")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation-json", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.ReservationJSON)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
	if j.OK {
		t.Errorf("Expected bungalow to be booked with response OK: %t, but got response OK: %t", false, j.OK)
	}

	// case #5: bungalow available

	// -- create request body
	postData = url.Values{}
	postData.Add("start", "2024-01-17")
	postData.Add("end", "2024-01-17")
	postData.Add("bungalow_id", "1")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation-json", strings.NewReader(postData.Encode()))
	// -- get context
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.ReservationJSON)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- parse json and recieve response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
	if !j.OK {
		t.Errorf("Expected OK value of %t, but got %t", true, j.OK)
	}

	// case #6: no request body

	// omit body
	// -- create request
	// have to use http.NewRequest, httptest.NewRequest doesn't trigger error when request.ParseForm is called
	req, _ = http.NewRequest("POST", "/reservation-json", nil)
	// -- get context
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.ReservationJSON)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- parse json and recieve response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
	if j.OK || j.Message != "Internal server error" {
		t.Errorf("Got response OK: %t, with empty request body", j.OK)
	}

	// case #7: database error
	// 2038-01-01
	postData = url.Values{}
	postData.Add("start", "2038-01-01")
	postData.Add("end", "2038-01-02")
	postData.Add("bungalow_id", "1")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation-json", strings.NewReader(postData.Encode()))
	// -- get context
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.ReservationJSON)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- parse json and recieve response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
	if j.OK {
		t.Errorf("Expected error thrown within ReservationJSON method, but got response OK: %t", j.OK)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

// PostReservation
func TestRepository_PostReservation(t *testing.T) {

	// -- test variables
	var postData url.Values
	var req *http.Request
	var ctx context.Context
	var rr *httptest.ResponseRecorder
	var handler http.Handler

	// case #1: invalid start date
	// -- create request body
	postData = url.Values{}
	postData.Add("start", "invalid")
	postData.Add("end", "2036-01-02")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostReservation)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #2: end date invalid
	// -- create request body
	postData = url.Values{}
	postData.Add("start", "2036-01-02")
	postData.Add("end", "invalid")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostReservation)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #3: bungalow not available
	// -- create request body
	postData = url.Values{}
	postData.Add("start", "2037-01-01")
	postData.Add("end", "2037-01-01")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostReservation)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status code: %d, but got status code: %d", http.StatusSeeOther, rr.Code)
	}

	// case #4: bungalow available
	// -- create request body
	postData = url.Values{}
	postData.Add("start", "2024-01-17")
	postData.Add("end", "2024-01-17")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation", strings.NewReader(postData.Encode()))
	// -- get context
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostReservation)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- parse json and recieve response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code: %d, but got status code: %d", http.StatusOK, rr.Code)
	}

	// case #5: no request body
	// -- create request
	// have to use http.NewRequest, httptest.NewRequest doesn't trigger error when request.ParseForm is called
	req, _ = http.NewRequest("POST", "/reservation", nil)
	// -- get context
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostReservation)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- parse json and recieve response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code: %d, but got status code: %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #6: database error
	postData = url.Values{}
	postData.Add("start", "2038-01-01")
	postData.Add("end", "2038-01-02")
	// -- create request
	req = httptest.NewRequest("POST", "/reservation", strings.NewReader(postData.Encode()))
	// -- get context
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostReservation)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- parse json and recieve response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code: %d, but got status code: %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

// ReservationOverview
func TestRepository_ReservationOverview(t *testing.T) {
	// case variables
	var req *http.Request
	var ctx context.Context
	var rr *httptest.ResponseRecorder
	var handler http.Handler

	// case #1: No session data
	// -- create request
	req = httptest.NewRequest("GET", "/reservation-overview", nil)
	// -- create ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.ReservationOverview)
	// -- make request
	handler.ServeHTTP(rr, req)
	// check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected response code: %d, got response code: %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #2: Invalid bungalow id
	// -- create reservation session data
	reservation := models.Reservation{
		BungalowID: 4,
		Bungalow: models.Bungalow{
			ID:           4,
			BungalowName: "The Solitude Shack",
		},
	}
	// -- create request
	req = httptest.NewRequest("GET", "/reservation-overview", nil)
	// -- create ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- add session data
	session.Put(ctx, "reservation", reservation)
	// -- create handler
	handler = http.HandlerFunc(Repo.ReservationOverview)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected response code: %d, got response code: %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #3: Nothing wrong

	// -- create reservation session data
	reservation = models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}
	// -- create request
	req = httptest.NewRequest("GET", "/reservation-overview", nil)
	// -- create ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- add session data
	session.Put(ctx, "reservation", reservation)
	// -- create handler
	handler = http.HandlerFunc(Repo.ReservationOverview)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected response code: %d, got response code: %d", http.StatusOK, rr.Code)
	}
}

// ChooseBungalow
func TestRepository_ChooseBungalow(t *testing.T) {
	// test variables
	var req *http.Request
	var ctx context.Context
	var rr *httptest.ResponseRecorder
	var handler http.Handler
	var reservation models.Reservation

	// case #1: Invalid id / cannot convert to int

	// -- create reservation session data
	reservation = models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}
	// -- create request
	req = httptest.NewRequest("GET", "/choose-bungalow/invalid", nil)
	// -- create ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
  // -- set request URI
  req.RequestURI = "/choose-bungalow/invalid"
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- add session data
	session.Put(ctx, "reservation", reservation)
	// -- create handler
	handler = http.HandlerFunc(Repo.ChooseBungalow)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected response code: %d, got response code: %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #2: No session

	// -- create request
	req = httptest.NewRequest("GET", "/choose-bungalow/1", nil)
	// -- create ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
  // -- set request URI
  req.RequestURI = "/choose-bungalow/1"
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.ChooseBungalow)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected response code: %d, got response code: %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #3: Nothing wrong

	// -- create reservation session data
	reservation = models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}
	// -- create request
	req = httptest.NewRequest("GET", "/choose-bungalow/1", nil)
	// -- create ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
  // -- add request URI
  req.RequestURI = "/choose-bungalow/1"
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- add session data
	session.Put(ctx, "reservation", reservation)
	// -- create handler
	handler = http.HandlerFunc(Repo.ChooseBungalow)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected response code: %d, got response code: %d", http.StatusSeeOther, rr.Code)
	}
}

// BookBungalow
func TestRepository_BookBungalow(t *testing.T) {

	// test variables
	var req *http.Request
	var ctx context.Context
	var rr *httptest.ResponseRecorder
	var handler http.Handler

	// case #1: No bungalow in db

	// -- create request
	req = httptest.NewRequest("GET", "/book-bungalow?s=2024-01-01&e=2024-01-02&id=4", nil)
	// -- create ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.BookBungalow)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected response code: %d, got response code: %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #2: OK

	// -- create request
	req = httptest.NewRequest("GET", "/book-bungalow?s=2024-01-01&e=2024-01-02&id=1", nil)
	// -- create ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.BookBungalow)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected response code: %d, got response code: %d", http.StatusSeeOther, rr.Code)
	}
}

// ShowLogin
func TestRepository_ShowLogin(t *testing.T) {
    
	// test variables
	var req *http.Request
	var ctx context.Context
	var rr *httptest.ResponseRecorder
	var handler http.Handler

	// case #1: OK

	// -- create request
	req = httptest.NewRequest("GET", "/user/login", nil)
	// -- create ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.ShowLogin)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected response code: %d, got response code: %d", http.StatusOK, rr.Code)
	}
}

// PostShowLogin
func TestRepository_PostShowLogin(t *testing.T) {

	// -- test variables
	var postData url.Values
	var req *http.Request
	var ctx context.Context
	var rr *httptest.ResponseRecorder
	var handler http.Handler

	// case #1: no form data in body
	// -- create request body
	// -- create request
	req, _ = http.NewRequest("POST", "/user/login", nil)
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostShowLogin)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #2: no email provided
	// -- create request body
	postData = url.Values{}
	postData.Add("password", "pass123")
	// -- create request
	req = httptest.NewRequest("POST", "/user/login", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostShowLogin)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
  error, ok := app.Session.Get(ctx, "error").(string)
  if !ok {
    t.Error("Error getting session info")
  }
	if error != "Invalid login credentials" {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #3: no password provided
	// -- create request body
	postData = url.Values{}
	postData.Add("email", "email@test.com")
	// -- create request
	req = httptest.NewRequest("POST", "/user/login", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostShowLogin)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
  error, ok = app.Session.Get(ctx, "error").(string)
  if !ok {
    t.Error("Error getting session info")
  }
	if error != "Invalid login credentials" {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #4: invalid email
	// -- create request body
	postData = url.Values{}
	postData.Add("email", "email")
  postData.Add("password", "pass123")
	// -- create request
	req = httptest.NewRequest("POST", "/user/login", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostShowLogin)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
  error, ok = app.Session.Get(ctx, "error").(string)
  if !ok {
    t.Error("Error getting session info")
  }
	if error != "Invalid login credentials" {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// case #5: not a user
	// -- create request body
	postData = url.Values{}
	postData.Add("email", "email@test.com")
  postData.Add("password", "pass123")
	// -- create request
	req = httptest.NewRequest("POST", "/user/login", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostShowLogin)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusTemporaryRedirect, rr.Code)
	}
  
	// case #6: OK
	// -- create request body
	postData = url.Values{}
	postData.Add("email", "validemail@test.com")
  postData.Add("password", "pass123")
	// -- create request
	req = httptest.NewRequest("POST", "/user/login", strings.NewReader(postData.Encode()))
	// -- get ctx
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// -- set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// -- create response recorder
	rr = httptest.NewRecorder()
	// -- create handler
	handler = http.HandlerFunc(Repo.PostShowLogin)
	// -- make request
	handler.ServeHTTP(rr, req)
	// -- check response
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusSeeOther, rr.Code)
	}
}
// Logout

// AdminDashboard

// AdminNewReservations

// AdminAllReservations

// AdminReservationsCalendar

// AdminShowReservation

// AdminPostShowReservation

// AdminProcessReservation

// AdminDeleteReservation

// AdminPostReservationsCalendar
