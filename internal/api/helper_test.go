package api_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func postForm(t *testing.T, h http.Handler, target string, data url.Values, bearer string) *httptest.ResponseRecorder {
	t.Helper()

	r := httptest.NewRequest(http.MethodPost, target, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if len(bearer) > 0 {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer))
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	t.Logf("response: %s", w.Body.String())
	return w
}

func get(t *testing.T, h http.Handler, target string, bearer string) *httptest.ResponseRecorder {
	t.Helper()

	r := httptest.NewRequest(http.MethodGet, target, nil)
	if len(bearer) > 0 {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer))
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	t.Logf("response: %s", w.Body.String())
	return w
}
