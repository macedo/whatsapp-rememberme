package handlers

import (
	"log"
	"net/http"
)

func AdminPageHandler(w http.ResponseWriter, r *http.Request) {
	if err := renderPage(w, r, "admin", ""); err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}
