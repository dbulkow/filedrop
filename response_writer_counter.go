package main

import (
	"net/http"
)

type ResponseWriterCounter struct {
	http.ResponseWriter
	writer http.ResponseWriter
}

func NewResponseWriterCounter(rw http.ResponseWriter) *ResponseWriterCounter {
	return &ResponseWriterCounter{writer: rw}
}

func (c *ResponseWriterCounter) Write(buf []byte) (int, error) {
	n, err := c.writer.Write(buf)
	downloadBytes.Add(float64(n))
	return n, err
}

func (c *ResponseWriterCounter) Header() http.Header {
	return c.writer.Header()
}

func (c *ResponseWriterCounter) WriteHeader(code int) {
	c.writer.WriteHeader(code)
}
