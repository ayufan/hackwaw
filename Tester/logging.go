package main

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"
)

type bytesWriter interface {
	Status() int
	Bytes() []byte
}

type loggingResponseWriter struct {
	rw      http.ResponseWriter
	bytes   bytes.Buffer
	status  int
	started time.Time
}

func newLoggingResponseWriter(rw http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		rw:      rw,
		started: time.Now(),
	}
}

func (l *loggingResponseWriter) Status() int {
	return l.status
}

func (l *loggingResponseWriter) Bytes() []byte {
	return l.bytes.Bytes()
}

func (l *loggingResponseWriter) Header() http.Header {
	return l.rw.Header()
}

func (l *loggingResponseWriter) Write(data []byte) (n int, err error) {
	if l.status == 0 {
		l.status = http.StatusOK
	}
	l.bytes.Write(data)
	n, err = l.rw.Write(data)
	return
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.rw.WriteHeader(status)
	if l.status == 0 {
		l.status = status
	}
}

func (l *loggingResponseWriter) Log(r *http.Request) {
	duration := time.Since(l.started)

	message := ""
	if l.status >= 400 {
		message = l.bytes.String()
		idx := bytes.IndexAny(l.Bytes(), "\r\n")
		if idx >= 0 {
			message = "message: " + string(l.Bytes()[:idx])
		}
	}

	log.Println("API:", r.Method, r.URL,
		"http-status:", l.status,
		"bytes:", l.bytes.Len(),
		"duration:", strconv.FormatFloat(duration.Seconds(), 'f', 2, 64),
		message)
}
