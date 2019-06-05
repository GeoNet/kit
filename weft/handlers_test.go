package weft_test

import (
	"bytes"
	"errors"
	"github.com/GeoNet/kit/weft"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeHandler(t *testing.T) {
	in := []struct {
		id        string
		f         weft.RequestHandler
		code      int
		surrogate string
	}{
		{
			id: loc(),
			f: func(r *http.Request, h http.Header, b *bytes.Buffer) error {
				return nil
			},
			code: http.StatusOK,
		},
		// returning an error will result in a 503
		{
			id: loc(),
			f: func(r *http.Request, h http.Header, b *bytes.Buffer) error {
				return errors.New("some error")
			},
			code:      http.StatusServiceUnavailable,
			surrogate: "max-age=10",
		},
		// an explicit 503 can also be returned
		{
			id: loc(),
			f: func(r *http.Request, h http.Header, b *bytes.Buffer) error {
				return weft.StatusError{Code: http.StatusServiceUnavailable, Err: errors.New("some error")}
			},
			code:      http.StatusServiceUnavailable,
			surrogate: "max-age=10",
		},
		// 500
		{
			id: loc(),
			f: func(r *http.Request, h http.Header, b *bytes.Buffer) error {
				return weft.StatusError{Code: http.StatusInternalServerError, Err: errors.New("some error")}
			},
			code:      http.StatusInternalServerError,
			surrogate: "max-age=10",
		},
		// 404 - no StatusError.Err needed
		{
			id: loc(),
			f: func(r *http.Request, h http.Header, b *bytes.Buffer) error {
				return weft.StatusError{Code: http.StatusNotFound}
			},
			code:      http.StatusNotFound,
			surrogate: "max-age=10",
		},
		{
			id: loc(),
			f: func(r *http.Request, h http.Header, b *bytes.Buffer) error {
				return weft.StatusError{Code: http.StatusMovedPermanently}
			},
			code:      http.StatusMovedPermanently,
			surrogate: "max-age=86400",
		},
		{
			id: loc(),
			f: func(r *http.Request, h http.Header, b *bytes.Buffer) error {
				return weft.StatusError{Code: http.StatusMovedPermanently}
			},
			code:      http.StatusMovedPermanently,
			surrogate: "max-age=86400",
		},
		{
			id: loc(),
			f: func(r *http.Request, h http.Header, b *bytes.Buffer) error {
				return weft.StatusError{Code: http.StatusGone}
			},
			code:      http.StatusGone,
			surrogate: "max-age=86400",
		},
		{
			id: loc(),
			f: func(r *http.Request, h http.Header, b *bytes.Buffer) error {
				return weft.StatusError{Code: http.StatusBadRequest}
			},
			code:      http.StatusBadRequest,
			surrogate: "max-age=86400",
		},
		{
			id: loc(),
			f: func(r *http.Request, h http.Header, b *bytes.Buffer) error {
				return weft.StatusError{Code: http.StatusMethodNotAllowed}
			},
			code:      http.StatusMethodNotAllowed,
			surrogate: "max-age=86400",
		},
	}

	// The TextError handler
	for _, v := range in {
		handler := weft.MakeHandler(v.f, weft.TextError)

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != v.code {
			t.Errorf("%s expected %d got %d", v.id, v.code, w.Code)
		}

		// only check the content and surrogate for codes that should change the content
		switch v.code {
		case http.StatusOK, http.StatusNoContent, http.StatusMovedPermanently, http.StatusSeeOther:
		default:
			if w.Code != http.StatusOK {
				if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
					t.Errorf("%s expected Content-Type %s got %s", v.id, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
				}

				if w.Header().Get("Surrogate-Control") != v.surrogate {
					t.Errorf("%s expected Surrogate-Control %s got %s", v.id, v.surrogate, w.Header().Get("Surrogate-Control"))
				}
			}
		}
	}

	// The HTMLError handler
	for _, v := range in {
		handler := weft.MakeHandler(v.f, weft.HTMLError)

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != v.code {
			t.Errorf("%s expected %d got %d", v.id, v.code, w.Code)
		}

		// only check the content and surrogate for codes that should change the content
		switch v.code {
		case http.StatusOK, http.StatusNoContent, http.StatusMovedPermanently, http.StatusSeeOther:
		default:
			if w.Code != http.StatusOK {
				if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
					t.Errorf("%s: expected Content-Type %s got %s", v.id, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
				}

				if w.Header().Get("Surrogate-Control") != v.surrogate {
					t.Errorf("%s: expected Surrogate-Control %s got %s", v.id, v.surrogate, w.Header().Get("Surrogate-Control"))
				}
			}
		}
	}

	// The UseError handler
	for _, v := range in {
		handler := weft.MakeHandler(v.f, weft.UseError)

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != v.code {
			t.Errorf("%s expected %d got %d", v.id, v.code, w.Code)
		}

		// only check the surrogate for codes that should change the content
		switch v.code {
		case http.StatusOK, http.StatusNoContent, http.StatusMovedPermanently, http.StatusSeeOther:
		default:
			if w.Code != http.StatusOK {
				if w.Header().Get("Surrogate-Control") != v.surrogate {
					t.Errorf("%s: expected Surrogate-Control %s got %s", v.id, v.surrogate, w.Header().Get("Surrogate-Control"))
				}
			}
		}
	}
}

func TestMakeDirectHandler(t *testing.T) {
	in := []struct {
		id        string
		f         weft.DirectRequestHandler
		code      int
		surrogate string
	}{
		{
			id: loc(),
			f: func(r *http.Request, w http.ResponseWriter) (int64, error) {
				n, err := w.Write([]byte{})
				return int64(n), err
			},
			code: http.StatusOK,
		},
		// returning an error will result in a 503
		{
			id: loc(),
			f: func(r *http.Request, w http.ResponseWriter) (int64, error) {
				return 0, errors.New("some error")
			},
			code:      http.StatusServiceUnavailable,
			surrogate: "max-age=10",
		},
		// an explicit 503 can also be returned
		{
			id: loc(),
			f: func(r *http.Request, w http.ResponseWriter) (int64, error) {
				return 0, weft.StatusError{Code: http.StatusServiceUnavailable, Err: errors.New("some error")}
			},
			code:      http.StatusServiceUnavailable,
			surrogate: "max-age=10",
		},
		// 500
		{
			id: loc(),
			f: func(r *http.Request, w http.ResponseWriter) (int64, error) {
				return 0, weft.StatusError{Code: http.StatusInternalServerError, Err: errors.New("some error")}
			},
			code:      http.StatusInternalServerError,
			surrogate: "max-age=10",
		},
		// 404 - no StatusError.Err needed
		{
			id: loc(),
			f: func(r *http.Request, w http.ResponseWriter) (int64, error) {
				return 0, weft.StatusError{Code: http.StatusNotFound}
			},
			code:      http.StatusNotFound,
			surrogate: "max-age=10",
		},
		{
			id: loc(),
			f: func(r *http.Request, w http.ResponseWriter) (int64, error) {
				return 0, weft.StatusError{Code: http.StatusMovedPermanently}
			},
			code:      http.StatusMovedPermanently,
			surrogate: "max-age=86400",
		},
		{
			id: loc(),
			f: func(r *http.Request, w http.ResponseWriter) (int64, error) {
				return 0, weft.StatusError{Code: http.StatusMovedPermanently}
			},
			code:      http.StatusMovedPermanently,
			surrogate: "max-age=86400",
		},
		{
			id: loc(),
			f: func(r *http.Request, w http.ResponseWriter) (int64, error) {
				return 0, weft.StatusError{Code: http.StatusGone}
			},
			code:      http.StatusGone,
			surrogate: "max-age=86400",
		},
		{
			id: loc(),
			f: func(r *http.Request, w http.ResponseWriter) (int64, error) {
				return 0, weft.StatusError{Code: http.StatusBadRequest}
			},
			code:      http.StatusBadRequest,
			surrogate: "max-age=86400",
		},
		{
			id: loc(),
			f: func(r *http.Request, w http.ResponseWriter) (int64, error) {
				return 0, weft.StatusError{Code: http.StatusMethodNotAllowed}
			},
			code:      http.StatusMethodNotAllowed,
			surrogate: "max-age=86400",
		},
	}

	// The TextError handler
	for _, v := range in {
		handler := weft.MakeDirectHandler(v.f, weft.TextError)

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != v.code {
			t.Errorf("%s: expected %d got %d", v.id, v.code, w.Code)
		}

		// only check the content and surrogate for codes that should change the content
		switch v.code {
		case http.StatusOK, http.StatusNoContent, http.StatusMovedPermanently, http.StatusSeeOther:
		default:
			if w.Code != http.StatusOK {
				if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
					t.Errorf("%s: expected Content-Type %s got %s", v.id, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
				}

				if w.Header().Get("Surrogate-Control") != v.surrogate {
					t.Errorf("%s: expected Surrogate-Control %s got %s", v.id, v.surrogate, w.Header().Get("Surrogate-Control"))
				}
			}
		}
	}

	// The HTMLError handler
	for _, v := range in {
		handler := weft.MakeDirectHandler(v.f, weft.HTMLError)

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != v.code {
			t.Errorf("%s: expected %d got %d", v.id, v.code, w.Code)
		}

		// only check the content and surrogate for codes that should change the content
		switch v.code {
		case http.StatusOK, http.StatusNoContent, http.StatusMovedPermanently, http.StatusSeeOther:
		default:
			if w.Code != http.StatusOK {
				if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
					t.Errorf("%s: expected Content-Type %s got %s", v.id, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
				}

				if w.Header().Get("Surrogate-Control") != v.surrogate {
					t.Errorf("%s: expected Surrogate-Control %s got %s", v.id, v.surrogate, w.Header().Get("Surrogate-Control"))
				}
			}
		}
	}
}

func TestGzip(t *testing.T) {
	fn := func(r *http.Request, h http.Header, b *bytes.Buffer) error {
		// write some content to test encoding sniffing and gzip
		b.Write([]byte(weft.ErrNotFound))

		return nil
	}

	handler := weft.MakeHandler(fn, weft.TextError)

	// no gzip

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Encoding") != "" {
		t.Error("did not expected encoded content without Accept-Encoding set")
	}

	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Error("incorrect content type")
	}

	// with gzip

	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Set("Accept-Encoding", weft.GZIP)
	req.Header.Set("Accept", "application/vnd.geo+json;version=2")
	w = httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Encoding") != weft.GZIP {
		t.Errorf("expected %s-encoded content, got %s", weft.GZIP, w.Header().Get("Content-Encoding"))
	}

	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Error("incorrect content type")
	}
}

