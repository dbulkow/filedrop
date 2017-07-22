package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func metrics() http.Handler {
	return prometheus.Handler()
}
