package routes

import (
	"fmt"
	"net/http"
)

func PostLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "login page")
}
