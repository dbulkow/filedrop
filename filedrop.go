package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var RootCmd = &cobra.Command{
	Use:   "filedrop",
	Short: "File sharing service",
	Long:  "File sharing service",
}

const (
	DefaultRoot   = "./downloads"
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

	if url == "" {
		url = fmt.Sprintf("https://%s:%s", listen, port)
	}

	if root == "" {
		root = DefaultRoot
	}

	flag.StringVarP(&root, "root", "r", root, "Storage directory")
	flag.StringVarP(&listen, "listen", "l", listen, "Listen address")
	flag.StringVarP(&port, "port", "p", port, "Port number")
	flag.StringVarP(&url, "url", "u", url, "Filedrop server URL to advertise")

	flag.Parse()

	storage = NewStorage(root)

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", makeGzipHandler(frontPage)))
	mux.HandleFunc("/status", status)
	mux.Handle("/metrics", metrics())

	srv := &http.Server{
		Addr:           listen + ":" + port,
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSNextProto:   nil,
	}

	log.Printf("Listening on http://%s\n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
