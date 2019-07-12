package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	activeDirs = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "filedrop_active_dirs",
		Help: "Number of directories being served",
	})
	activeFiles = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "filedrop_active_files",
		Help: "Number of files being served",
	})
)

func init() {
	prometheus.MustRegister(uploads)
	prometheus.MustRegister(uploadBytes)
	prometheus.MustRegister(downloads)
	prometheus.MustRegister(downloadBytes)
	prometheus.MustRegister(serveFrontPage)
	prometheus.MustRegister(serveFavIcon)
	prometheus.MustRegister(activeDirs)
	prometheus.MustRegister(activeFiles)
}

func metrics() http.Handler {
	return promhttp.Handler()
}
