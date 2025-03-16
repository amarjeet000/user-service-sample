package testutils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

type Header struct {
	Name  string
	Value string
}

func makeRequest(router *chi.Mux, method string, url string, headers []Header, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, body)
	for _, header := range headers {
		req.Header.Add(header.Name, header.Value)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func MakeGetRequestWithHeaders(router *chi.Mux, url string, headers []Header, body []byte) *httptest.ResponseRecorder {
	return makeRequest(router, http.MethodGet, url, headers, bytes.NewReader(body))
}
