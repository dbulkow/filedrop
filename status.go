package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func status(w http.ResponseWriter, r *http.Request) {
	storage.Lock()
	defer storage.Unlock()

	log.Printf("serving status to %s", r.RemoteAddr)

	now := time.Now()

	from, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("split: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	for _, md := range storage.Dirs {
		if md.From == from {
			rem := md.Expire.Sub(now)
			for _, f := range md.Files {
				fmt.Fprintf(w, "%s %v %s %t\n", md.Hash, rem-(rem%time.Minute), f.Name, f.downloaded)
			}
		}
	}
}
