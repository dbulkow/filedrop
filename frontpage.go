package main

import (
	"fmt"
	"log"
	"net/http"
)

//go:generate go run scripts/mkpage.go frontpage.html

func frontPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "favicon.ico" {
		log.Printf("serving favicon to %s", r.RemoteAddr)
		serveFavIcon.Inc()
		http.ServeFile(w, r, "favicon.png")
		return
	}

	log.Printf("serving form to %s", r.RemoteAddr)

	serveFrontPage.Inc()

	w.Header().Set("Etag", fmt.Sprintf("\"%s\"", frontpage_etag))
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "max-age=31536000") // 1 year
	w.Write([]byte(frontpage_html))
}
