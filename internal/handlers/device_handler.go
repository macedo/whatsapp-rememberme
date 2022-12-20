package handlers

import (
	"log"
	"net/http"
)

func NewDevicePageHandler(w http.ResponseWriter, r *http.Request) {
	if err := renderPage(w, r, "new_device", "admin"); err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}