func TestGzip2(t *testing.T) {
	fn := func(r *http.Request, h http.Header, b *bytes.Buffer) error {
		// write some content to test encoding sniffing and gzip
		b.Write([]byte(weft.ErrNotFound))

		h.Set("Content-Type", "application/vnd.geo+json;version=2")

		return nil
	}

	handler := weft.MakeHandler(fn, weft.TextError)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Set("Accept-Encoding", weft.GZIP)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Encoding") != weft.GZIP {
		t.Errorf("expected %s-encoded content, got %s", weft.GZIP, w.Header().Get("Content-Encoding"))
	}

	if w.Header().Get("Content-Type") != "application/vnd.geo+json;version=2" {
		t.Error("incorrect content type")
	}
}

func TestHandlers(t *testing.T) {
	in := []struct {
		id   string
		f    weft.RequestHandler
		code int
	}{
		{id: loc(), f: weft.NoMatch, code: http.StatusNotFound},
		{id: loc(), f: weft.Up, code: http.StatusOK},
		{id: loc(), f: weft.Soh, code: http.StatusOK},
	}

	for _, v := range in {
		handler := weft.MakeHandler(v.f, weft.TextError)

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != v.code {
			t.Errorf("expected %d got %d", v.code, w.Code)
		}

		req = httptest.NewRequest("POST", "http://example.com/foo", nil)
		w = httptest.NewRecorder()
		handler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected %d got %d", http.StatusMethodNotAllowed, w.Code)
		}
	}
}
