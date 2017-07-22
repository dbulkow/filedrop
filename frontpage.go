package main

import (
	"fmt"
	"net/http"
)

/*
 upload filename text field
 expire date/time picker - or - duration picker (1 day, week, month) no longer
 mail address list (separate by spaces) text field
 submit button
*/

func frontPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}
