package database

import "github.com/chisty/miniurl/model"

//DB is database interface
type DB interface {
	Save(data *model.MiniURL) error
	Get(id string) (*model.MiniURL, error)
}
