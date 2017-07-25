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

func postfile(w http.ResponseWriter, r *http.Request) {
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

	in, handler, err := r.FormFile("filename")
	if err != nil {
		log.Printf("file upload: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}
	defer in.Close()

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

	md := &MetaData{
		Type:     StorageFile,
		From:     from,
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

	nbytes, err := io.Copy(out, in)
	if err != nil {
		log.Printf("copy: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	uploads.Inc()
	uploadBytes.Add(float64(nbytes))

	log.Printf("uploaded %d bytes for %s as %s", nbytes, md.Filename, md.Hash)

	t, err := template.New("reply").Parse(replypage)
	if err != nil {
		log.Printf("template parse: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	shareURL := fmt.Sprintf("%s/%s", url, md.Hash)

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache")
	if err := t.Execute(w, shareURL); err != nil {
		log.Printf("template exec: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}
}
