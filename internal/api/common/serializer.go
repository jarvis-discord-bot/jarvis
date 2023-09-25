package common

import (
	"encoding/json"
	"net/http"
)

type serializer interface {
	Encode(w http.ResponseWriter, r *http.Request, v interface{}) error
	ContentType(w http.ResponseWriter, r *http.Request) string
	Decode(w http.ResponseWriter, r *http.Request, v interface{}) error
}

type jsonSerializer struct{}

var _ serializer = (*jsonSerializer)(nil)

// JSON is an Encoder for JSON.
var JSON serializer = (*jsonSerializer)(nil)

func (j *jsonSerializer) Encode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (j *jsonSerializer) ContentType(w http.ResponseWriter, r *http.Request) string {
	return "application/json; charset=utf-8"
}

func (j *jsonSerializer) Decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
