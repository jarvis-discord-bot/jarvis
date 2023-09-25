package common

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}, encoder serializer) {
	w.Header().Set("Content-Type", encoder.ContentType(w, r))
	w.WriteHeader(status)
	if data == nil {
		return
	}

	if err := encoder.Encode(w, r, data); err != nil {
		logrus.WithField("error", err.Error()).Warning("error encoding value")
	}
}
