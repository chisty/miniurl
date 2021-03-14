package cache

import "github.com/chisty/miniurl/model"

//Cache is the cache service interface
type Cache interface {
	Set(key string, value *model.MiniURL) error
	Get(key string) (*model.MiniURL, error)
}
