package model

import (
	"encoding/json"
	"io"
)

type MiniURL struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	CreatedOn string `json:"createdOn"`
}

func (l *MiniURL) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(l)
}

func (l *MiniURL) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(l)
}
