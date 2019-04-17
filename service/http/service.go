package http

import (
	"fmt"
	"net/http"
)

type Service struct{}

func (s Service) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	return http.ListenAndServe(":8080", nil)
}
