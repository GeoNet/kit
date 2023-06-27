package weft_test

import (
	"errors"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/GeoNet/kit/weft"
)

func TestStatus(t *testing.T) {
	in := []struct {
		id     string
		err    error
		status int
	}{
		{id: loc(), err: nil, status: http.StatusOK},
		{id: loc(), err: errors.New("an error"), status: http.StatusServiceUnavailable},
		{id: loc(), err: weft.StatusError{Err: errors.New("an error")}, status: http.StatusServiceUnavailable},
		{id: loc(), err: weft.StatusError{Code: http.StatusServiceUnavailable, Err: errors.New("an error")}, status: http.StatusServiceUnavailable},
		{id: loc(), err: weft.StatusError{Code: http.StatusServiceUnavailable}, status: http.StatusServiceUnavailable},
		{id: loc(), err: weft.StatusError{Code: http.StatusBadRequest}, status: http.StatusBadRequest},
		{id: loc(), err: weft.StatusError{Code: http.StatusBadRequest, Err: errors.New("error for the client")}, status: http.StatusBadRequest},
		{id: loc(), err: weft.StatusError{Code: http.StatusMethodNotAllowed}, status: http.StatusMethodNotAllowed},
	}

	for _, v := range in {
		if s := weft.Status(v.err); s != v.status {
			t.Errorf("%s expected status %d got %d", v.id, v.status, s)
		}
	}
}

func TestCheckQuery(t *testing.T) {
	in := []struct {
		id                 string
		url                string
		required, optional []string
		status             int
	}{
		{id: loc(), url: "http://test.com", required: []string{}, optional: []string{}, status: http.StatusOK},
		{id: loc(), url: "http://test.com", required: []string{}, optional: []string{"optional"}, status: http.StatusOK},
		{id: loc(), url: "http://test.com?optional=t", required: []string{}, optional: []string{"optional"}, status: http.StatusOK},
		{id: loc(), url: "http://test.com?optional=t", required: []string{}, optional: []string{"optional", "another"}, status: http.StatusOK},
		{id: loc(), url: "http://test.com", required: []string{"missing"}, optional: []string{"optional"}, status: http.StatusBadRequest},
		{id: loc(), url: "http://test.com", required: []string{"missing"}, optional: []string{}, status: http.StatusBadRequest},
		{id: loc(), url: "http://test.com?required=t", required: []string{"required"}, optional: []string{}, status: http.StatusOK},
		{id: loc(), url: "http://test.com?required=t", required: []string{"required"}, optional: []string{"optional"}, status: http.StatusOK},
		{id: loc(), url: "http://test.com?required=t&optional=t", required: []string{"required"}, optional: []string{"optional"}, status: http.StatusOK},
		{id: loc(), url: "http://test.com?required=t&optional=t&extra=t", required: []string{"required"}, optional: []string{"optional"}, status: http.StatusBadRequest},
	}

	for _, v := range in {
		r, err := http.NewRequest("GET", v.url, nil)
		if err != nil {
			t.Errorf("%s parsing URL: %s", v.id, err.Error())
		}

		if s := weft.Status(weft.CheckQuery(r, []string{"GET"}, v.required, v.optional)); s != v.status {
			t.Errorf("%s expected status %d got %d", v.id, v.status, s)
		}

		if s := weft.Status(weft.CheckQuery(r, []string{"POST"}, v.required, v.optional)); s != http.StatusMethodNotAllowed {
			t.Errorf("%s expected status %d got %d", v.id, http.StatusMethodNotAllowed, s)
		}

		// test with a cache buster added.
		if !strings.Contains(v.url, "?") {
			v.url = v.url + "?"
		}

		r, err = http.NewRequest("GET", v.url+";busta", nil)
		if err != nil {
			t.Errorf("%s parsing buster URL: %s", v.id, err.Error())
		}

		if s := weft.Status(weft.CheckQuery(r, []string{"GET"}, v.required, v.optional)); s != http.StatusBadRequest {
			t.Errorf("%s buster expected status %d got %d", v.id, http.StatusBadRequest, s)
		}

		if s := weft.Status(weft.CheckQuery(r, []string{"POST"}, v.required, v.optional)); s != http.StatusMethodNotAllowed {
			t.Errorf("%s buster expected status %d got %d", v.id, http.StatusMethodNotAllowed, s)
		}
	}
}

func loc() string {
	_, _, l, _ := runtime.Caller(1)
	return "L" + strconv.Itoa(l)
}
