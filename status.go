package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	human "github.com/dustin/go-humanize"
)

//go:generate go run scripts/mkpage.go status.html

func status(w http.ResponseWriter, r *http.Request) {
	storage.Lock()
	defer storage.Unlock()

	log.Printf("serving status to %s", r.RemoteAddr)

	from, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("split: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	type Files struct {
		Filename   string
		Size       string
		Downloaded string
	}

	type Dir struct {
		Hash   string
		Expire string
		Files  []*Files
	}

	type Status struct {
		Dirs []*Dir
	}

	status := Status{Dirs: make([]*Dir, 0)}

	for _, md := range storage.Dirs {
		if md.From == from {
			dir := &Dir{
				Expire: human.Time(md.Expire),
				Hash:   md.Hash,
			}

			dir.Files = make([]*Files, 0)
			for _, f := range md.Files {
				files := &Files{
					Filename: f.Name,
					Size:     human.Bytes(uint64(f.Size)),
				}
				if f.downloaded {
					files.Downloaded = "downloaded"
				}
				dir.Files = append(dir.Files, files)
			}

			status.Dirs = append(status.Dirs, dir)
		}
	}

	t, err := template.New("status").Parse(status_html)
	if err != nil {
		log.Printf("template parse: %v", err)
		w.Header().Set("Content-Type", "text/html")
		return
	}

	w.Header().Set("Etag", fmt.Sprintf("\"%s\"", status_etag))
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache")
	if err := t.Execute(w, status); err != nil {
		log.Printf("template exec: %v", err)
		return
	}
}
