package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"

	human "github.com/dustin/go-humanize"
)

//go:generate go run scripts/mkpage.go fileget.html

func fileget(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) > 2 {
		log.Printf("malformed URL %s", r.URL.Path)
		w.Header().Set("Contents-Type", "text/html")
		return
	}

	storage.Lock()
	defer storage.Unlock()

	if len(parts) == 2 {
		log.Printf("serving file %s to %s\n", r.URL.Path, r.RemoteAddr)

		md, ok := storage.Dirs[parts[0]]
		if ok == false {
			log.Printf("file not found %s", r.URL.Path)
			w.Header().Set("Content-Type", "text/html")
			return
		}

		var ctype string
		var file File
		var fnum int
		for fnum, file = range md.Files {
			if file.Name == parts[1] {
				ctype = file.Type
				break
			}
		}
		if ctype == "" {
			log.Printf("file \"%s\" not found %s", file.Name, r.URL.Path)
			w.Header().Set("Content-Type", "text/html")
			return
		}

		downloads.Inc()
		writer := NewResponseWriterCounter(w)

		writer.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", file.Name))
		writer.Header().Set("Content-Type", ctype)
		http.ServeFile(writer, r, path.Join(storage.Root, md.Hash, file.Name))

		md.Files[fnum].downloaded = true

		return
	}

	type Files struct {
		URL      string
		Filename string
		Type     string
		Size     string
	}

	type Retrieve struct {
		Expire string
		Files  []*Files
	}

	page := Retrieve{}

	md, ok := storage.Dirs[r.URL.Path]
	if ok {
		page.Expire = human.Time(md.Expire)
		page.Files = make([]*Files, 0)

		for _, file := range md.Files {
			f := &Files{
				URL:      r.URL.Path + "/" + file.Name,
				Filename: file.Name,
				Type:     file.Type,
				Size:     human.Bytes(uint64(file.Size)),
			}
			page.Files = append(page.Files, f)
		}
	}

	t, err := template.New("retrieve").Parse(fileget_html)
	if err != nil {
		log.Printf("template parse: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	w.Header().Set("Etag", fmt.Sprintf("\"%s\"", fileget_etag))
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache")
	if err := t.Execute(w, page); err != nil {
		log.Printf("template exec: %v", err)
		return
	}
}
