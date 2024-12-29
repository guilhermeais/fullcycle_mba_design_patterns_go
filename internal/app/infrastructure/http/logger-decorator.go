package http_handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type LoggerDecorator struct {
	Decoratee http.Handler
}

func (ld LoggerDecorator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	log.Printf("[%s %s]: \nHost: '%s'\nUser-Agent: '%s'\nBody: %s", r.Method, r.URL.Path, r.Host, r.Header.Get("user-agent"), string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	ld.Decoratee.ServeHTTP(w, r)
}
