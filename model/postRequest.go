package model

import (
	"encoding/json"
	"io"
)

type PostRequest struct {
	URL string `json:"url" validate:"required,url"`
}

func (rq *PostRequest) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(rq)
}

func (rq *PostRequest) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(rq)
}
