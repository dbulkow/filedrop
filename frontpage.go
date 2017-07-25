package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
)

//go:generate go run scripts/mkpage.go frontpage.html
//go:generate go run scripts/mkpage.go replypage.html

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

	if r.URL.Path != "" {
		storage.Lock()
		defer storage.Unlock()

		log.Printf("serving file %s to %s\n", r.URL.Path, r.RemoteAddr)

		md, ok := storage.Files[r.URL.Path]
		if ok == false {
			log.Printf("file not found %s", r.URL.Path)
			w.Header().Set("Content-Type", "text/html")
			return
		}

		downloads.Inc()
		writer := NewResponseWriterCounter(w)

		writer.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", md.Filename))
		http.ServeFile(writer, r, path.Join(storage.Root, r.URL.Path, md.Filename))
		return
	}

	log.Printf("serving form to %s", r.RemoteAddr)

	serveFrontPage.Inc()

	w.Header().Set("Etag", fmt.Sprintf("\"%s\"", frontpage_etag))
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "max-age=31536000") // 1 year
	w.Write([]byte(frontpage))
}
