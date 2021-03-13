package cache

import "github.com/chisty/miniurl/model"

type Cache interface {
	Set(key string, value *model.MiniURL) error
	Get(key string) (*model.MiniURL, error)
}
