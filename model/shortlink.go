package model

import (
	"encoding/json"
	"io"
)

type ShortLink struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	CreatedOn string `json:"createdOn"`
}

func (l *ShortLink) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(l)
}

func (l *ShortLink) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(l)
}
