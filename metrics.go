package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	uploads = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "filedrop_uploads",
		Help: "Current number of uploaded files",
	})
	uploadBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "filedrop_upload_bytes",
		Help: "Total number of upload bytes",
	})
	downloads = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "filedrop_downloads",
		Help: "Current number of downloaded files",
	})
	downloadBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "filedrop_download_bytes",
		Help: "Total number of download bytes",
	})
	serveFrontPage = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "filedrop_serve_front_page",
		Help: "Total times front page served",
	})
	serveFavIcon = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "filedrop_serve_fav_icon",
		Help: "Total times favicon.ico served",
	})
)

func init() {
	prometheus.MustRegister(uploads)
	prometheus.MustRegister(uploadBytes)
	prometheus.MustRegister(downloads)
	prometheus.MustRegister(downloadBytes)
	prometheus.MustRegister(serveFrontPage)
	prometheus.MustRegister(serveFavIcon)
}

func metrics() http.Handler {
	return prometheus.Handler()
}
