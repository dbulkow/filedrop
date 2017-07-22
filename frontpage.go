package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
)

//go:generate go run scripts/mkpage.go frontpage.html
//go:generate go run scripts/mkpage.go replypage.html

func frontPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "favicon.ico" {
		log.Printf("serving favicon to %s", r.RemoteAddr)
		http.ServeFile(w, r, "favicon.png")
		return
	}

	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "" {
			log.Printf("serving file %s to %s\n", r.URL.Path, r.RemoteAddr)

			md, ok := storage.Files[r.URL.Path]
			if ok == false {
				log.Printf("file not found %s", r.URL.Path)
				w.Header().Set("Content-Type", "text/html")
				return
			}

			w.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", md.Filename))
			http.ServeFile(w, r, path.Join(storage.Root, r.URL.Path, md.Filename))
			return
		}

		log.Printf("serving form to %s", r.RemoteAddr)

		w.Header().Set("Etag", fmt.Sprintf("\"%s\"", replypage_etag))
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Cache-Control", "max-age=31536000") // 1 year
		w.Write([]byte(frontpage))

	case http.MethodPost:
		r.ParseMultipartForm(32 << 20)

		duration, err := strconv.Atoi(r.FormValue("duration"))
		if err != nil {
			log.Printf("duration parse: %s", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}

		unit := r.FormValue("unit")
		switch unit {
		case "hour":
		case "day":
			duration *= 24
		case "week":
			duration *= 24 * 7
		default:
			log.Printf("unknown unit value: %s", unit)
			w.Header().Set("Content-Type", "text/html")
			return
		}

		in, handler, err := r.FormFile("filename")
		if err != nil {
			log.Printf("file upload: %v", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}
		defer in.Close()

		log.Println("uploading", handler.Filename)

		d, err := time.ParseDuration(fmt.Sprintf("%dh", duration))
		if err != nil {
			log.Printf("parse duration: %v", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}

		md := &MetaData{
			Filename: handler.Filename,
			Created:  time.Now(),
			Expire:   time.Now().Add(d),
		}

		out, err := storage.Create(md)
		if err != nil {
			log.Printf("storage create: %v", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}
		defer func() {
			if err := out.Close(); err != nil {
				log.Printf("close: %v", err)
			}
		}()

		if _, err := io.Copy(out, in); err != nil {
			log.Printf("copy: %v", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}

		t, err := template.New("reply").Parse(replypage)
		if err != nil {
			log.Printf("template parse: %v", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}

		shareURL := fmt.Sprintf("%s/%s", url, md.Hash)

		w.Header().Set("Content-Type", "text/html")
		if err := t.Execute(w, shareURL); err != nil {
			log.Printf("template exec: %v", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}
	}
}
