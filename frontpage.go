package main

import (
	"fmt"
	"net/http"
)

func frontPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}
