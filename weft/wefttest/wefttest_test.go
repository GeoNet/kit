package wefttest_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/GeoNet/kit/weft"
	wt "github.com/GeoNet/kit/weft/wefttest"
)

const (
	errContent  = "text/plain; charset=utf-8"
	maxAge86400 = "max-age=86400"
)

//expected csp header for normal responses
var normalCspHeader = weft.ReturnDefaultCSP()

//expected csp header for error responses
var errorCspHeader = weft.ReturnDefaultCSP()

// test server and handlers for running the tests

var ts *httptest.Server

var routes = wt.Requests{
	{ID: wt.L(), URL: "/soh/up"},
	{ID: wt.L(), URL: "/pet/store?petType=dog"},
}

func setup() {
	mux := http.NewServeMux()
	mux.HandleFunc("/soh/up", weft.MakeHandler(weft.Up, weft.TextError))
	mux.HandleFunc("/", weft.MakeHandler(queryHandler, weft.TextError))
	ts = httptest.NewServer(mux)
}

func teardown() {
	ts.Close()
}

func queryHandler(r *http.Request, h http.Header, b *bytes.Buffer) error {
	_, err := weft.CheckQueryValid(r, []string{"GET"}, []string{"petType"}, []string{}, validator)
	if err != nil {
		return err
	}

	return nil
}

func validator(u url.Values) error {
	if u.Get("petType") != "dog" {
		return weft.StatusError{Code: http.StatusBadRequest, Err: fmt.Errorf("got unxpected petType: %s", u.Get("petType"))}
	}

	return nil
}

func TestAllRoutes(t *testing.T) {
	setup()
	defer teardown()

	err := routes.DoAll(ts.URL)
	if err != nil {
		t.Error(err)
	}
}

func TestRoutes(t *testing.T) {
	setup()
	defer teardown()

	success := 0
	errors := 0

	for _, v := range routes {
		v.CSP = normalCspHeader
		_, err := v.Do(ts.URL)
		if err != nil {
			t.Errorf("TestRoutes %s", err.Error())
			errors++
		} else {
			success++
		}
	}

	t.Logf("TestRoutes success: %d errors: %d", success, errors)
}

func TestMethodNotAllowed(t *testing.T) {
	setup()
	defer teardown()

	success := 0
	errors := 0

	for _, v := range routes {
		v.Surrogate = maxAge86400
		v.Content = errContent
		v.CSP = errorCspHeader //strictCsp for error response

		i, err := v.MethodNotAllowed(ts.URL, []string{"GET"})
		if err != nil {
			t.Errorf("TestMethodNotAllowed %s", err.Error())
			errors++
		} else {
			success += i
		}
	}

	t.Logf("TestMethodNotAllowed success: %d errors: %d", success, errors)
}

func TestExtraParameter(t *testing.T) {
	setup()
	defer teardown()

	success := 0
	errors := 0

	for _, v := range routes {
		v.Surrogate = maxAge86400
		v.Content = errContent
		v.CSP = errorCspHeader //strictCsp for error response

		err := v.ExtraParameter(ts.URL, "extra", "parameter")
		if err != nil {
			t.Errorf("TestExtraParameter %s", err.Error())
			errors++
		} else {
			success++
		}
	}

	t.Logf("TestExtraParameter success: %d errors: %d", success, errors)
}

// TestFuzzRoutes tests routes with fuzzed query parameters.
// Fuzzing takes a while to run.  Fuzz tests can be excluded during other testing with:
//    go test -v -run 'Test[^Fuzz]'
func TestFuzzQuery(t *testing.T) {
	setup()
	defer teardown()

	success := 0
	errors := 0

	for _, v := range routes {
		v.Surrogate = maxAge86400
		v.Content = errContent
		v.CSP = errorCspHeader //strictCsp for error response
		i, err := v.FuzzQuery(ts.URL, wt.FuzzValues)
		if err != nil {
			t.Errorf("TestFuzzQuery %s", err.Error())
			errors++
		} else {
			success = success + i
		}
	}

	t.Logf("TestFuzzQuery success: %d errors: %d", success, errors)
}

func TestFuzzPath(t *testing.T) {
	setup()
	defer teardown()

	success := 0
	errors := 0
	for _, v := range routes {
		// will 404 or 400 so can't be sure of cache or content types.  Exclude them.
		v.Surrogate = ""
		v.Content = ""
		v.CSP = nil //no csp check
		i, err := v.FuzzPath(ts.URL, wt.FuzzValues)
		if err != nil {
			t.Errorf("TestFuzzPath %s", err.Error())
			errors++
		} else {
			success = success + i
		}
	}

	t.Logf("TestFuzzPath success: %d errors: %d", success, errors)
}
