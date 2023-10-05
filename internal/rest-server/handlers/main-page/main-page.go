package mainpage

import (
	"io/ioutil"
	"net/http"
)

const (
	pagePath = "static/main.html"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		html, err := ioutil.ReadFile(pagePath)
		if err != nil {
			http.Error(w, "Error sending page", http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte(html))
		if err != nil {
			http.Error(w, "Error sending page", http.StatusInternalServerError)
			return
		}
	}
}
