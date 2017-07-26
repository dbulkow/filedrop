package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

//go:generate go run scripts/mkpage.go filepost.html

func filepost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	storage.Lock()
	defer storage.Unlock()

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

	d, err := time.ParseDuration(fmt.Sprintf("%dh", duration))
	if err != nil {
		log.Printf("parse duration: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	from, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("split: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	files := make([]File, 0)
	for _, fh := range r.MultipartForm.File["filenames"] {
		f := File{
			Name: fh.Filename,
			Type: fh.Header["Content-Type"][0],
		}
		files = append(files, f)
	}

	md := &MetaData{
		Type:    StorageFile,
		From:    from,
		Files:   files,
		Created: time.Now(),
		Expire:  time.Now().Add(d),
	}

	if err := storage.Mkdir(md); err != nil {
		log.Printf("mkdir: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	for i, fh := range r.MultipartForm.File["filenames"] {
		in, err := fh.Open()
		if err != nil {
			log.Printf("file upload: %v", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}
		defer in.Close()

		out, err := storage.Create(md, fh.Filename)
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

		nbytes, err := io.Copy(out, in)
		if err != nil {
			log.Printf("copy: %v", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}

		uploads.Inc()
		uploadBytes.Add(float64(nbytes))

		log.Printf("uploaded %d bytes for %s as %s", nbytes, fh.Filename, md.Hash)

		md.Files[i].Size = nbytes
	}

	if err := storage.WriteMeta(md); err != nil {
		log.Printf("writemeta: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	t, err := template.New("reply").Parse(filepost_html)
	if err != nil {
		log.Printf("template parse: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	shareURL := fmt.Sprintf("%s/retrieve/%s", rootURL, md.Hash)

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache")
	if err := t.Execute(w, shareURL); err != nil {
		log.Printf("template exec: %v", err)
		return
	}
}
