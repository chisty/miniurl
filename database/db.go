package database

import "github.com/chisty/shortlink/model"

type DB interface {
	Save(data *model.ShortLink) error
	Get(id string) (*model.ShortLink, error)
}
