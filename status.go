package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func status(w http.ResponseWriter, r *http.Request) {
	log.Printf("serving status to %s", r.RemoteAddr)

	now := time.Now()

	for _, md := range storage.Files {
		if md.Expire.After(now) {
			fmt.Fprintf(w, "%s %s %s\n", md.Hash, md.Expire.Sub(now), md.Filename)
		}
	}
}
