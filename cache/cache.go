package cache

import "github.com/chisty/shortlink/model"

type Cache interface {
	Set(key string, value *model.ShortLink) error
	Get(key string) (*model.ShortLink, error)
}
