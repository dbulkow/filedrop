package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var RootCmd = &cobra.Command{
	Use:   "filedrop",
	Short: "File sharing service",
	Long:  "File sharing service",
}

const (
	DefaultPort   = "8080"
	DefaultListen = "127.0.0.1"
)

var (
	storage *Storage
	url     string
)

func main() {
	var (
		root   = os.Getenv("FILEDROP_ROOT")
		listen = os.Getenv("FILEDROP_ADDRESS")
		port   = os.Getenv("FILEDROP_PORT")
	)

	url = os.Getenv("FILEDROP_SERVER_URL")

	if listen == "" {
		listen = DefaultListen
	}

	if port == "" {
		port = DefaultPort
	}

	hostname, _ := os.Hostname()
	if url == "" {
		url = fmt.Sprintf("http://%s:%s", hostname, port)
	}

	flag.StringVarP(&root, "root", "r", root, "Storage directory")
	flag.StringVarP(&listen, "listen", "l", listen, "Listen address")
	flag.StringVarP(&port, "port", "p", port, "Port number")
	flag.StringVarP(&url, "url", "u", url, "Filedrop server URL to advertise")

	flag.Parse()

	storage = NewStorage(root)

	mux := http.NewServeMux()
	mux.Handle("/", prometheus.InstrumentHandler("server", http.StripPrefix("/", makeGzipHandler(frontPage))))
	mux.Handle("/postfile", prometheus.InstrumentHandler("postfile", http.StripPrefix("/postfile", makeGzipHandler(postfile))))
	mux.Handle("/status", prometheus.InstrumentHandler("status", makeGzipHandler(status)))
	mux.Handle("/metrics", prometheus.InstrumentHandler("metrics", metrics()))

	srv := &http.Server{
		Addr:           listen + ":" + port,
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSNextProto:   nil,
	}

	log.Printf("Advertising url \"%s\"", url)
	log.Printf("Listening on http://%s", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
