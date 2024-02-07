package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/cyla00/monero-escrow/views"
)

func main() {
	yearNow := time.Now().Year()
	http.Handle("/", templ.Handler(views.Index(strconv.Itoa(yearNow))))
	http.Handle("/sign-up", templ.Handler(views.Signup()))
	http.Handle("/sign-in", templ.Handler(views.Signin()))
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// })

	log.Fatal(http.ListenAndServe(":3000", nil))
}
