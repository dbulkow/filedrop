package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	DefaultPort   = "8433"
	DefaultListen = "127.0.0.1"
	DefaultCert   = "tls/cert.pem"
	DefaultKey    = "tls/key.pem"
)

var url string

func init() {
	var (
		listen = os.Getenv("FILEDROP_ADDRESS")
		port   = os.Getenv("FILEDROP_PORT")
		cert   = os.Getenv("FILEDROP_CERT")
		key    = os.Getenv("FILEDROP_KEY")
	)

	if listen == "" {
		listen = DefaultListen
	}

	if port == "" {
		port = DefaultPort
	}

	if cert == "" {
		cert = DefaultCert
	}

	if key == "" {
		key = DefaultKey
	}

	url = os.Getenv("FILEDROP_SERVER_URL")

	if url == "" {
		url = fmt.Sprintf("https://%s:%s", listen, port)
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start an https server",
		Long:  "Start an https server",
		Run:   serve,
	}

	serveCmd.Flags().StringVarP(&listen, "listen", "l", listen, "Listen address")
	serveCmd.Flags().StringVarP(&port, "port", "p", port, "Port number")
	serveCmd.Flags().StringVarP(&cert, "cert", "c", cert, "TLS cert file")
	serveCmd.Flags().StringVarP(&key, "key", "k", key, "TLS key file")
	serveCmd.Flags().StringVarP(&url, "url", "u", url, "Filedrop server URL to advertise")

	RootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	listen := cmd.Flag("listen").Value.String()
	port := cmd.Flag("port").Value.String()

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

	log.Printf("Listening on https://%s\n", srv.Addr)

	log.Fatal(srv.ListenAndServeTLS("tls/cert.pem", "tls/key.pem"))
}
