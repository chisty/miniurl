package service

import "github.com/chisty/miniurl/model"

//MiniURLSvc is the service interface.
type MiniURLSvc interface {
	Get(id string) (*model.MiniURL, error)
	Save(data *model.MiniURL) (*model.MiniURL, error)
}
