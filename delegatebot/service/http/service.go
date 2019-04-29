package http

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger logrus.FieldLogger
}

func New(logger logrus.FieldLogger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s Service) Run() error {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong\n")
	})

	return http.ListenAndServe(":8080", nil)
}
