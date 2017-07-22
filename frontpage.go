package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

//go:generate go run scripts/mkpage.go frontpage.html
//go:generate go run scripts/mkpage.go replypage.html

/*
 mail address of sender
 upload filename text field
 - button to switch between upload and paste text
 expire date/time picker - or - duration picker (1 day, week, month) no longer
 mail address list (separate by spaces) text field
 submit button

- verify all email addresses at submit (smtp.Verify)
*/

func frontPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Etag", fmt.Sprintf("\"%s\"", replypage_etag))
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Cache-Control", "max-age=31536000") // 1 year
		w.Write([]byte(frontpage))

	case http.MethodPost:
		r.ParseMultipartForm(32 << 20)

		//fmt.Println(r.PostFormValue("sender"))
		//fmt.Println(r.PostFormValue("recipients"))

		in, handler, err := r.FormFile("filename")
		if err != nil {
			log.Printf("file upload: %v", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}
		defer in.Close()

		log.Println("uploading", handler.Filename)

		out, err := os.OpenFile(path.Join("./downloads", handler.Filename), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("openfile: %v", err)
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

		w.Header().Set("Content-Type", "text/html")
		if err := t.Execute(w, handler.Filename); err != nil {
			log.Printf("template exec: %v", err)
			w.Header().Set("Content-Type", "text/html")
			return
		}
	}
}
